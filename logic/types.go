package logic

// User stores info about user
type Group struct {
	// ID     int
	Groupid int64 `gorm:"uniqueIndex"` // Telegram's chat id

}

// User stores info about user
type User struct {
	ID        int
	ChatID    int64 `gorm:"uniqueIndex"` // Telegram's chat id
	Wallmoney int64
	Bank      int64
	SubsCount uint          // Amount of current subscribtions
	Subs      []Publication `gorm:"many2many:user_subscribtion;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // User's subscriptions
	Admin     Admin         `gorm:"foreignKey:UserID"`
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
	Chatid     int64
	Msgid      int
	Nameid     int
	Createtime int `gorm:"default:now()"`
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
	Playid string
	Nameid int
	Chatid int64
	Userid int64 //用户
	Bet    float64
	Score  int64

	Createtime int `gorm:"default:now()"`
	Status     int `gorm:"default:1"` //下注
}

func (AddScore) TableName() string {
	return "scorelogs"
}

// User stores info about user
type Scorelogs struct {
	Userid      int64
	Playid      string
	Chatid      int64
	Nameid      int
	Bet         int
	Changescore int64
	Score       int64
	Status      int
	Createtime  int `gorm:"default:now()"`
	Endtime     int `gorm:"default:now()"`
	Details     string
}
