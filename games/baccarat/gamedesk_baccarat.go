package baccarat

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/aoyako/telegram_2ch_res_bot/storage"

	"github.com/leekchan/accounting"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

const (
	GAME_PLAYER = 2
	MAX_COUNT   = 3

	//索引定义
	INDEX_PLAYER = 0 //闲家索引
	INDEX_BANKER = 1 //庄家索引

	AREA_XIAN        = 0 //闲家索引
	AREA_PING        = 1 //平家索引
	AREA_ZHUANG      = 2 //庄家索引
	AREA_XIAN_TIAN   = 3 //闲天王
	AREA_ZHUANG_TIAN = 4 //庄天王
	AREA_TONG_DUI    = 5 //同点平
	AREA_XIAN_DUI    = 6 //闲对子
	AREA_ZHUANG_DUI  = 7 //庄对子
	AREA_MAX         = 8 //最大区域
)

//区域倍数multiple
const (
	MULTIPLE_XIAN        = 200  //闲家倍数
	MULTIPLE_PING        = 900  //平家倍数
	MULTIPLE_ZHUANG      = 200  //庄家倍数 195
	MULTIPLE_XIAN_TIAN   = 300  //闲天王倍数
	MULTIPLE_ZHUANG_TIAN = 300  //庄天王倍数
	MULTIPLE_TONG_DIAN   = 3300 //同点平倍数
	MULTIPLE_XIAN_PING   = 1200 //闲对子倍数
	MULTIPLE_ZHUANG_PING = 1200 //庄对子倍数

)

//掩码
const (
	ID_CONTROL_MASK_PLAYER      = 0x01
	ID_CONTROL_MARK_PING        = 0x02
	ID_CONTROL_MARK_BANKER      = 0x04
	ID_CONTROL_MASK_PAIR_PLAYER = 0x08
	ID_CONTROL_MARK_PAIR_BANKER = 0x10
)

// 投注选择庄胜🔴、闲胜🔵、和🟢、庄对🟠、闲对🟣其中一个
var betsinfo map[int]string = map[int]string{-1: "🕒未选择", 0: "🔵闲", 1: "🟢平", 2: "🔴庄", 6: "🟣闲对", 7: "🟠庄对"}

//百家乐
type Baccarat struct {
	games.GameDesk
	m_cbCardCount      [2]byte    //扑克数目
	m_cbTableCardArray [2][3]byte //桌面扑克

}

func (g *Baccarat) AddScore(player games.PlayInfo, score float64) (int64, error) {
	return g.GameDesk.AddScore(player, score)

}
func (g *Baccarat) Bet(userid int64, area int) (bool, error) {
	return g.GameDesk.Bet(userid, area)
}

func (g *Baccarat) EndGame() error {

	g.UnInitTable()
	g.GameStation = games.GS_TK_FREE

	return nil
}

func (g *Baccarat) SettleGame(userid int64) ([]logic.Scorelogs, error) {

	_, err := g.GameDesk.SettleGame(userid)
	if err != nil {
		return nil, err
	}

	//结算
	g.CalculateScore()
	g.GameDesk.WriteChangeScore(g.PlayID, g.ChatID, g.M_lUserReturnScore) //回写数据库

	return nil, nil
}

//结算信息
func (g *Baccarat) GetSettleInfos() (logic.Records, error) {
	betinfo := &logic.BaccaratRecords{}
	ac := accounting.Accounting{Symbol: "$"}
	//天地 投注选择庄胜🔴、闲胜🔵、和🟢、庄对🟠、闲对🟣其中一个
	for i := GAME_PLAYER - 1; i >= 0; i-- {
		var str string
		if i == INDEX_BANKER {
			str += "🔴庄 "
			str += GetCardTimesEmoj(g.m_cbTableCardArray[i], MAX_COUNT)
			str += " "
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])

		} else if i == INDEX_PLAYER {
			str += "🔵闲 "
			str += GetCardTimesEmoj(g.m_cbTableCardArray[i], MAX_COUNT)
			str += " "
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])

		}
		betinfo.Detail = append(betinfo.Detail, str)
	}
	for k := range g.Players {
		change := logic.ChangeScore{}
		change.UserName = g.Players[k].Name
		change.Title = g.Players[k].Title
		change.FmtArea = betsinfo[g.Areas[k]]

		if v, ok := g.M_lUserWinScore[k]; ok {
			if g.M_lUserWinScore[k] > 0 { //赢钱了

				str := fmt.Sprintf("*赢* \\+%s", ac.FormatMoney(v))
				change.FmtChangescore = str
			} else {
				str := fmt.Sprintf("*输* ~\\%s~", ac.FormatMoney(v))
				change.FmtChangescore = str
			}
		} else {
			str := fmt.Sprintf("*返回* \\+%s", ac.FormatMoney(g.Bets[k]))
			change.FmtChangescore = str
		}

		betinfo.Change = append(betinfo.Change, change)
	}
	ways, count := g.GetRecords()
	betinfo.Ways = ways
	betinfo.WaysCount = count

	return betinfo, nil
}

