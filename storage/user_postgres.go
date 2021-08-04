package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// UserPostgres is an implementation of storage.User
type UserPostgres struct {
	db  *gorm.DB
	cfg *InitDatabase
}

// NewUserPostgres constructor of UserPostgres struct
func NewUserPostgres(db *gorm.DB, cfg *InitDatabase) *UserPostgres {
	return &UserPostgres{
		db:  db,
		cfg: cfg,
	}
}

// Register adds user in databse
func (userStorage *UserPostgres) Register(user *logic.User) error {
	var count int64
	userStorage.db.Model(&logic.User{}).Where("userid = ?", user.Userid).Count(&count)
	if count == 0 {
		result := userStorage.db.Create(user)

		if result.Error != nil {
			return result.Error
		}

		return nil
	}

	return errors.New("User already exists")
}

// Unregister removes user from database
func (userStorage *UserPostgres) Unregister(user *logic.User) error {
	result := userStorage.db.Delete(user)
	return result.Error
}

// GetUserByChatID returns user by chat id
func (userStorage *UserPostgres) GetUserByChatID(chatID int64) (*logic.User, error) {
	var user logic.User
	var count int64
	userStorage.db.Model(&logic.User{}).Where("userid = ?", chatID).Count(&count)
	if count == 0 {
		return nil, errors.New("No user found")
	}
	result := userStorage.db.Where("userid = ?", chatID).First(&user)
	return &user, result.Error
}

// Update user
func (userStorage *UserPostgres) Update(user *logic.User) error {
	result := userStorage.db.Save(user)
	return result.Error
}

// GetUserByID returns user by it's id
func (userStorage *UserPostgres) GetUserByID(userID int64) (*logic.User, error) {
	var user logic.User
	var count int64
	userStorage.db.Model(&logic.User{}).Where("id = ?", userID).Count(&count)
	if count == 0 {
		return nil, errors.New("No user found")
	}
	result := userStorage.db.Where("id = ?", userID).First(&user)

	return &user, result.Error
}

// GetUsersByPublication returns subscribers of publication
func (userStorage *UserPostgres) GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) {
	var users []logic.User
	result := userStorage.db.Model(&pub).Association("Users").Find(&users)

	return users, result
}

// IsUserAdmin checks if user has administrator privileges
func (userStorage *UserPostgres) IsUserAdmin(user *logic.User) bool {
	// var count int64
	// userStorage.db.Model(&logic.Admin{}).Where("user_id = ?", user.ID).Count(&count)
	return false
}

// IsChatAdmin checks if user has administrator privileges by chatID
func (userStorage *UserPostgres) IsChatAdmin(chatID int64) bool {
	user, err := userStorage.GetUserByChatID(chatID)
	if err != nil {
		return false
	}
	return userStorage.IsUserAdmin(user)
}

func contains(slice []int64, val int64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func (userStorage *UserPostgres) Balance(userID, chatid int64) (*logic.Leaderboard, error) {
	var user logic.User
	var board logic.Leaderboard
	// uid := fmt.Sprintf("%d%d", userID, chatid)
	userStorage.db.Where("userid = ?", userID).First(&user)
	board.Userid = user.ChatID
	board.Score = user.Wallmoney
	board.Win = 0

	return &board, nil

}

//转账
func (userStorage *UserPostgres) Transfer(userID string, targetid string, payload int64) (int64, error) {

	var sourceuser logic.User
	var targetuser logic.User
	// var scorelog logic.Signlogs
	var ncount int64
	var rax float64 = 0.08

	tx := userStorage.db.Begin()
	defer tx.Commit()

	userStorage.db.Where("userid = ?", userID).First(&sourceuser).Count(&ncount)
	if ncount == 0 {
		return 0, errors.New("用户不存在")
	}
	userStorage.db.Where("userid = ?", targetid).First(&targetuser).Count(&ncount)
	if ncount == 0 {
		return 0, errors.New("用户不存在")
	}
	if sourceuser.Wallmoney < payload {
		return 0, errors.New("金额不足")
	}
	//游戏中不能转账
	userStorage.db.Model(&logic.Gamerounds{}).Where("chatid = ? and status=1", userID).Count(&ncount)
	if ncount > 0 {
		return 0, errors.New("游戏中无法转账")
	}

	raxscore := float64(payload) * rax //税率
	score := (float64(payload) - raxscore)

	result := userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).Update("wallmoney", gorm.Expr("wallmoney-?", payload))
	if result.Error != nil {
		tx.Rollback()
		return 0, errors.New("发生错误")
	}

	result = userStorage.db.Model(&logic.User{}).Where("userid = ?", targetid).Update("wallmoney", gorm.Expr("wallmoney+?", score))
	if result.Error != nil {
		tx.Rollback()
		return 0, errors.New("发生错误")
	}
	cashlog := logic.Cashlogs{
		Orderid:     OrderID(),
		Userid:      userID,
		Targetid:    targetid,
		Changescore: payload,
		Score:       sourceuser.Wallmoney,
		Btype:       1,
	}
	userStorage.db.Create(&cashlog)

	targetcashlog := logic.Cashlogs{
		Orderid:     OrderID(),
		Userid:      targetid,
		Targetid:    userID,
		Changescore: payload,
		Score:       targetuser.Wallmoney,
		Btype:       2,
	}
	userStorage.db.Create(&targetcashlog)

	return int64(raxscore), nil

}

