package telegram

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/leekchan/accounting"

	telebot "gopkg.in/tucnak/telebot.v2"
)

func (tb *TgBot) SendHtmlMessage(msg string, menu *telebot.ReplyMarkup, m *telebot.Message) (*telebot.Message, error) {
	return tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyMarkup: menu, ParseMode: telebot.ModeMarkdownV2})
}

func (tb *TgBot) EditHtmlMessage(m *telebot.Message, msg string, menu *telebot.ReplyMarkup) (*telebot.Message, error) {
	replay := &telebot.ReplyMarkup{InlineKeyboard: m.ReplyMarkup.InlineKeyboard}
	if menu != nil {
		replay = menu
	}
	return tb.Bot.Edit(m, msg, &telebot.SendOptions{ReplyMarkup: replay, ParseMode: telebot.ModeMarkdownV2})

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

// 转账
func GamesZZ(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		if len(m.Payload) == 0 || m.ReplyTo == nil {
			msg := TemplateNiuniu_transerror()
			tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return
		}

		payload, err := strconv.ParseInt(m.Payload, 10, 64)
		if err != nil {
			msg := TemplateNiuniu_transerror()
			tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return
		}

		rax, result := tb.Controller.Transfer(int64(m.Sender.ID), int64(m.ReplyTo.Sender.ID), payload)
		if result != nil {
			tb.Bot.Send(m.Chat, "金额不足，转账失败！", &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return
		}
		//
		fmtmsg := fmt.Sprintf("转账成功，实际到账：%d，手续费：%d", payload-rax, rax)
		tb.Bot.Send(m.Chat, fmtmsg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})

	}
}

// 红包
func GamesHB(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())

		start := tb.Games.NewGames(games.GAME_NIUNIU, m.Chat.ID)
		//
		if !start {
			msg := TemplateNiuniu_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else { //可以开启新局
			msg := TemplateNiuniu_Text()
			reply := TemplateNiuniu_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			tb.Games.GameBegin(games.GAME_NIUNIU, message.Chat.ID, message.ID)

		}

	}
}

// 红包
func GamesRank(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())

		start := tb.Games.NewGames(games.GAME_NIUNIU, m.Chat.ID)
		//
		if !start {
			msg := TemplateNiuniu_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else { //可以开启新局
			msg := TemplateNiuniu_Text()
			reply := TemplateNiuniu_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			tb.Games.GameBegin(games.GAME_NIUNIU, message.Chat.ID, message.ID)

		}

	}
}

// 胜场
func GamesWins(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())

		start := tb.Games.NewGames(games.GAME_NIUNIU, m.Chat.ID)
		//
		if !start {
			msg := TemplateNiuniu_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else { //可以开启新局
			msg := TemplateNiuniu_Text()
			reply := TemplateNiuniu_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			tb.Games.GameBegin(games.GAME_NIUNIU, message.Chat.ID, message.ID)

		}

	}
}

// 存款
func GamesDeposit(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())

		start := tb.Games.NewGames(games.GAME_NIUNIU, m.Chat.ID)
		//
		if !start {
			msg := TemplateNiuniu_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else { //可以开启新局
			msg := TemplateNiuniu_Text()
			reply := TemplateNiuniu_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			tb.Games.GameBegin(games.GAME_NIUNIU, message.Chat.ID, message.ID)

		}

	}
}

// 取钱
func GamesWithdraw(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())

		start := tb.Games.NewGames(games.GAME_NIUNIU, m.Chat.ID)
		//
		if !start {
			msg := TemplateNiuniu_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else { //可以开启新局
			msg := TemplateNiuniu_Text()
			reply := TemplateNiuniu_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			tb.Games.GameBegin(games.GAME_NIUNIU, message.Chat.ID, message.ID)

		}

	}
}

// /start endpoint
func NiuniuBet(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())

		start := tb.Games.NewGames(games.GAME_NIUNIU, m.Chat.ID)
		//
		if !start {
			msg := TemplateNiuniu_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else { //可以开启新局
			msg := TemplateNiuniu_Text()
			reply := TemplateNiuniu_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			tb.Games.GameBegin(games.GAME_NIUNIU, message.Chat.ID, message.ID)

		}

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

		err := tb.Controller.Register(int64(m.Sender.ID))
		if err != nil {
			log.Println("插入用户失败: ", m.Chat.ID)
		}

		// help(tb)(m)
	}
}

// /start endpoint
func LeaveGroups(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		err := tb.Controller.Unregister(int64(m.Sender.ID))
		if err != nil {

		}

		// help(tb)(m)
	}
}

