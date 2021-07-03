package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aoyako/telegram_2ch_res_bot/initialize"
)

func main() {
	log.Println("Starting...")
	bot, apicnt, games, duration := initialize.App()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go bot.Bot.Start()
	go initialize.StartPolling(apicnt, duration)

	go games.LoadGames()
	log.Println("Started")

	<-quit
	log.Println("Quit")
}

// package main

// import (
// 	"time"

// 	"gopkg.in/tucnak/telebot.v2"
// 	tb "gopkg.in/tucnak/telebot.v2"
// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// )

// func main() {
// 	b, _ := tb.NewBot(tb.Settings{Token: "1854419870:AAHD0YMshI2FFYLoLVILMHykiWYqttU7Te8",
// 		Poller: &telebot.LongPoller{Timeout: 30 * time.Second}})

// 	// Universal markup builders.

// 	menu := &tb.ReplyMarkup{
// 		InlineKeyboard: [][]tb.InlineButton{
// 			{{
// 				Data: "btn",
// 				Text: "Hi Telebot!",
// 			}, {
// 				Data: "btn1",
// 				Text: "Hi Telebot!2",
// 			},
// 			},
// 		},
// 	}
// 	// selector = &tb.ReplyMarkup{}

// 	menu.Text("<b>bold</b>, <strong>bold</strong><i>italic</i>, <em>italic</em><u>underline</u>")

// 	// menu.Contact("Send phone number")
// 	// menu.Location("Send location")

// 	// Inline buttons.
// 	//
// 	// Pressing it will cause the client to
// 	// send the bot a callback.
// 	//
// 	// Make sure Unique stays unique as per button kind,
// 	// as it has to be for callback routing to work.
// 	//
// 	// btnPrev = selector.Data("⬅", "prev", ...)
// 	// btnNext = selector.Data("➡", "next", ...)
// 	// )

// 	// menu.Reply(
// 	// 	menu.Row(btnHelp),
// 	// 	menu.Row(btnSettings),
// 	// )
// 	// selector.Inline(
// 	// 	selector.Row(btnPrev, btnNext),
// 	// )

// 	// Command: /start <PAYLOAD>
// 	b.Handle("/start", func(m *tb.Message) {
// 		if !m.Private() {
// 			return
// 		}

// 		b.Send(m.Sender, "<b>bold</b>, <strong>bold</strong><i>italic</i>, <em>italic</em><u>underline</u>", menu, tb.ModeHTML)
// 	})

// 	// // On reply button pressed (message)
// 	// b.Handle(&btnHelp, func(m *tb.Message) {...})

// 	// // On inline button pressed (callback)
// 	// b.Handle(&btnPrev, func(c *tb.Callback) {
// 	// 	// ...
// 	// 	// Always respond!
// 	// 	b.Respond(c, &tb.CallbackResponse{...})
// 	// })

// 	b.Start()
// }

// // package main

// // import (
// // 	"log"
// // 	"os"
// // 	"os/signal"
// // 	"syscall"

// // 	"github.com/aoyako/telegram_2ch_res_bot/initialize"
// // )

// // func main() {
// // 	log.Println("Starting...")
// // 	bot, apicnt, duration := initialize.App()

// // 	quit := make(chan os.Signal)
// // 	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

// // 	go bot.Bot.Start()
// // 	go initialize.StartPolling(apicnt, duration)

// // 	log.Println("Started")

// // 	<-quit
// // 	log.Println("Quit")
// // }