//下注信息
//获取下注列表,还么有选择,只能获取下注筹码的人
func (g *Baccarat) GetSelectInfos() (*logic.Select, error) {

	sel := &logic.Select{}

	bets := make([]logic.Bets, 0)

	for k, _ := range g.Bets {
		var bet logic.Bets
		bet.Userid = k
		bet.Title = g.Players[k].Title //头衔
		bet.UserName = g.Players[k].Name

		if g.Areas[k] != 0 {
			bet.FmtBetArea = "✅" + betsinfo[g.Areas[k]]
		} else {
			bet.FmtBetArea = betsinfo[g.Areas[k]]
		}

		bets = append(bets, bet)
	}
	sel.Players = bets

	ncountdown := time.Until(g.BetCountDownTime)
	if int(ncountdown.Seconds()) < 0 {
		sel.Countdown = 0
	} else {
		sel.Countdown = int(ncountdown.Seconds())
	}

	return sel, nil
}

//下注信息
func (g *Baccarat) GetBetInfos() (bets []logic.Bets, err error) {
	return g.GameDesk.GetBetInfos()

}

func (g *Baccarat) InitTable(playid string, nameid int, chatid int64) {
	g.GameDesk.InitTable(playid, nameid, chatid)
}

func (g *Baccarat) DispatchTableCard() {

	// 洗牌
	nums := GenerateRandomNumber(0, len(m_cbCardListData), len(m_cbCardListData))
	var ncard int
	for i := 0; i < GAME_PLAYER; i++ {

		for j := 0; j < MAX_COUNT; j++ {
			ncard++
			g.m_cbTableCardArray[i][j] = 0
			g.m_cbTableCardArray[i][j] = m_cbCardListData[nums[ncard]]

		}
	}
	// fmt.Printf("组:%d,发牌:%d", g.ChatID, g.m_cbTableCardArray)
	//7 2
	// [[55 11 56] [6 44 10]]
	// 		[0]	61 '='	unsigned char
	// 		[1]	18 '\x12'	unsigned char
	// 		[2]	40 '('	unsigned char
	// -		[1]	0x0cfb2d49 "\x19\x18\a...	unsigned char[3]
	// 		[0]	25 '\x19'	unsigned char
	// 		[1]	24 '\x18'	unsigned char
	// 		[2]	7 '\a'	unsigned char

	// g.m_cbTableCardArray[INDEX_PLAYER] = [3]byte{61, 18, 40}
	// g.m_cbTableCardArray[INDEX_BANKER] = [3]byte{25, 24, 7}

	//首次发牌
	g.m_cbCardCount[INDEX_PLAYER] = 2
	g.m_cbCardCount[INDEX_BANKER] = 2

	//计算点数
	cbBankerCount := GetCardListPip(g.m_cbTableCardArray[INDEX_BANKER], g.m_cbCardCount[INDEX_BANKER])
	cbPlayerTwoCardCount := GetCardListPip(g.m_cbTableCardArray[INDEX_PLAYER], g.m_cbCardCount[INDEX_PLAYER])
	//闲家补牌
	var cbPlayerThirdCardValue byte //第三张牌点数
	if cbPlayerTwoCardCount <= 5 && cbBankerCount < 8 {
		//计算点数
		g.m_cbCardCount[INDEX_PLAYER]++
		cbPlayerThirdCardValue = GetCardPip(g.m_cbTableCardArray[INDEX_PLAYER][2])
	}
	//庄家补牌
	if cbPlayerTwoCardCount < 8 && cbBankerCount < 8 {
		switch cbBankerCount {
		case 0:
		case 1:
		case 2:
			g.m_cbCardCount[INDEX_BANKER]++

		case 3:
			if (g.m_cbCardCount[INDEX_PLAYER] == 3 && cbPlayerThirdCardValue != 8) || g.m_cbCardCount[INDEX_PLAYER] == 2 {
				g.m_cbCardCount[INDEX_BANKER]++
			}
			break

		case 4:
			if (g.m_cbCardCount[INDEX_PLAYER] == 3 && cbPlayerThirdCardValue != 1 && cbPlayerThirdCardValue != 8 && cbPlayerThirdCardValue != 9 && cbPlayerThirdCardValue != 0) || g.m_cbCardCount[INDEX_PLAYER] == 2 {
				g.m_cbCardCount[INDEX_BANKER]++
			}
			break
		case 5:
			if (g.m_cbCardCount[INDEX_PLAYER] == 3 && cbPlayerThirdCardValue != 1 && cbPlayerThirdCardValue != 2 && cbPlayerThirdCardValue != 3 && cbPlayerThirdCardValue != 8 && cbPlayerThirdCardValue != 9 && cbPlayerThirdCardValue != 0) || g.m_cbCardCount[INDEX_PLAYER] == 2 {
				g.m_cbCardCount[INDEX_BANKER]++
			}
			break

		case 6:
			if g.m_cbCardCount[INDEX_PLAYER] == 3 && (cbPlayerThirdCardValue == 6 || cbPlayerThirdCardValue == 7) {
				g.m_cbCardCount[INDEX_BANKER]++
			}
			break

			//不须补牌
		case 7:
		case 8:
		case 9:
			break
		default:
			break
		}
	}
}

