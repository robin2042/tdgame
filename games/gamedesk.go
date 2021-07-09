package games

import (
	"errors"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

type GameTable interface {
	GetChatID() int64
	GetPlayID() string
	SetMsgID(int)   //获取游戏状态
	GetStatus() int //获取游戏状态
	StartGame(int64) (bool, error)
	Bet(int64, int) (bool, error) //用户,下注区域
	GetBetInfos() (*logic.Select, error)
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

	Bets     map[PlayInfo]int64 //下注额
	Areas    map[PlayInfo]int   //下注区域
	Changes  map[PlayInfo]int64 //胜负
	Historys map[PlayInfo]int64 //历史开奖记录

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

//开始
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

	return true, nil
}

//开始
func (g *GameDesk) GetStatus() int {
	return g.GameStation
}

func (g *GameDesk) DispatchTableCard() {
	nums := GenerateRandomNumber(0, 54, 54)
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
