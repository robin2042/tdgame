package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// UserPostgres is an implementation of storage.User
type GroupMysql struct {
	db *gorm.DB
}

// NewUserPostgres constructor of UserPostgres struct
func NewGroupMysql(db *gorm.DB) *GroupMysql {
	return &GroupMysql{
		db: db,
	}
}

// Register adds user in databse
func (groupStorage *GroupMysql) GroupRegister(user *logic.Group) error {
	return nil
}