//开始
func (g *Baccarat) StartGame(userid int64) (bool, error) {
	result, err := g.GameDesk.StartGame(userid)
	if err != nil {
		return result, err
	}
	//发牌
	g.DispatchTableCard()
	return true, nil
}

//结算

func (g *Baccarat) CalculateScore() int64 {

	//计算牌点
	cbPlayerCount := GetCardListPip(g.m_cbTableCardArray[INDEX_PLAYER], g.m_cbCardCount[INDEX_PLAYER])
	cbBankerCount := GetCardListPip(g.m_cbTableCardArray[INDEX_BANKER], g.m_cbCardCount[INDEX_BANKER])

	//系统输赢
	// LONGLONG lSystemScore = 0l;
	var cbWinArea [AREA_MAX]bool
	var cbWinner int
	g.DeduceWinner(&cbWinArea) //判断赢家

	if cbWinArea[AREA_XIAN] {
		cbWinner = ID_CONTROL_MASK_PLAYER
	}
	if cbWinArea[AREA_PING] {
		cbWinner = ID_CONTROL_MARK_PING
	}
	if cbWinArea[AREA_ZHUANG] {
		cbWinner = ID_CONTROL_MARK_BANKER
	}
	if cbWinArea[AREA_XIAN_DUI] {
		cbWinner = ID_CONTROL_MASK_PAIR_PLAYER
	}
	if cbWinArea[AREA_ZHUANG_DUI] {
		cbWinner = ID_CONTROL_MARK_PAIR_BANKER
	}

	//区域倍率
	var cbMultiple [AREA_MAX]int = [AREA_MAX]int{MULTIPLE_XIAN, MULTIPLE_PING, MULTIPLE_ZHUANG,
		MULTIPLE_XIAN_TIAN, MULTIPLE_ZHUANG_TIAN, MULTIPLE_TONG_DIAN,
		MULTIPLE_XIAN_PING, MULTIPLE_ZHUANG_PING}

	fmt.Println(cbMultiple, cbPlayerCount, cbBankerCount)

	for k := range g.GameDesk.Bets {
		area := g.Areas[k]

		if cbWinArea[area] { //投注这里
			g.M_lUserReturnScore[k] += (g.Bets[k] * int64((cbMultiple[area]-100)/100.00)) + g.Bets[k] //赢钱
			g.M_lUserWinScore[k] += (g.Bets[k] * int64((cbMultiple[area]-100)/100.00))

		}
		//总的分数
		for k := range g.Players {
			//没有下注
			if _, ok := g.M_lUserWinScore[k]; !ok {
				g.M_lUserReturnScore[k] = g.Bets[k]
			}
		}
		key := fmt.Sprintf("%d%d", g.ChatID, g.NameID)
		g.Rdb.RPush(key, cbWinner)
		fmt.Println(g.M_lUserWinScore)

	}

	return 0

}

