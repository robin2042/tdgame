package storage

import (
	"errors"

	"tdgames/logic"

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
	var count int64
	groupStorage.db.Model(&logic.Group{}).Where("groupid = ?", user.Groupid).Count(&count)
	if count == 0 {
		result := groupStorage.db.Create(user)

		if result.Error != nil {
			return result.Error
		}

		return nil
	}

	return errors.New("User already exists")
}

// Register adds user in databse
func (groupStorage *GroupMysql) UnGroupRegister(user *logic.Group) error {

	result := groupStorage.db.Delete(user)
	return result.Error
}
