package games

import (
	"errors"
	"fmt"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/leekchan/accounting"
)

var (
	ID_TIAN_MARK  byte = 0x01
	ID_DI_MARK    byte = 0x02
	ID_XUAN_MARK  byte = 0x04
	ID_HUANG_MARK byte = 0x08
	ID_QUAN_SHU   byte = 0x10
)

const (
	//索引定义
	INDEX_BANKER  = 0 //庄家索引
	INDEX_PLAYER1 = 1 //天
	INDEX_PLAYER2 = 2 //地
	INDEX_PLAYER3 = 3 //玄
	INDEX_PLAYER4 = 4 //黄

)

var betsinfo map[int]string = map[int]string{0: "🕒未选择", 1: "🐲青龙", 2: "🐯白虎", 3: "🦚朱雀", 4: "🐢玄武"}

type GameTable interface {
	GetChatID() int64
	GetPlayID() string
	SetMsgID(int)   //获取游戏状态
	GetStatus() int //获取游戏状态
	StartGame(int64) (bool, error)
	SettleGame() ([]logic.Scorelogs, error)
	EndGame() error

	Bet(int64, int) (bool, error)           //用户,下注区域
	GetStartInfos() (*logic.Select, error)  //显示下注人员
	GetBetInfos() ([]logic.Bets, error)     //下注信息
	GetSelectInfos() (*logic.Select, error) //显示下注人员
	GetSettleInfos() (*logic.Records, error)
	// WriteUserScore([]logic.ChangeScore) error //批量执行写分
}

type GameDesk struct {
	GameTable
	MsgID              int    //消息ID
	PlayID             string //局号
	ChatID             int64  //桌台号
	NameID             int
	GameStation        int       //游戏状态
	LastBetTime        time.Time //最后一次下注时间
	BeginTime          time.Time //开局时间
	StartTime          time.Time //开始游戏时间
	NextStartTime      time.Time
	m_cbTableCardArray [5][5]byte         //牌
	Players            map[int64]PlayInfo //在线用户

	Bets  map[int64]int64 //下注额
	Areas map[int64]int   //下注区域

	// Changes         map[PlayInfo]int64 //胜负
	Historys        map[PlayInfo]int64 //历史开奖记录
	m_cbTimers      [5]int             //牛几倍率
	m_lUserWinScore map[int64]int64    //赢钱

	m_lUserReturnScore map[int64]int64 //赢钱
	m_GameRecordArrary []byte          //路子

}

func (g *GameDesk) InitTable(playid string, nameid int, chatid int64) {
	g.PlayID = playid

	g.NameID = nameid
	g.ChatID = chatid

	g.Players = make(map[int64]PlayInfo) //在线用户
	g.Bets = make(map[int64]int64)
	g.Areas = make(map[int64]int)

	// g.Changes = make(map[PlayInfo]int64)
	g.m_lUserWinScore = make(map[int64]int64)
	g.m_lUserReturnScore = make(map[int64]int64)
	g.GameStation = GS_TK_FREE
}

//清理表
func (g *GameDesk) UnInitTable() {

	for pi := range g.Areas {
		delete(g.Areas, pi)
	}
	for pi := range g.Bets {
		delete(g.Bets, pi)
	}

	for pi := range g.Players {
		delete(g.Players, pi)
	}

}

//GameTable
func (g *GameDesk) SetPlayID(playid string) {
	g.PlayID = playid
}

func (g *GameDesk) GetChatID() int64 {
	return g.ChatID
}

//GameTable
func (g *GameDesk) GetPlayID() string {
	return g.PlayID
}

//下注信息
func (g *GameDesk) GetBetInfos() ([]logic.Bets, error) {
	s := make([]logic.Bets, 0)
	ac := accounting.Accounting{Symbol: "$"}

	for k, v := range g.Bets {
		var bet logic.Bets
		bet.Userid = k
		bet.UserName = g.Players[k].Name
		bet.Bet = v
		bet.FmtBet = ac.FormatMoney(v)
		s = append(s, bet)
	}
	return s, nil
}

