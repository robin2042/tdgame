package games

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
	"github.com/leekchan/accounting"
)

const (
	GAME_NIUNIU = 40022000
)

var betsinfo map[int]string = map[int]string{-1: "未选择", 0: "青龙", 1: "白虎", 2: "朱雀", 3: "玄武"}

// Controller struct is used to access database
const (

	//游戏状态
	GS_TK_FREE    = iota //等待开始
	GS_TK_BET            //下注状态
	GS_TK_PLAYING        //游戏进行
)

type GameTable interface {
	GetChatID() int64
	GetPlayID() string
	SetMsgID(int)   //获取游戏状态
	GetStatus() int //获取游戏状态
	StartGame(int64) (bool, error)
	Bet(int64, int64, int) (bool, string)
	GetBetInfos() (*logic.Select, error)
}

type PlayInfo struct {
	Name   string
	UserID int64
	Title  string //头衔，富可敌国 小康之家
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
	m_cbTableCardArray [5][5]byte //牌

	Bets     map[PlayInfo]int64 //下注额
	Area     map[PlayInfo]int   //下注区域
	Changes  map[PlayInfo]int64 //胜负
	Historys map[PlayInfo]int64 //历史开奖记录

}

type GameManage interface {
	LoadGames()
}

type Games interface {
	GameBegin(nameid, msgid int, chatid int64) int
	GameEnd(nameid, chatid int64) int
	GetTable(nameid int, chatid int64) GameTable //桌台
	Bet(table GameTable, userid int64, area int) (bool, error)
	AddScore(GameTable, PlayInfo, float64) (int64, int64, error) //下注额 下注总额 错误
	BetInfos(chatid int64) ([]logic.Bets, error)
}

type GameMainManage struct {
	Games
	stg    *storage.Storage
	Tables map[int64]GameTable // chatid<-->table

}

// NewController constructor of Controller
func NewGameManager(stg *storage.Storage) Games {

	return &GameMainManage{
		stg:    stg,
		Tables: map[int64]GameTable{},
	}
}

//下注
func (g *GameMainManage) LoadGames() (bool, error) {
	// if g.bGameStation != GS_TK_CALL {
	// 	return true, nil
	// }

	return true, nil
}

func (g *GameMainManage) GetTable(nameid int, chatid int64) GameTable {
	table := g.Tables[int64(chatid)]
	if table != nil {
		return table
	}

	table = CreateTable(nameid, chatid)
	g.Tables[chatid] = table

	return table
}

func (g *GameMainManage) SaveGameRounds(nameid int, chatid int64, playid string) bool {

	return g.stg.IsChatAdmin(chatid)

}

func (g *GameMainManage) GameBegin(nameid, msgid int, chatid int64) int {

	table := g.GetTable(GAME_NIUNIU, chatid)
	if table.GetStatus() != GS_TK_FREE { //存在就返回
		return table.GetStatus()
	}

	table.SetMsgID(msgid)

	round := &logic.Gamerounds{
		Playid: GenerateID(nameid, chatid),
		Chatid: chatid,
		Msgid:  msgid,
		Nameid: nameid,
		Status: GS_TK_BET,
	}
	g.stg.SaveGameRound(round)

	return GS_TK_FREE

}

//游戏结束，清理用户下注信息
func (g *GameMainManage) GameEnd(nameid, chatid int64) int {
	table := g.GetTable(GAME_NIUNIU, chatid)
	gamedesk := table.(*GameDesk)

	for pi := range gamedesk.Area {
		delete(gamedesk.Area, pi)
	}

	for pi := range gamedesk.Changes {
		delete(gamedesk.Changes, pi)
	}

	for pi := range gamedesk.Bets {
		delete(gamedesk.Bets, pi)
	}

	return 0
}

func (g *GameMainManage) Bet(table GameTable, userid int64, area int) (bool, error) {
	gamedesk := table.(*GameDesk)
	if gamedesk.GetStatus() != GS_TK_PLAYING {
		return false, nil
	}

	return true, nil

}

func (g *GameMainManage) AddScore(table GameTable, player PlayInfo, score float64) (int64, int64, error) {
	gamedesk := table.(*GameDesk)

	addscore := &logic.AddScore{
		Playid: gamedesk.PlayID,
		Chatid: gamedesk.ChatID,
		Userid: player.UserID,
		Nameid: gamedesk.NameID,
		Bet:    score,
	}

	betscore, err := g.stg.AddScore(addscore)
	if err != nil {
		return 0, 0, err
	}
	gamedesk.LastBetTime = time.Now()
	gamedesk.Bets[player] += betscore //下注
	gamedesk.Area[player] = -1        //未选择

	return betscore, gamedesk.Bets[player], err

}

//获取下注列表
func (g *GameMainManage) BetInfos(chatid int64) ([]logic.Bets, error) {
	table := g.Tables[chatid]
	gamedesk := table.(*GameDesk)

	s := make([]logic.Bets, 0, len(gamedesk.Bets))
	ac := accounting.Accounting{Symbol: "$"}

	for k, v := range gamedesk.Bets {
		var bet logic.Bets
		bet.Userid = k.UserID
		bet.UserName = k.Name
		bet.Bet = v
		bet.FmtBet = ac.FormatMoney(v)
		s = append(s, bet)
	}

	return s, nil

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
	for k, v := range g.Area {
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

//投注
//数据库先扣除
func (g *GameDesk) Bet(userid int64, score int64, area int) (bool, string) {

	return true, ""
}

func CreateTable(nameid int, chatid int64) GameTable {
	playid := GenerateID(nameid, chatid)

	table := new(GameDesk)
	table.SetPlayID(playid)
	table.NameID = nameid
	table.ChatID = chatid
	table.Bets = make(map[PlayInfo]int64)
	table.Area = make(map[PlayInfo]int)
	table.Changes = make(map[PlayInfo]int64)

	table.GameStation = GS_TK_FREE
	return table
}
func GenerateID(nameid int, chatid int64) string {
	strchatid := strconv.FormatInt(chatid, 10)
	timeUnix := time.Now().Unix()
	playid := fmt.Sprintf("%s%d", strchatid, timeUnix)

	return playid
}
