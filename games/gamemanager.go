package games

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logger"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

var (
	GAME_Title map[int]string = map[int]string{5000000: "一贫如洗", 100000000: "专业杀猪", 5000000000: "西厂总管", 10000000000: "富可敌国", 100000000000: "宇宙首富"}
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
	Name      string
	UserID    int64
	WallMoney int64
	BetCount  int    //可以更改三次下注
	Title     string //头衔，富可敌国 小康之家
}

type GameManage interface {
	LoadGames()
}

type Games interface {
	NewGames(nameid, chatid int64) bool //判断上一句时间
	GameBegin(nameid int, chatid int64, msgid int) int
	GameEnd(nameid, chatid int64, msgid int) error
	GetTable(nameid int, chatid int64, msgid int) GameTable //桌台
	Bet(table GameTable, userid int64, area int) (bool, error)
	AddScore(GameTable, PlayInfo, float64) (int64, error) //下注额 下注总额 错误
	BetInfos(chatid int64, msgid int) ([]logic.Bets, error)
	WriteGameRounds(string, int) error
	WriteUserScore(string, []logic.Scorelogs) error
	WriteUserRecords(string, []logic.Scorelogs) error
	GetRecords(nameid, chatid int64) (*logic.Way, int)
}

type GameMainManage struct {
	Games
	stg    *storage.Storage
	rdb    *storage.CloudStore
	Tables map[string]GameTable // chatid<-->table

}

// NewController constructor of Controller
func NewGameManager(stg *storage.Storage, rds *storage.CloudStore) Games {

	return &GameMainManage{
		stg:    stg,
		rdb:    rds,
		Tables: map[string]GameTable{},
	}
}

//下注
func (g *GameMainManage) LoadGames() (bool, error) {
	// if g.bGameStation != GS_TK_CALL {
	// 	return true, nil
	// }

	return true, nil
}

func (g *GameMainManage) GetTable(nameid int, chatid int64, msgid int) GameTable {
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

	table := g.GetTable(GAME_NIUNIU, chatid, msgid)
	if table.GetStatus() != GS_TK_FREE { //存在就返回
		return table.GetStatus()
	}

	table.SetMsgID(msgid)

	round := &logic.Gamerounds{
		Playid: playid,
		Chatid: chatid,
		Msgid:  msgid,
		Nameid: nameid,
		Status: GS_TK_BET,
	}
	g.stg.SaveGameRound(round)

	return GS_TK_FREE

}

//判断能否开局
func (g *GameMainManage) NewGames(nameid, chatid int64) bool {

	start := g.stg.NewGames(int(nameid), chatid)
	if start == nil {
		return true
	}
	return false
}

//游戏结束，清理用户下注信息
func (g *GameMainManage) GameEnd(nameid, chatid int64, msgid int) error {

	table := g.GetTable(GAME_NIUNIU, chatid, msgid)
	scores := table.EndGame()
	logger.Info("回写数据库:", scores) //回写数据库
	delete(g.Tables, table.GetPlayID())

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
	return GetNiuniu_Record(g.rdb, nameid, chatid)

}

func (g *GameMainManage) AddScore(table GameTable, player PlayInfo, score float64) (int64, error) {

	board, _ := g.stg.Balance(player.UserID)
	player.WallMoney = board.Score //拿到钱

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

func CreateTable(nameid int, chatid int64, msgid int) GameTable {
	playid := fmt.Sprintf("%d%d", chatid, msgid)

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
