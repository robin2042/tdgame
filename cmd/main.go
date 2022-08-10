package main

import (
	"os"
	"os/signal"
	"syscall"
	"tdgames/initialize"
)

var strjetton = []string{"大单", "小双", "大双", "小单", "小", "大", "单", "双"}

// const letter_goodness []string {"大","小","单","双","大单","大双","小单","小双"}

// func main() {

// 	fmt.Println(strjetton, betspeed)

// 	str := strings.Split("小单311 大单300 大双120", " ")
// 	// fmt.Println(strings.FindString("Hello World! world")) //Hello World!

// 	fmt.Println(str)

// 	for i := 0; i < len(str); i++ {
// 		for i := 0; i < len(strjetton); i++ {
// 			x := strings.Index(str[0], strjetton[i])
// 			if x >= 0 { //找到了
// 				fmt.Println(str[0])
// 				str = str[1:]
// 				fmt.Println(str)
// 			}

// 		}
// 	}
// 	// str1 := strings.Contains("小单311 大单300 大双120", "单")
// 	// fmt.Println(str1)

// 	// f := func(c rune) bool {
// 	// 	return unicode.IsNumber(c)
// 	// }

// 	// FieldsFunc() function splits the string passed
// 	// on the return values of the function f
// 	// String will therefore be split when a number
// 	// is encontered and returns all non-numbers
// 	// fmt.Printf("Fields are:%q\n",
// 	// 	strings.FieldsFunc(str[0], f))

// 	// b, err := tb.NewBot(tb.Settings{
// 	// 	Token:  "5469368758:AAGvszUEx83jSqqEhkib3mioKrsFCM0xQvA",
// 	// 	Poller: &tb.LongPoller{Timeout: 10 * time.Second},
// 	// })

// 	// if err != nil {
// 	// 	return
// 	// }

// 	// b.Handle(tb.OnText, func(m *tb.Message) {
// 	// 	b.Send(m.Chat, "hello world12345")
// 	// })

// 	// b.Start()
// }

func main() {

	bot, _ := initialize.App()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go bot.Bot.Start()
	// go initialize.StartPolling(apicnt, duration)

	<-quit

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
