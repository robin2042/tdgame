package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// UserPostgres is an implementation of storage.User
type GamesMysql struct {
	db *gorm.DB
}

// NewUserPostgres constructor of UserPostgres struct
func NewGamesMysql(db *gorm.DB) *GamesMysql {
	return &GamesMysql{
		db: db,
	}
}

// Register adds user in databse
func (groupStorage *GamesMysql) SaveGameRound(game *logic.Gamerounds) error {

	result := groupStorage.db.Create(game)
	return result.Error
}
