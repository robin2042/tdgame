package games

import (
	"errors"
	"fmt"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
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

var betsinfo map[int]string = map[int]string{0: "🕒未选择", 1: "🐉青龙", 2: "🐅白虎", 3: "🦚朱雀", 4: "🐢玄武"}

type GameTable interface {
	GetChatID() int64
	GetPlayID() string
	SetMsgID(int)   //获取游戏状态
	GetStatus() int //获取游戏状态
	StartGame(int64) (bool, error)
	EndGame() (bool, error)
	Bet(int64, int) (bool, error) //用户,下注区域
	GetBetInfos() (*logic.Select, error)
	GetSettleInfos() (*logic.Records, error)
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

	Bets              map[PlayInfo]int64 //下注额
	Areas             map[PlayInfo]int   //下注区域
	m_lUserTianScore  map[PlayInfo]int   //天
	m_lUserDiScore    map[PlayInfo]int   //地
	m_lUserXuanScore  map[PlayInfo]int   //玄
	m_lUserHuangScore map[PlayInfo]int   //黄

	Changes         map[PlayInfo]int64 //胜负
	Historys        map[PlayInfo]int64 //历史开奖记录
	m_cbTimers      [5]int             //牛几倍率
	m_lUserWinScore map[int64]int64    //赢钱

	m_lUserReturnScore map[int64]int64 //赢钱

}

func (g *GameDesk) InitTable(playid string, nameid int, chatid int64) {
	g.PlayID = playid

	g.NameID = nameid
	g.ChatID = chatid

	g.Players = make(map[int64]PlayInfo) //在线用户
	g.Bets = make(map[PlayInfo]int64)
	g.Areas = make(map[PlayInfo]int)
	g.m_lUserTianScore = make(map[PlayInfo]int)
	g.m_lUserDiScore = make(map[PlayInfo]int)
	g.m_lUserXuanScore = make(map[PlayInfo]int)
	g.m_lUserHuangScore = make(map[PlayInfo]int)

	g.Changes = make(map[PlayInfo]int64)
	g.m_lUserWinScore = make(map[int64]int64)
	g.m_lUserReturnScore = make(map[int64]int64)
	g.GameStation = GS_TK_FREE
}

//清理表
func (g *GameDesk) UnInitTable() {

	for pi := range g.Areas {
		delete(g.Areas, pi)
	}

	for pi := range g.Changes {
		delete(g.Changes, pi)
	}

	for pi := range g.Bets {
		delete(g.Bets, pi)
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
func (g *GameDesk) GetBetInfos() (*logic.Select, error) {
	betinfo := &logic.Select{}
	betinfo.Players = make([]logic.Bets, 0)
	for k, v := range g.Areas {
		bet := logic.Bets{}
		bet.UserName = k.Name
		bet.FmtBetArea = betsinfo[v]
		betinfo.Players = append(betinfo.Players, bet)

	}
	return betinfo, nil
}

//结算信息
func (g *GameDesk) GetSettleInfos() (*logic.Records, error) {
	betinfo := &logic.Records{}
	var str string
	for i := 0; i < MAX_COUNT; i++ {

		if i == INDEX_BANKER {
			str += "🎴庄家"
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])
			str += "<br>"
		} else if i == INDEX_PLAYER1 {
			str += "🐉青龙"
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])
			str += ""
		} else if i == INDEX_PLAYER2 {
			str += "🐅白虎"
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])
			str += "<br>"
		} else if i == INDEX_PLAYER3 {
			str += "🦚朱雀"
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])
			str += ""
		} else if i == INDEX_PLAYER4 {
			str += "🐢玄武"
			str += GetCardValueEmoj(g.m_cbTableCardArray[i])
			str += "<br>"
		}

	}
	betinfo.Detail = str //牌局

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
		if i.UserID == userid {
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
	//结算
	g.CalculateScore()

	return true, nil
}

//结束游戏

func (g *GameDesk) EndGame() (bool, error) {

	return true, nil

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
	g.Areas[user] = area
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
	//计算积分
	//遍历下注人员
	for k, v := range g.Areas {
		if v == INDEX_PLAYER1 {
			if (ID_TIAN_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k.UserID] += g.Bets[k] * int64(g.m_cbTimers[1])
				g.m_lUserReturnScore[k.UserID] += g.Bets[k]

			} else {
				lUserLostScore[k.UserID] -= g.Bets[k] * int64(g.m_cbTimers[0])

			}
		}

		if v == INDEX_PLAYER2 {
			if (ID_DI_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k.UserID] += g.Bets[k] * int64(g.m_cbTimers[2])
				g.m_lUserReturnScore[k.UserID] += g.Bets[k]
			} else {
				lUserLostScore[k.UserID] -= g.Bets[k] * int64(g.m_cbTimers[0])
				// lBankerWinScore += m_lUserDiScore[i]*m_cbTimers[0] ;
			}
		}
		if v == INDEX_PLAYER3 {
			if (ID_XUAN_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k.UserID] += g.Bets[k] * int64(g.m_cbTimers[3])
				g.m_lUserReturnScore[k.UserID] += g.Bets[k]
			} else {
				lUserLostScore[k.UserID] -= g.Bets[k] * int64(g.m_cbTimers[0])
			}

		}
		if v == INDEX_PLAYER4 {
			if (ID_HUANG_MARK & cbWinner) > 0 {
				g.m_lUserWinScore[k.UserID] += g.Bets[k] * int64(g.m_cbTimers[4])
				g.m_lUserReturnScore[k.UserID] += g.Bets[k]
			} else {
				lUserLostScore[k.UserID] -= g.Bets[k] * int64(g.m_cbTimers[0])

			}
		}

		g.m_lUserWinScore[k.UserID] += lUserLostScore[k.UserID] //总成绩
		fmt.Println(lUserLostScore)                             //总输赢

	}

}
