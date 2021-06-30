package games

import "github.com/aoyako/telegram_2ch_res_bot/storage"

// Controller struct is used to access database
type GameManager struct {
}

// NewController constructor of Controller
func NewGameManager(stg *storage.Storage) *GameManager {
	return &GameManager{}
}