//结算信息
func (g *GameDesk) GetSettleInfos() (*logic.Records, error) {
	betinfo := &logic.Records{}
	ac := accounting.Accounting{Symbol: "$"}

	for i := 0; i < MAX_COUNT; i++ {
		var str string
		if i == INDEX_BANKER {
			str += "🎴庄家 "
			str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
			str += " "
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])

		} else if i == INDEX_PLAYER1 {
			str += "🐲青龙 "
			str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
			str += " "
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])

		} else if i == INDEX_PLAYER2 {
			str += "🐯白虎 "
			str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
			str += " "
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])

		} else if i == INDEX_PLAYER3 {
			str += "🦚朱雀 "
			str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
			str += " "
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])

		} else if i == INDEX_PLAYER4 {
			str += "🐢玄武 "
			str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
			str += " "
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])

		}
		betinfo.Detail = append(betinfo.Detail, str)
	}
	for k := range g.Players {
		change := logic.ChangeScore{}
		change.UserName = g.Players[k].Name

		change.FmtArea = betsinfo[g.Areas[k]]

		if g.m_lUserWinScore[k] > 0 { //赢钱了

			str := fmt.Sprintf("*赢* \\+%s", ac.FormatMoney(g.m_lUserWinScore[k]))
			change.FmtChangescore = str
		} else {
			str := fmt.Sprintf("*输* ~\\%s~", ac.FormatMoney(g.m_lUserWinScore[k]))
			change.FmtChangescore = str
		}
		betinfo.Change = append(betinfo.Change, change)
	}

	betinfo.WaysCount = len(g.m_GameRecordArrary) //路子

	//天地玄黄
	for _, v := range g.m_GameRecordArrary {
		if (ID_TIAN_MARK & v) > 0 {
			betinfo.Ways.Tian += "● "
		} else {
			betinfo.Ways.Tian += "○ "
		}
		if (ID_DI_MARK & v) > 0 {
			betinfo.Ways.Di += "● "
		} else {
			betinfo.Ways.Di += "○ "
		}
		if (ID_XUAN_MARK & v) > 0 {
			betinfo.Ways.Xuan += "● "
		} else {
			betinfo.Ways.Xuan += "○ "
		}
		if (ID_HUANG_MARK & v) > 0 {
			betinfo.Ways.Huang += "● "
		} else {
			betinfo.Ways.Huang += "○ "
		}

	}

	return betinfo, nil
}

//开始
func (g *GameDesk) StartGame(userid int64) (bool, error) {
	if g.GameStation != GS_TK_FREE {
		return false, errors.New("已经开局请等待本局结束！")
	}
	if time.Now().Before(g.LastBetTime.Add(time.Second * 6)) {
		return false, errors.New("所有用户无操作6s后才能开始游戏")
	}

	var bfind bool
	for i := range g.Bets {
		if i == userid {
			bfind = true
			break
		}
	}
	if !bfind {
		return false, errors.New("您没有参与此游戏，无权更改游戏状态")
	}
	//记录牌局
	g.GameStation = GS_TK_PLAYING

	//发牌
	g.DispatchTableCard()

	return true, nil
}

//回写数据库
func (g *GameDesk) SettleGame() ([]logic.Scorelogs, error) {

	//结算
	g.CalculateScore()

	scores := make([]logic.Scorelogs, 0)
	fmt.Println(scores)

	for k, v := range g.m_lUserWinScore {
		score := logic.Scorelogs{
			Userid:      k,
			Playid:      g.PlayID,
			Chatid:      g.ChatID,
			Nameid:      g.NameID,
			Bet:         g.Bets[k],
			Changescore: g.m_lUserWinScore[k],
			Score:       g.Players[k].WallMoney,
			Status:      2,
		}
		fmt.Println(k, v, score)
		scores = append(scores, score)
	}

	return scores, nil
}

//结束游戏,清理本局变量

func (g *GameDesk) EndGame() error {

	g.UnInitTable()
	g.GameStation = GS_TK_FREE

	return nil
}

//开始
func (g *GameDesk) GetStatus() int {
	return g.GameStation
}

func (g *GameDesk) DispatchTableCard() {
	nums := GenerateRandomNumber(0, 52, 52)
	var ncard int
	for i := 0; i < GAME_PLAYER; i++ {

		for j := 0; j < MAX_COUNT; j++ {
			ncard++
			g.m_cbTableCardArray[i][j] = m_cbCardListData[nums[ncard]]

		}
	}

}

//开始
func (g *GameDesk) GetMsgID() int {
	return g.MsgID
}

//开始
func (g *GameDesk) SetMsgID(m int) {
	g.MsgID = m
}

