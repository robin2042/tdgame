package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB constructor of gorm.DB databse with postgresql database
func NewPostgresDB(cfg Config) (*gorm.DB, error) {
	connStr := formatPostgresConfig(cfg)
	return gorm.Open(mysql.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

// Formats config struct to meet gorm's expectations
func formatPostgresConfig(cfg Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}

// MigrateDatabase migrates database
func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&logic.User{}, &logic.Admin{}, &logic.Publication{}, &logic.Info{})

	if err != nil {
		log.Fatalf("Error migrating database")
	}

	var count int64
	db.Find(&logic.Info{}).Count(&count)
	if count == 0 {
		db.Create(&logic.Info{
			LastPost: uint64(time.Now().Unix()),
		})
	} else {
		var info logic.Info
		db.Find(&logic.Info{}).First(&info)
		info.LastPost = uint64(time.Now().Unix())
		db.Save(info)
	}
}
