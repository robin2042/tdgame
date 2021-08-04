package gamemanage

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/aoyako/telegram_2ch_res_bot/games/baccarat"
	"github.com/aoyako/telegram_2ch_res_bot/games/niuniu"
	"github.com/aoyako/telegram_2ch_res_bot/logger"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

var (
	GAME_SCORE []int64  = []int64{500000, 5000000, 100000000, 3000000000, 5000000000, 10000000000, 100000000000}
	GAME_Title []string = []string{"一贫如洗", "专业杀猪", "小康之家", "腰缠万贯", "西厂总管", "富可敌国", "宇宙首富"}
)

type GameMainManage struct {
	games.Games
	stg    *storage.Storage
	rdb    *storage.CloudStore
	Tables map[string]games.GameTable // chatid<-->table

}

// NewController constructor of Controller
func NewGameManager(stg *storage.Storage, rds *storage.CloudStore) games.Games {

	return &GameMainManage{
		stg:    stg,
		rdb:    rds,
		Tables: map[string]games.GameTable{},
	}
}

//下注
func (g *GameMainManage) LoadGames() (bool, error) {
	// if g.bGameStation != GS_TK_CALL {
	// 	return true, nil
	// }

	return true, nil
}

func (g *GameMainManage) GetTable(nameid int, chatid int64, msgid int) games.GameTable {
	playid := fmt.Sprintf("%d%d", chatid, msgid)
	table := g.Tables[playid]
	if table != nil {
		return table
	}

	table = CreateTable(nameid, chatid, msgid)
	g.Tables[playid] = table
	table.SetRdb(g.rdb)
	table.SetDB(g.stg)

	return table
}

func (g *GameMainManage) GameBegin(nameid int, chatid int64, msgid int) int {
	playid := fmt.Sprintf("%d%d", chatid, msgid)

	table := g.GetTable(nameid, chatid, msgid)
	if table.GetStatus() != games.GS_TK_FREE { //存在就返回
		return table.GetStatus()
	}

	table.SetMsgID(msgid)

	round := &logic.Gamerounds{
		Playid: playid,
		Chatid: chatid,
		Msgid:  msgid,
		Nameid: nameid,
		Status: games.GS_TK_BET,
	}
	g.stg.SaveGameRound(round)

	return games.GS_TK_FREE

}

//判断能否开局
func (g *GameMainManage) NewGames(nameid, chatid int64) bool {

	start := g.stg.NewGames(int(nameid), chatid)
	return start == nil
}

//游戏结束，清理用户下注信息
func (g *GameMainManage) GameEnd(nameid, chatid int64, msgid int) error {

	table := g.GetTable(games.GAME_NIUNIU, chatid, msgid)
	scores := table.EndGame()
	logger.Info("回写数据库:", scores) //回写数据库
	delete(g.Tables, table.GetPlayID())

	return nil
}

//投注金额
func (g *GameMainManage) Bet(table games.GameTable, userid int64, area int) (bool, error) {

	if table.GetStatus() != games.GS_TK_PLAYING {
		return false, errors.New("已经开局,无法更改选择")
	}
	return table.Bet(userid, area)

}

func (g *GameMainManage) BetInfos(chatid int64, msgid int) ([]logic.Bets, error) {
	playid := fmt.Sprintf("%d%d", chatid, msgid)
	table := g.Tables[playid]
	return table.GetBetInfos()

}

//写分
func (g *GameMainManage) WriteUserScore(playid string, scores []logic.Scorelogs) error {
	return nil
}

//写分
func (g *GameMainManage) WriteUserRecords(playid string, scores []logic.Scorelogs) error {
	return g.stg.WriteUserRecords(playid, scores)
}

func (g *GameMainManage) GetRecords(nameid, chatid int64) (*logic.Way, int) {
	// return GetNiuniu_Record(g.rdb, nameid, chatid)
	return nil, 0
}

func (g *GameMainManage) AddScore(table games.GameTable, player games.PlayInfo, score float64) (int64, error) {

	board, _ := g.stg.Balance(player.UserID, table.GetChatID())
	player.WallMoney = board.Score //拿到钱
	player.Title = GetTitle(board.Score)

	ebet, err := table.AddScore(player, score)
	if err != nil {
		return 0, err
	} else {
		addscore := &logic.AddScore{
			Playid: table.GetPlayID(),
			Chatid: table.GetChatID(),
			Userid: player.UserID,
			Nameid: table.GetNameID(),
			Bet:    float64(ebet),
			Score:  player.WallMoney,
		}
		g.stg.AddScore(addscore)
	}

	logger.Info("下注：", player.UserID, ebet)

	return ebet, nil
}

func CreateTable(nameid int, chatid int64, msgid int) games.GameTable {
	playid := fmt.Sprintf("%d%d", chatid, msgid)
	var table games.GameTable
	if nameid == games.GAME_NIUNIU {
		table = new(niuniu.Niuniu)
	} else if nameid == games.GAME_BACCARAT {
		table = new(baccarat.Baccarat)
	}

	table.InitTable(playid, nameid, chatid)

	return table
}
func GenerateID(nameid int, chatid int64) string {
	strchatid := strconv.FormatInt(chatid, 10)
	timeUnix := time.Now().Unix()
	playid := fmt.Sprintf("%s%d", strchatid, timeUnix)

	return playid
}

func GetTitle(score int64) string {
	for i := 0; i < len(GAME_SCORE); i++ {
		if GAME_SCORE[i] <= score {
			return GAME_Title[i]
		}
	}

	return GAME_Title[0]
}