// IsUserAdmin checks if user has administrator privileges
func (userStorage *UserPostgres) Sign(userID int, chatid int64, sign int) (int64, bool) {
	var user logic.User
	scorelog := logic.Signlogs{
		Userid: int64(userID),
		Chatid: chatid,
	}
	var ncount int64
	uid := fmt.Sprintf("%d%d", userID, chatid)
	userStorage.db.Where("uid = ?", uid).First(&user).Count(&ncount)
	if ncount == 0 {
		user.Uid = uid
		user.Userid = int64(userID)
		user.ChatID = chatid
		userStorage.Register(&user)
	}

	//没有签到过
	if err := userStorage.db.Where("userid  = ? order by createtime desc ", userID, chatid).Find(&scorelog).RowsAffected; err == 0 {

		result := userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).Update("wallmoney", gorm.Expr("wallmoney+?", sign))
		if result.Error != nil {
			return 0, false
		}

		scorelog.Score = user.Wallmoney
		scorelog.Sign = sign

		// 处理错误...
		userStorage.db.Create(scorelog)
		return user.Wallmoney + int64(sign), true
	} else {

		timer, _ := time.ParseInLocation("2006-01-02 15:04:05", scorelog.Createtime, time.Local)

		if time.Since(timer).Seconds() <= 150 {
			return 0, false
		}
		userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).Update("wallmoney", gorm.Expr("wallmoney+?", sign))
		var scorelog logic.Signlogs
		scorelog.Score = user.Wallmoney
		scorelog.Sign = sign
		scorelog.Userid = int64(userID)
		scorelog.Chatid = chatid

		// 处理错误...
		userStorage.db.Create(scorelog)
	}

	return user.Wallmoney, true
}

// 存钱
func (userStorage *UserPostgres) Deposit(userID int, payload int64) (int64, error) {

	tx := userStorage.db.Begin()
	defer tx.Commit()

	//游戏不能存钱
	var user logic.User
	var ncount int64
	result := userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).First(&user)
	if result.RowsAffected <= 0 {
		return 0, errors.New("not found")
	}
	if user.Wallmoney < payload {
		return user.Wallmoney, errors.New("钱不够")
	}

	//游戏中不能转账
	userStorage.db.Model(&logic.Gamerounds{}).Where("chatid = ? and status=1", userID).Count(&ncount)
	if ncount > 0 {
		return 0, errors.New("游戏中无法转账")
	}

	// logic.User{wallmoney: gorm.Expr("wallmoney-?", payload), Bank:gorm.Expr("bank+?", payload)}

	result = userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).Updates(map[string]interface{}{"wallmoney": gorm.Expr("wallmoney-?", payload), "Bank": gorm.Expr("bank+?", payload)})
	if result.Error != nil {
		tx.Rollback()
		return 0, errors.New("发生错误")
	}

	return user.Wallmoney - payload, nil
}

//取款
func (userStorage *UserPostgres) DrawMoney(userID int, payload int64) (int64, error) {

	tx := userStorage.db.Begin()
	defer tx.Commit()

	//游戏不能存钱
	var user logic.User

	result := userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).First(&user)
	if result.RowsAffected <= 0 {
		return 0, errors.New("not found")
	}
	if user.Bank < payload {
		return user.Bank, errors.New("钱不够")
	}

	// logic.User{wallmoney: gorm.Expr("wallmoney-?", payload), Bank:gorm.Expr("bank+?", payload)}

	result = userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).Updates(map[string]interface{}{"wallmoney": gorm.Expr("wallmoney+?", payload), "Bank": gorm.Expr("bank-?", payload)})
	if result.Error != nil {
		tx.Rollback()
		return 0, errors.New("发生错误")
	}

	return user.Wallmoney + payload, nil
}

func OrderID() string {
	orderid := fmt.Sprintf("%d", time.Now().Nanosecond())

	return orderid
}
