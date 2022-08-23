package gamemanage

import (
	"errors"
	"fmt"
	"strconv"
	"tdgames/games"
	"tdgames/games/dice"
	"tdgames/logger"
	"tdgames/logic"
	"tdgames/storage"
	"tdgames/telegram"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/tucnak/telebot.v2"
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

//创建停盘消息
func CloseBetTimer(tb *telegram.TgBot, ntimer int) {

	groupid := viper.GetInt64("tg.groupid")
	c := &telebot.Chat{
		ID: int64(groupid),
	}

	timer := time.NewTimer(time.Duration(ntimer) * time.Second)
	go func() {
		<-timer.C //等等定时器、
		table := tb.Games.GetTable(games.GAME_DICE, groupid, 0)
		fmt.Println(table)
		CloseBet_SendBack(tb, 0, c)
		RandDiceTimer(tb, 2) //封盘后10秒发送骰子

	}()

	fmt.Println("停盘消息")

}

//倒计时
func CountDownTimer(tb *telegram.TgBot, ntimer time.Duration) {
	groupid := viper.GetInt64("tg.groupid")
	c := &telebot.Chat{
		ID: int64(groupid),
	}

	table := tb.Games.GetTable(games.GAME_DICE, groupid, 0)
	fmt.Println(table)
	timer := time.NewTimer(ntimer)
	fmt.Println(timer)
	go func() {
		<-timer.C //等等定时器
		CountDown_SendBack(tb, 0, c)
		CloseBetTimer(tb, 10) //停盘消息

	}()
}

func RandDiceTimer(tb *telegram.TgBot, ntimer int) {
	groupid := viper.GetInt64("tg.groupid")
	c := &telebot.Chat{
		ID: int64(groupid),
	}

	table := tb.Games.GetTable(games.GAME_DICE, groupid, 0)
	fmt.Println(table)
	timer := time.NewTimer(time.Duration(ntimer) * time.Second)

	go func() {
		<-timer.C //等等定时器
		RandDice_SendBack(tb, 0, c)
		//发送消息
		LotteryGame_SendBack(tb, 0, c)
		//开启下一轮游戏
		//开奖结果
		//最近10期

	}()
}

//开奖结果
func LotteryGame_SendBack(tb *telegram.TgBot, groupid int, c *telebot.Chat) {
	table := tb.Games.GetTable(games.GAME_DICE, int64(groupid), 0)

	lottery := table.GetLotteryInfo()

	str := telegram.TemplateDice_LotteryText(lottery)

	tb.SendChatMessage(str, nil, c)

	//发送最后10个
	rhistory := "history"

	lasthistory, _ := tb.Rds.GetLrange(rhistory, 0, dice.TEN_HISTORY)
	dicehistory := logic.DiceHistory{
		Last:    lottery.Wins,
		Records: lasthistory,
	}
	historystr := telegram.TemplateDice_HistoryText(dicehistory)
	tb.SendChatMessage(historystr, nil, c)

}

//停盘消息
func CloseBet_SendBack(tb *telegram.TgBot, groupid int, c *telebot.Chat) {
	table := tb.Games.GetTable(games.GAME_DICE, int64(groupid), 0)
	period := table.GetPeriodInfo()
	bets, err := table.GetBetInfos()
	fmt.Println(bets, err)
	jettoninfo := logic.DiceJettonInfo{
		Info: period,
		Bets: bets,
	}

	str := telegram.TemplateDice_CloseBetText(jettoninfo)

	m, err := tb.SendChatMessage(str, nil, c)
	fmt.Println(m, err)
}

//封盘
func CountDown_SendBack(tb *telegram.TgBot, groupid int, c *telebot.Chat) {
	table := tb.Games.GetTable(games.GAME_DICE, int64(groupid), 0)
	str := telegram.TemplateDice_CountDownText()
	if table != nil {
		fmt.Println(table)
	}
	m, err := tb.SendChatMessage(str, nil, c)
	fmt.Println(m, err)
}

//随机发送骰子
func RandDice_SendBack(tb *telegram.TgBot, groupid int, c *telebot.Chat) {
	dice := &telebot.Dice{
		Type: telebot.Cube.Type,
	}

	first, _ := tb.DiceToMessage(c, dice)
	if first.Dice.Type != telebot.Cube.Type {
		return
	}
	second, _ := tb.DiceToMessage(c, dice)
	if first.Dice.Type != telebot.Cube.Type {
		return
	}
	three, _ := tb.DiceToMessage(c, dice)
	if first.Dice.Type != telebot.Cube.Type {
		return
	}

	table := tb.Games.GetTable(games.GAME_DICE, int64(groupid), 0)
	table.SettleGame(first.Dice.Value, second.Dice.Value, three.Dice.Value) //结算游戏

}

//启动游戏
func InitStart(tb *telegram.TgBot) {

	groupid := viper.GetInt64("tg.groupid")
	m := &telebot.Chat{
		ID: int64(groupid),
	}

	table := tb.Games.GetTable(games.GAME_DICE, groupid, 0)
	fmt.Println(table)
	periond, lasttime, _ := table.InitPeriodInfo()
	logger.Info("当前%s期,开局时间:%d!", periond, lasttime)
	newgametimer := time.NewTimer(time.Duration(lasttime) * time.Second)

	if table.GetStatus() > games.GS_TK_BET {

		tb.SendChatMessage("已经开局，请等待结束\\!", nil, m)
	}
	start := tb.Games.NewGames(games.GAME_DICE, m.ID)
	if start {
		logger.Info("新游戏开局:%s", m.ID)
	}

	go func() {
		logger.Info("当前时间为:", time.Now())
		fmt.Println("准备开局，当前时间为:", time.Now())
		<-newgametimer.C

		msg := telegram.TemplateDice_Text(periond)

		// reply := telegram.TemplateDice_Bet(tb)
		message := telebot.Message{Chat: m}

		_, err := tb.SendHtmlMessage(msg, nil, &message)
		if err != nil {
			logger.Error(err.Error())
		}
		start := tb.Games.NewGames(games.GAME_DICE, m.ID)
		if !start {
			logger.Error(err.Error())
		}

		tb.Games.GameBegin(games.GAME_DICE, message.Chat.ID, message.ID)
		table := tb.Games.GetTable(games.GAME_DICE, message.Chat.ID, message.ID)
		// table.SetPeriodInfo(periond)
		// fmt.Println(table)
		lasttime, _ := table.GetGameTimeSecond()
		// fmt.Println(lasttime)
		CountDownTimer(tb, lasttime) //30 倒计时
	}()

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
	// playid := "-7306125820"
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

func (g *GameMainManage) BetInfos(chatid int64, msgid int) ([]string, error) {
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

func (g *GameMainManage) AddScore(messageid string, table games.GameTable, player games.PlayInfo, area, score int) (int64, error) {

	board, _ := g.stg.Balance(player.UserID, table.GetChatID())
	player.WallMoney = board.Score //拿到钱
	// player.Title = GetTitle(area, score)

	ebet, err := table.AddScore(player, area, score)
	if err != nil {
		return 0, err
	} else {
		addscore := &logic.AddScore{
			Messageid: messageid,
			Playid:    table.GetPlayID(),
			Chatid:    table.GetChatID(),
			Userid:    player.UserID,
			Nameid:    table.GetNameID(),
			Bet:       int64(score),
			Score:     player.WallMoney,
		}
		g.stg.AddScore(addscore)
	}

	logger.Info("下注：", player.UserID, ebet)

	return ebet, nil
}

func CreateTable(nameid int, chatid int64, msgid int) games.GameTable {
	playid := fmt.Sprintf("%d%d", chatid, msgid)
	var table games.GameTable
	// if nameid == games.GAME_NIUNIU {
	// 	table = new(niuniu.Niuniu)
	// } else if nameid == games.GAME_BACCARAT {
	// 	table = new(baccarat.Baccarat)
	// } else
	if nameid == games.GAME_DICE {
		table = new(dice.Dice)
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

func GetTitle(area, score int) string {
	for i := 0; i < len(GAME_SCORE); i++ {
		// if GAME_SCORE[i] <= score {
		// 	return GAME_Title[i]
		// }
	}

	return GAME_Title[0]
}
