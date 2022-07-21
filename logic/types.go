package logic

// User stores info about user
type Group struct {
	// ID     int
	Groupid int64 `gorm:"uniqueIndex"` // Telegram's chat id

}

// User stores info about user
type User struct {
	Uid       string
	Userid    int64
	ChatID    int64 `gorm:"uniqueIndex"` // Telegram's chat id
	Wallmoney int64
	Bank      int64
	SubsCount uint // Amount of current subscribtions
	// Subs      []Publication `gorm:"many2many:user_subscribtion;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // User's subscriptions
	// Admin     Admin         `gorm:"foreignKey:UserID"`
}

// Admin stores info about admins
type Admin struct {
	ID     int
	UserID uint64
}

// Publication stores info about origin of data sent to user
type Publication struct {
	ID        int
	Board     string // 2ch board name
	Tags      string // Array of strings to search in thread title
	IsDefault bool   // Publication owner
	Type      string // File formats
	Alias     string // String alias
	Users     []User `gorm:"many2many:user_subscribtion;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Info stores addition information about bot
type Info struct {
	ID       int
	LastPost uint64 // Time of the latest post
}

// User stores info about user
type Gamerounds struct {
	Playid     string
	Chatid     int64 //用户ID
	Msgid      int
	Nameid     int
	Winarea    int
	Createtime string `gorm:"default:now()"`
	Status     int
}

// User stores info about user
type Bet struct {
	Playid string
	Nameid int
	Chatid int64 //用户
	Bet    int64
	Score  int64

	Createtime int `gorm:"default:now()"`
	Status     int `gorm:"default:1"` //下注
}

// User stores info about user
type AddScore struct {
	Messageid string //唯一性ID
	Uid       string //用户id
	Playid    string
	Nameid    int
	Chatid    int64
	Userid    int64 //用户
	Bet       float64
	Score     int64

	Createtime int `gorm:"default:now()"`
	Status     int `gorm:"default:1"` //下注
}

func (AddScore) TableName() string {
	return "scorelogs"
}

// User stores info about user
type Bets struct {
	Userid     int64
	UserName   string
	Title      string
	Bet        int64
	FmtBet     string //下注额格式化
	BetArea    int
	FmtBetArea string //下注格式化

}

// User stores info about user
type Scorelogs struct {
	Userid      int64
	Playid      string
	Chatid      int64
	Nameid      int
	Bet         int64
	Changescore int64
	Score       int64
	Area        int
	Status      int
	Createtime  string `gorm:"default:now()"`
	Endtime     string `gorm:"default:now()"`
	Details     string
}

// 签到表
// id int AI PK
// userid int
// chatid bigint
// score bigint
// createtime timestamp
type Signlogs struct {
	Userid     int64
	Chatid     int64
	Sign       int //签到金额
	Score      int64
	Createtime string `gorm:"default:now()"`
}

type Leaderboard struct {
	Userid int64
	Score  int64
	Win    int64
	Grades int //名次
}
type Records interface {
}

// User stores info about user
type Selects interface {
}

// // //开奖记录
// type Records struct {
// 	Detail    []string //庄闲牌
// 	Change    []ChangeScore
// 	Ways      *Way //路子
// 	WaysCount int
// }

//开奖记录
type BaccaratRecords struct {
	Records
	Detail    []string //庄闲牌
	Change    []ChangeScore
	Ways      string //路子
	WaysCount int
}

// User stores info about user
type Select struct {
	Selects
	Countdown int    //倒计时
	Players   []Bets //选择区域
	Ways      Way    //路子
	WaysCount int
}

// User stores info about user
type NiuNiuSelect struct {
	Selects
	Countdown int    //倒计时
	Players   []Bets //选择区域
	Ways      string //路子
	WaysCount int
}

// User stores info about user
type BaccaratSelect struct {
	Selects
	Countdown int    //倒计时
	Players   []Bets //选择区域
	Ways      string //路子
	WaysCount int
}

// ：🐯白虎 赢 +$4,0000
// 扣钱
type ChangeScore struct {
	UserName       string //名字
	Title          string //头衔
	Area           int    //下注
	FmtArea        string //下注格式化
	Winscore       int64
	Returncore     int64 //退回 金币
	FmtChangescore string
}

// 扣钱
type Way struct {
	Tian  string
	Di    string
	Xuan  string
	Huang string
}

//开奖记录
type History struct {
	Win []int
}

func (History) TableName() string {
	return "gamerounds"
}

type Cashlogs struct {
	Orderid     string
	Userid      string
	Targetid    string
	Changescore int64
	Score       int64
	Btype       int
	Delete      int
	Createtime  string `gorm:"default:now()"`
	Modifytime  string `gorm:"default:now()"`
}
