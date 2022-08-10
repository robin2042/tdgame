package telegram

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"tdgames/controller"
	"tdgames/downloader"
	"tdgames/games"
	"time"
	"unsafe"

	"tdgames/logic"

	"github.com/xfrr/goffmpeg/transcoder"

	telebot "gopkg.in/tucnak/telebot.v2"
)

var fileHandlersQueue = make(chan bool, 100)

// MessageSender defines interface for bot-sender
type MessageSender interface {
	Send(r telebot.Recipient, value interface{}, args ...interface{}) (*telebot.Message, error)
	Edit(msg telebot.Editable, what interface{}, options ...interface{}) (*telebot.Message, error)
	Respond(c *telebot.Callback, resp ...*telebot.CallbackResponse) error
	Handle(interface{}, interface{})
	Start()
}

// TgBot represents telegram bot view
type TgBot struct {
	Me         *telebot.User
	Bot        MessageSender
	Controller *controller.Controller
	Downloader *downloader.Downloader
	Games      games.Games
}

// NewTelegramBot constructor of TelegramBot
func NewTelegramBot(token string, cnt *controller.Controller, d *downloader.Downloader, g games.Games) *TgBot {
	settings := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 30 * time.Second},
	}

	// If token is empty, do not send request
	// Developers of telebot lib made "offline" mode unaccessible
	// so reflect and unsafe is used to change that field
	if token == "" {
		rs := reflect.ValueOf(settings)
		rs2 := reflect.New(rs.Type()).Elem()
		rs2.Set(rs)
		rsf := rs2.FieldByName("offline")
		rsf = reflect.NewAt(rsf.Type(), unsafe.Pointer(rsf.UnsafeAddr())).Elem()
		rsf.SetBool(true)

		settings = rs2.Interface().(telebot.Settings)
	}

	bot, err := telebot.NewBot(settings)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return nil
	}

	return &TgBot{
		Me:         bot.Me,
		Bot:        bot,
		Controller: cnt,
		Downloader: d,
		Games:      g,
	}
}

// SetupHandlers to default values
func SetupHandlers(tb *TgBot) {
	//组
	tb.Bot.Handle(telebot.OnAddedToGroup, OnBotAddGroups(tb))
	tb.Bot.Handle(telebot.OnUserJoined, EnterGroups(tb))
	tb.Bot.Handle(telebot.OnUserLeft, LeaveGroups(tb))
	tb.Bot.Handle(telebot.OnCallback, Callback(tb))
	tb.Bot.Handle(telebot.OnText, Ontext(tb))
	//games
	// 	hl -【欢乐牛牛🎴】
	// zz -【转账💰】
	// hb -【红包🧧】
	// rank -【富豪榜🏆】
	// sheng -【胜场榜🚩】
	// ck -【存款💵】
	// qk -【取款💴】
	tb.Bot.Handle("/hl", NiuniuBet(tb))   //百人牛牛
	tb.Bot.Handle("/bj", BaccaratBet(tb)) //百家乐
	tb.Bot.Handle("/hh", BaccaratBet(tb)) //红黑
	tb.Bot.Handle("/lh", BaccaratBet(tb)) //龙虎
	tb.Bot.Handle("/bc", BaccaratBet(tb)) //奔驰宝马
	tb.Bot.Handle("/sl", BaccaratBet(tb)) //森林舞会
	tb.Bot.Handle("/sl", FruitBet(tb))    //水果机
	tb.Bot.Handle("/lp", RouletteBet(tb)) //轮盘

	//功能
	tb.Bot.Handle("/zz", GamesZZ(tb))
	tb.Bot.Handle("/hb", GamesHB(tb))
	tb.Bot.Handle("/rank", GamesRank(tb))
	tb.Bot.Handle("/sheng", GamesWins(tb))
	tb.Bot.Handle("/ck", GamesDeposit(tb))
	tb.Bot.Handle("/qk", GamesWithdraw(tb))

	tb.Bot.Handle("/relief", LeaveGroups(tb))
	tb.Bot.Handle("/kj", KAIJU(tb))

}

// Send files to users
func (tb *TgBot) Send(users []*logic.User, path, caption string) {
	if len(users) == 0 {
		return
	}
	var file telebot.Sendable
	if strings.HasSuffix(path, ".mp4") {
		file = &telebot.Video{File: telebot.FromURL(path), Caption: caption}
	}

	if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") || strings.HasSuffix(path, ".gif") {
		file = &telebot.Photo{File: telebot.FromURL(path), Caption: caption}
	}

	if strings.HasSuffix(path, ".webm") {
		err := tb.Downloader.Save(path)
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			err := tb.Downloader.Free(path)
			if err != nil {
				log.Println(err)
			}
			err = tb.Downloader.Free(strings.TrimSuffix(path, ".webm") + ".mp4")
			if err != nil {
				log.Println(err)
			}
		}()

		newVidPath, err := convertWebmToMp4(tb.Downloader, path)
		if err != nil {
			log.Println(err)
			return
		}

		file = &telebot.Video{File: telebot.FromDisk(newVidPath), Caption: caption}
	}

	for _, user := range users {
		for {
			fileHandlersQueue <- true

			_, err := tb.Bot.Send(&telebot.Chat{
				ID: int64(user.ChatID),
			}, file)

			<-fileHandlersQueue

			if err != nil {
				if e, ok := err.(telebot.FloodError); ok {
					time.Sleep(time.Duration(e.RetryAfter) * time.Second)
					continue
				} else {
					log.Println(err)
				}
			}
			break
		}
	}
}

func convertWebmToMp4(d *downloader.Downloader, path string) (string, error) {
	trans := new(transcoder.Transcoder)

	vidPath := d.Get(path)
	newVidPath := strings.TrimSuffix(vidPath, ".webm") + ".mp4"

	err := trans.Initialize(vidPath, newVidPath)
	if err != nil {
		return "", err
	}
	done := trans.Run(false)
	err = <-done
	if err != nil {
		log.Println(err)
		return "", err
	}

	return newVidPath, nil
}
