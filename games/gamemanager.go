package games

import (
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/storage"

	"gopkg.in/tucnak/telebot.v2"
)

const (
	GAME_NIUNIU = 40022000
)

// Controller struct is used to access database
const (

	//游戏状态
	GS_TK_FREE    = iota //等待开始
	GS_TK_CALL           //叫庄状态
	GS_TK_SCORE          //下注状态
	GS_TK_PLAYING        //游戏进行
)

type GameTable interface {
	GetStatus() int //获取游戏状态
	GameBegin()
	Bet()

	GameEnd()
}

type GameDesk struct {
	bGameStation int
	StartTime    time.Time
}

type GameManage interface {
	LoadGames()
}

type Games interface {
	GetTable(nameid, chatid int) GameTable //返回桌台

	HandleMessage(m *telebot.Message) (bool, error)
	// GetStatus() int                              //获取游戏状态
	// Bet(userid int64, score int64) (bool, error) // bet
	// StartGame()                                  //开始
	// EndGame()                                    //结束
}

type GameMainManage struct {
	Games
	Tables map[int64]GameTable // chatid<-->table
}

// NewController constructor of Controller
func NewGameManager(stg *storage.Storage) Games {
	return &GameMainManage{}
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

//下注
func (g *GameMainManage) GameBegin(nameid int, m *telebot.Message) (bool, error) {
	// if g.bGameStation != GS_TK_CALL {
	// 	return true, nil
	// }

	return true, nil
}

//下注
func (g *GameMainManage) HandleMessage(m *telebot.Message) (bool, error) {
	// if g.bGameStation != GS_TK_CALL {
	// 	return true, nil
	// }

	return true, nil
}

func (g *GameMainManage) GetTable(nameid, chatid int) GameTable {

	return nil
}
