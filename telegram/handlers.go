package telegram

import (
	"fmt"
	"log"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	telebot "gopkg.in/tucnak/telebot.v2"
)

func (tb *TgBot) SendHtmlMessage(msg string, menu *telebot.ReplyMarkup, m *telebot.Message) (*telebot.Message, error) {
	return tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyMarkup: menu, ParseMode: telebot.ModeHTML})
}

func (tb *TgBot) EditHtmlMessage(m *telebot.Message, msg string) (*telebot.Message, error) {

	replay := &telebot.ReplyMarkup{InlineKeyboard: m.ReplyMarkup.InlineKeyboard}
	fmt.Println(replay)

	return tb.Bot.Edit(m, msg, &telebot.SendOptions{ReplyMarkup: replay, ParseMode: telebot.ModeHTML})

	//return tb.Bot.Edit(m, msg)
}

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
func NiuniuBet(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		start := tb.Games.GameBegin(games.GAME_NIUNIU, m.ID, m.Chat.ID)
		if start != games.GS_TK_FREE { //已经开局
			msg := TemplateNiuniu_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else {
			msg := TemplateNiuniu_Text()
			reply := TemplateNiuniu_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			fmt.Println(message.ID)

		}

	}
}

// /subs endpoint
func subs(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		// tb.Games.HandleMessage(m)
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

// /start endpoint
func Callback(tb *TgBot) func(c *telebot.Callback) {
	return func(m *telebot.Callback) {
		fmt.Println(m)
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

// /下注
func Niuniu_BetCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID)

		fmt.Println(c.MessageID, table.GetMsgID())
		// tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID)
		// m := telebot.StoredMessage{
		// 	MessageID: fmt.Sprintf("%d", c.Message.ID),
		// 	ChatID:    c.Message.Chat.ID,
		// }

		tb.EditHtmlMessage(c.Message, "update text")
		// fmt.Println(a, b)
	}
}