func (g *GameDesk) Bet(userid int64, area int) (bool, error) {
	user, v := g.Players[userid]
	if !v {
		return false, errors.New("您没有下注")
	}
	if user.BetCount >= 3 {
		return false, errors.New("您已选择无法更改")
	}
	g.Areas[userid] = area
	user.BetCount++

	return true, nil
}
func (g *GameDesk) CalculateScore() {

	lUserLostScore := make(map[int64]int64)

	//推断赢家
	var cbWinner byte

	for i := 1; i <= INDEX_PLAYER4; i++ {
		var cbMarkType byte
		switch i {
		case 1:
			cbMarkType = ID_TIAN_MARK
		case 2:
			cbMarkType = ID_DI_MARK
		case 3:
			cbMarkType = ID_XUAN_MARK
		case 4:
			cbMarkType = ID_HUANG_MARK
		}
		if CompareCard(g.m_cbTableCardArray[i], g.m_cbTableCardArray[INDEX_BANKER], MAX_COUNT) {

			cbWinner |= cbMarkType
		} else {
			cbWinner = (cbWinner & (^cbMarkType + 1))
		}

	}

	for i := 0; i < MAX_COUNT; i++ {
		g.m_cbTimers[i] = GetTimes(g.m_cbTableCardArray[i], 5, MAX_MULTIPLE)
	}
	if len(g.m_GameRecordArrary) > 100 {
		g.m_GameRecordArrary = nil
	}
	g.m_GameRecordArrary = append(g.m_GameRecordArrary, cbWinner)

	//计算积分
	//遍历下注人员
	for k, v := range g.Areas {
		if v == INDEX_PLAYER1 {
			if (ID_TIAN_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k] += g.Bets[k] * int64(g.m_cbTimers[1])
				g.m_lUserReturnScore[k] += g.Bets[k]

			} else {
				lUserLostScore[k] -= g.Bets[k] * int64(g.m_cbTimers[0])

			}
		}

		if v == INDEX_PLAYER2 {
			if (ID_DI_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k] += g.Bets[k] * int64(g.m_cbTimers[2])
				g.m_lUserReturnScore[k] += g.Bets[k]

			} else {
				lUserLostScore[k] -= g.Bets[k] * int64(g.m_cbTimers[0])

			}
		}
		if v == INDEX_PLAYER3 {
			if (ID_XUAN_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k] += g.Bets[k] * int64(g.m_cbTimers[3])
				g.m_lUserReturnScore[k] += g.Bets[k]

			} else {
				lUserLostScore[k] -= g.Bets[k] * int64(g.m_cbTimers[0])

			}

		}
		if v == INDEX_PLAYER4 {
			if (ID_HUANG_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k] += g.Bets[k] * int64(g.m_cbTimers[4])
				g.m_lUserReturnScore[k] += g.Bets[k]

			} else {
				lUserLostScore[k] -= g.Bets[k] * int64(g.m_cbTimers[0])

			}
		}

		g.m_lUserWinScore[k] += lUserLostScore[k] //总成绩

		fmt.Println(lUserLostScore) //总输赢

	}

}

//获取下注列表,还么有选择,只能获取下注筹码的人
func (g *GameDesk) GetStartInfos() (*logic.Select, error) {

	sel := &logic.Select{}

	bets := make([]logic.Bets, 0)

	for k, _ := range g.Bets {
		var bet logic.Bets
		bet.Userid = k
		bet.UserName = g.Players[k].Name

		bet.FmtBetArea = betsinfo[g.Areas[k]]

		bets = append(bets, bet)
	}
	sel.Players = bets
	sel.Countdown = 60
	//天地玄黄
	for _, v := range g.m_GameRecordArrary {
		if (ID_TIAN_MARK & v) > 0 {
			sel.Ways.Tian += "● "
		} else {
			sel.Ways.Tian += "○ "
		}
		if (ID_DI_MARK & v) > 0 {
			sel.Ways.Di += "● "
		} else {
			sel.Ways.Di += "○ "
		}
		if (ID_XUAN_MARK & v) > 0 {
			sel.Ways.Xuan += "● "
		} else {
			sel.Ways.Xuan += "○ "
		}
		if (ID_HUANG_MARK & v) > 0 {
			sel.Ways.Huang += "● "
		} else {
			sel.Ways.Huang += "○ "
		}

	}
	return sel, nil
}

//获取下注列表,还么有选择,只能获取下注筹码的人
func (g *GameDesk) GetSelectInfos() (*logic.Select, error) {

	sel := &logic.Select{}

	bets := make([]logic.Bets, 0)

	for k, _ := range g.Bets {
		var bet logic.Bets
		bet.Userid = k
		bet.UserName = g.Players[k].Name

		if g.Areas[k] != 0 {
			bet.FmtBetArea = "✅" + betsinfo[g.Areas[k]]
		} else {
			bet.FmtBetArea = betsinfo[g.Areas[k]]
		}

		bets = append(bets, bet)
	}
	sel.Players = bets
	sel.Countdown = 60

	return sel, nil
}
