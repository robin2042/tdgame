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

// Controller struct is used to access database
const (

	//游戏状态
	GS_TK_FREE    = iota //等待开始
	GS_TK_BET            //下注状态
	GS_TK_PLAYING        //游戏进行
)

type PlayInfo struct {
	Name     string
	UserID   int64
	BetCount int    //可以更改三次下注
	Title    string //头衔，富可敌国 小康之家
}

type GameManage interface {
	LoadGames()
}

type Games interface {
	GameBegin(nameid, msgid int, chatid int64) int
	GameEnd(nameid, chatid int64) error
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
func (g *GameMainManage) GameEnd(nameid, chatid int64) error {
	table := g.GetTable(GAME_NIUNIU, chatid)
	table.EndGame()

	// gamedesk := table.(*GameDesk)

	return nil
}

func (g *GameMainManage) Bet(table GameTable, userid int64, area int) (bool, error) {
	gamedesk := table.(*GameDesk)
	if gamedesk.GetStatus() != GS_TK_PLAYING {
		return false, errors.New("已经开局,无法更改选择")
	}
	gamedesk.Bet(userid, area)

	return true, nil

}

func (g *GameMainManage) AddScore(table GameTable, player PlayInfo, score float64) (int64, int64, error) {
	gamedesk := table.(*GameDesk)
	_, v := gamedesk.Players[player.UserID]

	//第一次增加
	if !v {
		gamedesk.Players[player.UserID] = player
	}

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

func CreateTable(nameid int, chatid int64) GameTable {
	playid := GenerateID(nameid, chatid)

	table := new(GameDesk)
	table.InitTable(playid, nameid, chatid)

	return table
}
func GenerateID(nameid int, chatid int64) string {
	strchatid := strconv.FormatInt(chatid, 10)
	timeUnix := time.Now().Unix()
	playid := fmt.Sprintf("%s%d", strchatid, timeUnix)

	return playid
}
