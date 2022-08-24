package games

import (
	"tdgames/logic"
	"time"

	"tdgames/storage"
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

type GameTable interface {
	SetRdb(*storage.CloudStore)
	SetDB(*storage.Storage)
	InitTable(string, int, int64)
	GetChatID() int64
	GetPlayID() string
	GetNameID() int
	SetMsgID(string) //获取游戏状态
	GetStatus() int  //获取游戏状态
	StartGame(int64) (bool, error)
	SettleGame(int, int, int) ([]logic.Scorelogs, error)
	EndGame() error
	AddScore(player PlayInfo, area, score int) (int64, error)
	Bet(int64, int) (bool, error)           //用户,下注区域
	GetStartInfos() (logic.Selects, error)  //显示下注人员
	GetBetInfos() ([]string, error)         //下注信息
	GetBetInfo(int64) ([]string, int)       //下注信息
	GetSelectInfos() (logic.Selects, error) //显示下注人员
	GetSettleInfos() ([]string, error)
	GetPeriodInfo() logic.PeriodInfo //开局信息
	GetLotteryInfo() logic.LotteryInfo
	InitPeriodInfo(period string) (logic.PeriodInfo, int, error) //设置开局信息
	GetTitlesInfo() (string, error)                              //获取标题信息
	GetBalance(uid int64) int64                                  //获取用户余额
	GetGameTimeSecond() (time.Duration, time.Duration)           //倒计时时间,停盘,停止时间
}

type Games interface {
	NewGames(nameid, chatid int64) bool //判断上一句时间
	GameBegin(nameid int, chatid int64, msgid string) int
	GameEnd(nameid, chatid int64, msgid string) error
	GetTable(nameid int, chatid int64, msgid string) GameTable //桌台
	Bet(table GameTable, userid int64, area int) (bool, error)
	AddScore(string, GameTable, PlayInfo, int, int) (int64, error) //下注额 下注总额 错误
	BetInfos(chatid int64, msgid int) ([]string, error)
	WriteGameRounds(string, int) error
	WriteUserScore(string, []logic.Scorelogs) error
	WriteUserRecords(string, []logic.Scorelogs) error
	GetRecords(nameid, chatid int64) (*logic.Way, int)
}

const (
	GAME_NIUNIU   = 40022000
	GAME_BACCARAT = 40023000 //百家乐
	GAME_REDBLACK = 40024000 //红黑
	GAME_DRAGEN   = 40025000 //龙虎
	GAME_BENZIN   = 40026000 //奔驰宝马
	GAME_FRUIT    = 40027000 //水果机
	GAME_ROULE    = 40028000 //轮盘
	GAME_DICE     = 40029000
)

// Controller struct is used to access database
const (

	//游戏状态
	GS_TK_FREE    = iota //等待开始
	GS_TK_BET            //下注状态
	GS_TK_PLAYING        //游戏进行
)
