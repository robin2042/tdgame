package storage

import (
	"errors"
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
	userStorage.db.Model(&logic.User{}).Where("userid = ?", user.ChatID).Count(&count)
	if count == 0 {
		result := userStorage.db.Create(user)

		if result.Error != nil {
			return result.Error
		}

		// Adds admin record, if admin added
		if _, contains := contains(userStorage.cfg.Admin, user.ChatID); contains {
			result := userStorage.db.Create(&logic.Admin{
				UserID: uint64(user.ID),
			})

			if result.Error != nil {
				return result.Error
			}
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
	var count int64
	userStorage.db.Model(&logic.Admin{}).Where("user_id = ?", user.ID).Count(&count)
	return count != 0
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

func (userStorage *UserPostgres) Balance(userID int64) (*logic.Leaderboard, error) {
	var user logic.User
	var board logic.Leaderboard

	userStorage.db.Where("userid = ?", userID).First(&user)
	board.Userid = user.ChatID
	board.Score = user.Wallmoney
	board.Win = 0

	return &board, nil

}

func (userStorage *UserPostgres) Transfer(userID int64, targetid int64, payload int64) error {

	var sourceuser logic.User
	var targetuser logic.User
	// var scorelog logic.Signlogs
	var ncount int64
	var rax float64 = 0.08

	tx := userStorage.db.Begin()
	defer tx.Commit()

	userStorage.db.Where("userid = ?", userID).First(&sourceuser).Count(&ncount)
	if ncount == 0 {
		return errors.New("用户不存在")
	}
	userStorage.db.Where("userid = ?", userID).First(&targetuser).Count(&ncount)
	if ncount == 0 {
		return errors.New("用户不存在")
	}
	if sourceuser.Wallmoney < payload {
		return errors.New("金额不足")
	}
	score := sourceuser.Wallmoney * int64(rax)

	result := userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).Update("wallmoney", gorm.Expr("wallmoney-?", score))
	if result.Error != nil {
		tx.Rollback()
		return errors.New("发生错误")
	}

	result = userStorage.db.Model(&logic.User{}).Where("userid = ?", targetid).Update("wallmoney", gorm.Expr("wallmoney+?", score))
	if result.Error != nil {
		tx.Rollback()
		return errors.New("发生错误")
	}
	
	return nil

}

// IsUserAdmin checks if user has administrator privileges
func (userStorage *UserPostgres) Sign(userID int, chatid int64, sign int) (int64, bool) {
	var user logic.User
	var scorelog logic.Signlogs
	var ncount int64
	userStorage.db.Where("userid = ?", userID).First(&user).Count(&ncount)
	if ncount == 0 {
		user.Userid = int64(userID)
		user.ChatID = chatid
		userStorage.Register(&user)
	}

	//没有签到过
	if err := userStorage.db.Where("userid  = ? order by createtime desc ", userID).Find(&scorelog).RowsAffected; err == 0 {

		result := userStorage.db.Model(&logic.User{}).Where("userid = ?", userID).Update("wallmoney", gorm.Expr("wallmoney+?", sign))
		if result.Error != nil {
			return 0, false
		}

		scorelog.Score = user.Wallmoney
		scorelog.Sign = sign
		scorelog.Userid = int64(userID)

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

		// 处理错误...
		userStorage.db.Create(scorelog)
	}

	return user.Wallmoney, true
}