//获取下注列表,还么有选择,只能获取下注筹码的人
func (g *Baccarat) GetStartInfos() (logic.Selects, error) {

	sel := &logic.BaccaratSelect{}

	bets := make([]logic.Bets, 0)

	for k, _ := range g.Bets {
		var bet logic.Bets
		bet.Userid = k
		bet.UserName = g.Players[k].Name
		bet.Title = g.Players[k].Title
		bet.FmtBetArea = betsinfo[g.Areas[k]]

		bets = append(bets, bet)
	}
	sel.Players = bets
	ncountdown := time.Until(g.BetCountDownTime)
	sel.Countdown = int(ncountdown.Seconds())
	way, count := g.GetRecords()
	sel.Ways = way
	sel.WaysCount = count
	return sel, nil
}

//推断赢家
func (g *Baccarat) DeduceWinner(pWinArea *[AREA_MAX]bool) {
	//计算牌点
	cbPlayerCount := GetCardListPip(g.m_cbTableCardArray[INDEX_PLAYER], g.m_cbCardCount[INDEX_PLAYER])
	cbBankerCount := GetCardListPip(g.m_cbTableCardArray[INDEX_BANKER], g.m_cbCardCount[INDEX_BANKER])

	//胜利区域--------------------------
	//平
	if cbPlayerCount == cbBankerCount {
		pWinArea[AREA_PING] = true

		// 同平点
		if g.m_cbCardCount[INDEX_PLAYER] == g.m_cbCardCount[INDEX_BANKER] {
			var wCardIndex byte
			for wCardIndex = 1; wCardIndex < g.m_cbCardCount[INDEX_PLAYER]; wCardIndex++ {
				cbBankerValue := games.GetCardValue(g.m_cbTableCardArray[INDEX_BANKER][wCardIndex])
				cbPlayerValue := games.GetCardValue(g.m_cbTableCardArray[INDEX_PLAYER][wCardIndex])
				if cbBankerValue != cbPlayerValue {
					break
				}
			}

			if wCardIndex == g.m_cbCardCount[INDEX_PLAYER] {
				pWinArea[AREA_TONG_DUI] = true
			}
		}
	} else if cbPlayerCount < cbBankerCount { // 庄
		pWinArea[AREA_ZHUANG] = true

		//天王判断
		if cbBankerCount == 8 || cbBankerCount == 9 {
			pWinArea[AREA_ZHUANG_TIAN] = true
		}
	} else // 闲
	{
		pWinArea[AREA_XIAN] = true

		//天王判断
		if cbPlayerCount == 8 || cbPlayerCount == 9 {
			pWinArea[AREA_XIAN_TIAN] = true
		}
	}
	//对子判断
	if games.GetCardValue(g.m_cbTableCardArray[INDEX_PLAYER][0]) == games.GetCardValue(g.m_cbTableCardArray[INDEX_PLAYER][1]) {
		pWinArea[AREA_XIAN_DUI] = true
	}
	if games.GetCardValue(g.m_cbTableCardArray[INDEX_BANKER][0]) == games.GetCardValue(g.m_cbTableCardArray[INDEX_BANKER][1]) {
		pWinArea[AREA_ZHUANG_DUI] = true
	}
}

func GetBaccarat_Record(s *storage.CloudStore, nameid, chatid int64) (string, int) {

	key := fmt.Sprintf("%d%d", chatid, nameid)
	var way string //路子
	scores, err := s.LRange(key, 0, 10)
	if err != nil {
		return "", 0
	}

	if len(scores) >= 15 {
		s.Del(key)
	}

	for _, v := range scores {
		//fmt.Println(v)
		k, _ := strconv.Atoi(v)
		// //天地玄黄
		if byte(k) == ID_CONTROL_MASK_PLAYER {
			way += "🔵 "
		} else if byte(k) == ID_CONTROL_MARK_PING {
			way += "🟢 "
		} else if byte(k) == ID_CONTROL_MARK_BANKER {
			way += "🔴 "
		}

	}

	return way, len(scores)
}

//获取游戏记录
func (g *Baccarat) GetRecords() (string, int) {
	return GetBaccarat_Record(g.Rdb, int64(g.NameID), g.ChatID) //游戏ID
}
