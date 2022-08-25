package initialize

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"tdgames/controller"
	"tdgames/downloader"
	"tdgames/gamemanage"
	"tdgames/storage"
	"tdgames/telegram"

	"github.com/spf13/viper"
)

// App initializes application
func App() (*telegram.TgBot, uint64) {
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config file: %s", err.Error())
	}
	fmt.Print(viper.Get("niuniu_start"))

	db, err := storage.NewMysqlDB(storage.Config{
		Host:     viper.GetString("db.host"),
		Port:     os.Getenv("DB_PORT"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		log.Fatalf("Error creating database: %s", err.Error())
	}

	admins := storage.InitDatabase{
		Admin: stringToInt64Slice(viper.GetStringSlice("tg.admin_id")),
	}

	Storage := storage.NewStorage(db, &admins)
	Rds := storage.ExampleNewClient()
	controller := controller.NewController(Storage)
	games := gamemanage.NewGameManager(Storage, Rds)

	bot := telegram.NewTelegramBot(os.Getenv("BOT_TOKEN"), controller, downloader.NewDownloader(
		viper.GetString("disk.path"),
		viper.GetUint64("disk.size")),
		games, Rds)

	telegram.SetupHandlers(bot)
	gamemanage.InitStart(bot, false) //初始化

	fmt.Println("链接bot成功!")

	return bot, viper.GetUint64("polling.time")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func stringToInt64Slice(data []string) []int64 {
	result := make([]int64, len(data))
	for key := range data {
		tmp, err := strconv.ParseInt(data[key], 10, 64)
		if err != nil {
			panic(err)
		}
		result[key] = tmp
	}
	return result
}
