package games

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

const (
	GAME_NIUNIU = 40022000
)

// Controller struct is used to access database
const (

	//游戏状态
	GS_TK_FREE    = iota //等待开始
	GS_TK_BET            //下注状态
	GS_TK_PLAYING        //游戏进行
)

type GameTable interface {
	GetMsgID() int  //获取游戏状态
	SetMsgID(int)   //获取游戏状态
	GetStatus() int //获取游戏状态
	StartGame() (bool, string)
	Bet(int64, int64, int) (bool, string)

	// GameEnd()
}

type PlayInfo struct {
	Name   string
	UserID int64
	Title  string //头衔，富可敌国 小康之家
}

type GameDesk struct {
	GameTable
	MsgID         int //消息ID
	PlayID        string
	ChatID        int64
	NameID        int
	GameStation   int
	StartTime     time.Time
	NextStartTime time.Time
	Bets          map[PlayInfo]int64 //下注额
	Area          map[PlayInfo]int64 //下注区域
	Changes       map[PlayInfo]int64 //胜负

}

type GameManage interface {
	LoadGames()
}

type Games interface {
	GameBegin(nameid, msgid int, chatid int64) int
	GetTable(nameid int, chatid int64) GameTable //桌台
	Bet(table *GameDesk, userid int64, area int) bool
	AddScore(table GameTable, userid int64, score float64) error
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
	_, playid := table.StartGame() //新开局

	fmt.Println(playid)
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

func (g *GameMainManage) Bet(table *GameDesk, userid int64, area int) bool {
	addscore := &logic.AddScore{
		Playid: table.PlayID,
		Chatid: userid,
		Nameid: table.NameID,
	}
	g.stg.AddScore(addscore)
	return true

}

func (g *GameMainManage) AddScore(table GameTable, userid int64, score float64) error {
	gamedesk := table.(*GameDesk)

	addscore := &logic.AddScore{
		Playid: gamedesk.PlayID,
		Chatid: gamedesk.ChatID,
		Userid: userid,
		Nameid: gamedesk.NameID,
		Bet:    score,
	}
	g.stg.AddScore(addscore)
	return nil

}

//GameTable
func (g *GameDesk) SetPlayID(playid string) {
	g.PlayID = playid
}

//开始
func (g *GameDesk) StartGame() (bool, string) {
	if g.GameStation != GS_TK_FREE {
		return false, ""
	}
	//记录牌局
	later := time.Now()
	g.StartTime = later
	g.NextStartTime = later.Add(time.Second * 90) //90S后
	g.GameStation = GS_TK_BET
	return true, ""
}

//开始
func (g *GameDesk) GetStatus() int {
	return g.GameStation
}

// GetMsgID() int  //获取游戏状态
// 	SetMsgID(int)   //获取游戏状态

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
	table.GameStation = GS_TK_FREE
	return table
}
func GenerateID(nameid int, chatid int64) string {
	strchatid := strconv.FormatInt(chatid, 10)
	timeUnix := time.Now().Unix()
	playid := fmt.Sprintf("%s%d", strchatid, timeUnix)

	return playid
}
