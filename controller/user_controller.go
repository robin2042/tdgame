package controller

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// UserController is an implementation of controller.User
type UserController struct {
	stg *storage.Storage
}

// NewUserController constructor of UserController struct
func NewUserController(stg *storage.Storage) *UserController {
	return &UserController{stg: stg}
}

// Register performs user registration
func (ucon *UserController) Register(chatID int64) error {

	user := &logic.User{
		Userid: chatID,
	}
	return ucon.stg.Register(user)
}

// Unregister performs user deregistration
func (ucon *UserController) Unregister(chatID int64) error {
	user := &logic.User{
		Userid: chatID,
	}
	return ucon.stg.Unregister(user)
}

// GetUsersByPublication returns subscribers of publication
func (ucon *UserController) GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) {
	users, err := ucon.stg.GetUsersByPublication(pub)

	return users, err
}

// GetUsersByPublication returns subscribers of publication
func (ucon *UserController) Sign(userid int, chatid int64, sign int) (int64, bool) {
	return ucon.stg.Sign(userid, chatid, sign)

}

// GetUsersByPublication returns subscribers of publication
func (ucon *UserController) Balance(chatID int64) (*logic.Leaderboard, error) {
	return ucon.stg.Balance(chatID)

}

// 转账
func (ucon *UserController) Transfer(userID int64, targetid int64, payload int64) error {

	return nil

}