// /start endpoint
func Callback(tb *TgBot) func(c *telebot.Callback) {
	return func(m *telebot.Callback) {
		//fmt.Println(m)
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

		table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID)
		if table.GetStatus() > games.GS_TK_BET {
			reply := telebot.CallbackResponse{Text: "已经开局，请等待结束！", ShowAlert: true}
			tb.Bot.Respond(c, &reply)
		}
		floatvar, _ := strconv.ParseFloat(c.Data, 64)
		//fmt.Println(floatvar)

		player := games.PlayInfo{
			Name:   fmt.Sprintf("%s %s", c.Sender.FirstName, c.Sender.LastName),
			UserID: int64(c.Sender.ID),
		}

		totalscore, err := tb.Games.AddScore(table, player, floatvar)

		if err != nil {
			reply := telebot.CallbackResponse{Text: "余额不足，请通过签到获取资金后下注", ShowAlert: true}
			tb.Bot.Respond(c, &reply)
		} else {
			bets, _ := tb.Games.BetInfos(c.Message.Chat.ID, c.Message.ID)
			//下注成功
			SendBetMessage(tb, c, totalscore)
			players := TemplateNiuniu_BetText(bets)
			tb.EditHtmlMessage(c.Message, players, nil)
		}

		// tb.EditHtmlMessage(c.Message, "update text")
		//fmt.Println(score, totalscore)
	}
}

// /开始游戏
func Niuniu_StartCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID)
		start, err := table.StartGame(int64(c.Sender.ID))
		if !start {
			reply := telebot.CallbackResponse{Text: err.Error(), ShowAlert: true}
			tb.Bot.Respond(c, &reply)
			return
		}
		betsinfo, _ := table.GetStartInfos()
		//fmt.Println(betsinfo)

		msg := TemplateNiuniu_SelectText(betsinfo)
		reply := TemplateNiuniu_Select(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

// /签到
func Niuniu_BalanceCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		ac := accounting.Accounting{Symbol: "$"}
		name := c.Sender.FirstName + " " + c.Sender.LastName

		board, _ := tb.Controller.Balance(int64(c.Sender.ID))
		str := fmt.Sprintf("%s\n\t\t当前余额:%s", name, ac.FormatMoney(board.Score))

		reply := telebot.CallbackResponse{Text: str, ShowAlert: true}
		tb.Bot.Respond(c, &reply)

		// score, err := tb.Controller.Sign(int64(c.Sender.ID), sign)

	}
}

// /签到表
func Niuniu_SignCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		sign := 700000

		score, err := tb.Controller.Sign(c.Sender.ID, c.Message.Chat.ID, sign)
		if !err {
			reply := telebot.CallbackResponse{Text: "150秒内限定签到一次", ShowAlert: true}
			tb.Bot.Respond(c, &reply)
		} else {
			ac := accounting.Accounting{Symbol: "$"}

			str := fmt.Sprintf("签到成功\n\t\t系统赠送了您:%s\n\t\t当前总余额:%s\n\t\t每间隔150秒可再次点击签到领取", ac.FormatMoney(sign), ac.FormatMoney(score))
			reply := telebot.CallbackResponse{Text: str, ShowAlert: true}
			tb.Bot.Respond(c, &reply)

		}

		// table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID)
		// floatvar, _ := strconv.ParseFloat(c.Data, 64)
		// fmt.Println(floatvar)

		// score, err := tb.Controller.Register()
		// if err != nil {
		// 	reply := telebot.CallbackResponse{Text: "余额不足，请通过签到获取资金后下注", ShowAlert: true}
		// 	tb.Bot.Respond(c, &reply)
		// } else {
		// 	fmt.Println(score)
		// }

		// tb.EditHtmlMessage(c.Message, "update text")
		// fmt.Println(a, b)
	}
}

// 选择青龙白虎朱雀玄武
func Niuniu_SelectCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID)
		data, _ := strconv.Atoi(c.Data)

		tb.Games.Bet(table, int64(c.Sender.ID), data)
		//fmt.Println(success, err)

		betsinfo, _ := table.GetSelectInfos()
		//fmt.Println(betsinfo)

		msg := TemplateNiuniu_SelectText(betsinfo)
		reply := TemplateNiuniu_Select(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

// 结束游戏
func Niuniu_EndGameCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID)
		// playid := fmt.Sprintf("%d%d", c.Message.Chat.ID, c.Message.ID)

		//写分
		_, err := table.SettleGame(int64(c.Sender.ID))
		if err != nil {
			reply := telebot.CallbackResponse{Text: err.Error(), ShowAlert: true}
			tb.Bot.Respond(c, &reply)
			return
		}

		//回写数据库
		// fmt.Println(betsinfo)
		//获取游戏记录
		records, _ := table.GetSettleInfos()
		way, waycount := tb.Games.GetRecords(games.GAME_NIUNIU, c.Message.Chat.ID)

		records.Ways = way
		records.WaysCount = waycount

		tb.Games.GameEnd(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID) //结算游戏

		msg := TemplateNiuniu_EndGameText(records)
		// fmt.Println(msg)
		reply := TemplateNiuniu_EndGameReplyMarkUp(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

func SendBetMessage(tb *TgBot, c *telebot.Callback, score int64) {
	ac := accounting.Accounting{Symbol: "$"}

	str := fmt.Sprintf("下注成功\n您当前下注总额:%s\n请在无人跟注后点击开始游戏！", ac.FormatMoney(score))
	reply := telebot.CallbackResponse{Text: str, ShowAlert: true}
	tb.Bot.Respond(c, &reply)

}
