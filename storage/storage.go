package storage

import (
	"tdgames/logic"

	"gorm.io/gorm"
)

// User interface defines methods for User Storage
type Group interface {
	GroupRegister(user *logic.Group) error   // Adds user in databse
	UnGroupRegister(user *logic.Group) error // Removes user from database
	// GetUserByChatID(chatID int64) (*logic.User, error)                  // Returns user by chat id
	// Update(user *logic.User) error                                      // Updates user
	// GetUserByID(userID int64) (*logic.User, error)                      // Returns user by it's id
	// GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) // Returns owner of publication
	// IsUserAdmin(user *logic.User) bool
	// IsChatAdmin(userID int64) bool
}

// User interface defines methods for User Storage
type User interface {
	Register(user *logic.User) error                                    // Adds user in databse
	Unregister(user *logic.User) error                                  // Removes user from database
	GetUserByChatID(chatID int64) (*logic.User, error)                  // Returns user by chat id
	Update(user *logic.User) error                                      // Updates user
	GetUserByID(userID int64) (*logic.User, error)                      // Returns user by it's id
	GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) // Returns owner of publication
	IsUserAdmin(user *logic.User) bool
	IsChatAdmin(userID int64) bool
	Sign(userID int, chatid int64, sign int) (int64, bool)                 //签到
	Balance(userID, chatid int64) (*logic.Leaderboard, error)              //余额
	Transfer(userID string, targetid string, payload int64) (int64, error) //转账
	Deposit(userID int, payload int64) (int64, error)                      // 存钱
	DrawMoney(userID int, payload int64) (int64, error)                    //取款
}

// Subscription interface defines methods for Publicaiton Storage
type Subscription interface {
	Add(user *logic.User, publication *logic.Publication) error       // Adds new subscription to user with publication
	AddDefault(publication *logic.Publication) error                  // Adds new subscription to user with publication
	Remove(publication *logic.Publication) error                      // Removes existing sybscription
	Disonnect(user *logic.User, publication *logic.Publication) error // Disonnect user from publication
	Update(user *logic.User, publication *logic.Publication) error    // Updates selected subscription
	GetSubsByUser(user *logic.User) ([]logic.Publication, error)      // Returns list of user's subscriptions
	GetAllSubs() []logic.Publication                                  // Returns all publications
	GetAllDefaultSubs() []logic.Publication
	Connect(user *logic.User, publication *logic.Publication) error
}

// Info interface definces methods for Info Storage
type Info interface {
	GetLastTimestamp() uint64    // Returns time of the latest post
	SetLastTimestamp(tsp uint64) // Sets time of the latest post
}

// Info interface definces methods for Info Storage
type Games interface {
	NewGames(nameid int, chatid int64) error
	SaveGameRound(game *logic.Gamerounds) error
	AddScore(game *logic.AddScore) (int64, error)
	BetInfos(playid string) ([]logic.Scorelogs, error)
	WriteChangeScore(string, int64, map[int64]int64) error
	WriteUserRecords(string, []logic.Scorelogs) error
	GetRecords(nameid int, chatid int64) []logic.Records
}

// Storage struct is used to access database
type Storage struct {
	User
	Subscription
	Info
	Group
	Games
}

// NewStorage constructor of Storage
func NewStorage(db *gorm.DB, cfg *InitDatabase) *Storage {
	return &Storage{
		User:         NewUserPostgres(db, cfg),
		Subscription: NewSubscriptionPostgres(db),
		Info:         NewInfoPostgres(db),
		Group:        NewGroupMysql(db),
		Games:        NewGamesMysql(db),
	}
}
