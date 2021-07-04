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
	GetStatus() int //获取游戏状态
	StartGame() (bool, string)
	// Bet()

	// GameEnd()
}

type GameDesk struct {
	PlayID      string
	GameStation int
	StartTime   time.Time
}

type GameManage interface {
	LoadGames()
}

type Games interface {
	GameBegin(nameid int, chatid int64) bool
	GetTable(nameid int, chatid int64) GameTable //返回桌台

	// GetStatus() int                              //获取游戏状态
	// Bet(userid int64, score int64) (bool, error) // bet
	// StartGame()                                  //开始
	// EndGame()                                    //结束
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
func (g *GameMainManage) Bet(userid int64, score int64) (bool, error) {
	// if g.bGameStation != GS_TK_CALL {
	// 	return true, nil
	// }

	return true, nil
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

func (g *GameMainManage) GameBegin(nameid int, chatid int64) bool {

	table := g.GetTable(GAME_NIUNIU, chatid)
	_, playid := table.StartGame() //新开局
	fmt.Println(playid)
	round := &logic.GameRounds{
		Playid: GenerateID(nameid, chatid),
		Chatid: chatid,
		Nameid: nameid,
		Status: GS_TK_BET,
	}
	g.stg.SaveGameRound(round)

	return table.GetStatus() != GS_TK_FREE

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

	g.GameStation = GS_TK_BET
	return true, ""
}

//开始
func (g *GameDesk) GetStatus() int {
	return g.GameStation
}

func CreateTable(nameid int, chatid int64) GameTable {
	playid := GenerateID(nameid, chatid)

	table := new(GameDesk)
	table.SetPlayID(playid)
	table.GameStation = GS_TK_FREE
	return table
}
func GenerateID(nameid int, chatid int64) string {
	strchatid := strconv.FormatInt(chatid, 10)
	timeUnix := time.Now().Unix()
	playid := fmt.Sprintf("%s%d", strchatid, timeUnix)

	return playid
}
