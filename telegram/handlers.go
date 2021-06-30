package telegram

import (
	"fmt"
	"log"

	telebot "gopkg.in/tucnak/telebot.v2"
)

// /start endpoint
func start(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Register(m.Chat.ID)
		if err != nil {

		}

		// help(tb)(m)
	}
}

// /start endpoint
func NiuniuStart(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		bet := TemplateNiuniu_Bet()
		fmt.Println(bet)
	}
}

// /subs endpoint
func subs(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		// subs, err := tb.Controller.Subscription.GetSubsByChatID(m.Chat.ID)
		// if err != nil {
		// 	tb.Bot.Send(m.Sender, "Bad request")
		// 	return
		// }
		// result := fmt.Sprintf("Your subs:%s", marshallSubs(subs, true))
		// tb.Bot.Send(m.Sender, result)
	}
}

// /add bot to groups
func OnBotAddGroups(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.GroupRegister(m.Chat.ID)
		if err != nil {
			log.Println(err)
		}

	}
}

// /start endpoint
func EnterGroups(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {

		err := tb.Controller.Register(int64(m.Chat.ID))
		if err != nil {
			log.Println("插入用户失败: ", m.Chat.ID)
		}

		// help(tb)(m)
	}
}

// /start endpoint
func LeaveGroups(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Unregister(m.Chat.ID)
		if err != nil {

		}

		// help(tb)(m)
	}
}

// /救济金
func Relief(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Unregister(m.Chat.ID)
		if err != nil {

		}

		// help(tb)(m)
	}
}

// /start endpoint
func Niuniu_StartGame(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Register(m.Chat.ID)
		if err != nil {

		}

		// help(tb)(m)
	}
}

// /start endpoint
func Niuniu_Bet(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Register(m.Chat.ID)
		if err != nil {

		}

		// help(tb)(m)
	}
}

// /start endpoint
func Niuniu_EndGame(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Register(m.Chat.ID)
		if err != nil {

		}

		// help(tb)(m)
	}
}
