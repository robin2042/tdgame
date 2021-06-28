package telegram

import (
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
		err := tb.Controller.Register(m.Chat.ID)
		if err != nil {

		}

		// help(tb)(m)
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
