package telegram

import (
	"fmt"
	"log"
	"tdgames/logic"

	"tdgames/games"

	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/leekchan/accounting"
	telebot "gopkg.in/tucnak/telebot.v2"
)

const (
	CK_MONEY = 1000000
)

func (tb *TgBot) SendChatMessage(msg string, menu *telebot.ReplyMarkup, m *telebot.Chat) (*telebot.Message, error) {
	if menu == nil {
		return tb.Bot.Send(m, msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	} else {
		return tb.Bot.Send(m, msg, &telebot.SendOptions{ReplyMarkup: menu, ParseMode: telebot.ModeMarkdownV2})
	}

}

func (tb *TgBot) SendHtmlMessage(msg string, menu *telebot.ReplyMarkup, m *telebot.Message) (*telebot.Message, error) {
	if menu == nil {
		return tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2})
	} else {
		return tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyMarkup: menu, ParseMode: telebot.ModeMarkdownV2})
	}

}

func (tb *TgBot) DiceToMessage(m *telebot.Message, dice *telebot.Dice) (*telebot.Message, error) {

	return tb.Bot.Send(m.Chat, dice)

}

func (tb *TgBot) ReplyToMessage(msg string, m *telebot.Message) (*telebot.Message, error) {

	return tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m})

}

func (tb *TgBot) EditHtmlMessage(m *telebot.Message, msg string, menu *telebot.ReplyMarkup) (*telebot.Message, error) {
	replay := &telebot.ReplyMarkup{InlineKeyboard: m.ReplyMarkup.InlineKeyboard}
	if menu != nil {
		replay = menu
	}
	return tb.Bot.Edit(m, msg, &telebot.SendOptions{ReplyMarkup: replay, ParseMode: telebot.ModeMarkdownV2})

	//return tb.Bot.Edit(m, msg)
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

		userid := fmt.Sprintf("%d%d", m.Sender.ID, m.Chat.ID)
		targetid := fmt.Sprintf("%d%d", m.ReplyTo.Sender.ID, m.Chat.ID)

		rax, result := tb.Controller.Transfer(userid, targetid, payload)

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

		// reply := telebot.CallbackResponse{Text: "敬请期待"}
		tb.SendHtmlMessage("敬请期待", nil, m)

	}
}

// 红包
func GamesRank(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		tb.SendHtmlMessage("敬请期待", nil, m)

	}
}

// 胜场
func GamesWins(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		tb.SendHtmlMessage("敬请期待", nil, m)

	}
}

// 存款
func GamesDeposit(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		if !m.Private() {

			tb.Bot.Send(m.Chat, "存款命令只能私聊使用", &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return

		}
		payload, _ := strconv.ParseInt(m.Payload, 10, 64)
		if len(m.Payload) == 0 || payload < 1000000 {
			msg := "您的当前存款为：$0\n存款指令：\\/ck 金额\n存款金额 \\>\\= \\$100,0000"
			tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return
		}
		ac := accounting.Accounting{Symbol: "$"}
		score, derr := tb.Controller.Deposit(m.Sender.ID, payload)
		fmt.Println(score, derr)
		if derr != nil {

			msg := fmt.Sprintf("存款金额不能大于当前余额，当前余额：%s", ac.FormatMoney(score))

			tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return
		}
		msg := fmt.Sprintf("已成功存入%s，当前余额：%s", ac.FormatMoney(payload), ac.FormatMoney(score))
		tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})

	}
}

// 取钱
func GamesWithdraw(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())
		if !m.Private() {

			tb.Bot.Send(m.Chat, "取款命令只能私聊使用", &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return

		}
		payload, _ := strconv.ParseInt(m.Payload, 10, 64)
		if len(m.Payload) == 0 || payload < 1000000 {
			msg := "您的当前存款为：$0\n取款指令：/qk 金额\n金额 \\>\\= $100,0000\n每次取款手续费：0%%"
			tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return
		}

		ac := accounting.Accounting{Symbol: "$"}
		score, derr := tb.Controller.DrawMoney(m.Sender.ID, payload)
		fmt.Println(score, derr)
		if derr != nil {

			msg := fmt.Sprintf("取款金额不能大于当前余额，当前余额：%s", ac.FormatMoney(score))

			tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})
			return
		}
		msg := fmt.Sprintf("已成功取款%s，当前余额：%s", ac.FormatMoney(payload), ac.FormatMoney(score))
		tb.Bot.Send(m.Chat, msg, &telebot.SendOptions{ReplyTo: m, ParseMode: telebot.ModeMarkdownV2})

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

// 轮盘
func RouletteBet(tb *TgBot) func(m *telebot.Message) {
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

// 百家乐
func BaccaratBet(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//fmt.Println(m.MessageSig())

		start := tb.Games.NewGames(games.GAME_BACCARAT, m.Chat.ID)
		//
		if !start {
			msg := TemplateBaccarat_limit()
			tb.SendHtmlMessage(msg, nil, m)
		} else { //可以开启新局
			msg := TemplateBaccarat_Text()
			reply := TemplateBaccarat_Bet(tb)
			message, _ := tb.SendHtmlMessage(msg, reply, m)

			tb.Games.GameBegin(games.GAME_BACCARAT, message.Chat.ID, message.ID)

		}

	}
}

// fruit 水果机
func FruitBet(tb *TgBot) func(m *telebot.Message) {
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

		err := tb.Controller.Register(int64(m.Sender.ID), m.Chat.ID)
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
func KAIJU(tb *TgBot) func(m *telebot.Message) {
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

// 判断文本消息
func Ontext(tb *TgBot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		//如果是机器人自己，退出

		messageid := fmt.Sprintf("%d%d", m.Unixtime, m.ID)

		str := strings.Split(m.Text, " ")

		arrbet := games.SplitBet(str)
		if len(arrbet) <= 0 {
			return
		}
		table := tb.Games.GetTable(games.GAME_DICE, m.Chat.ID, 0)
		if table.GetStatus() > games.GS_TK_BET {
			return
			// reply := telebot.CallbackResponse{Text: "已经开局，请等待结束！", ShowAlert: true}
			// tb.Bot.Respond(c, &reply)
		}
		for i := 0; i < len(arrbet); i++ {
			intvar := arrbet[i].Score
			area := arrbet[i].Bet
			//fmt.Println(floatvar)

			player := games.PlayInfo{
				Name:   fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName),
				UserID: int64(m.Sender.ID),
			}
			fmt.Println(intvar, player)
			_, err := tb.Games.AddScore(messageid, table, player, area, intvar)
			if err != nil {
				return
				// TemplateDice_BetText()
			}
			strbets, _ := table.GetBetInfo(int64(m.Sender.ID))
			period := table.GetPeriodInfo()
			// table.GetPeriodInfo()
			balance := table.GetBalance(int64(m.Sender.ID))
			diceinfo := logic.DiceJettonInfo{
				Info:    period,
				Bets:    strbets,
				Balance: balance,
			}
			fmt.Println(diceinfo)
			strdice := TemplateDice_BetText(diceinfo)

			tb.ReplyToMessage(strdice, m)

		}
		RandDice_SendBack(tb, 0, m)

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

//随机发送骰子
func RandDice_SendBack(tb *TgBot, serialno int, m *telebot.Message) {
	dice := &telebot.Dice{
		Type: telebot.Cube.Type,
	}

	first, _ := tb.DiceToMessage(m, dice)
	if first.Dice.Type != telebot.Cube.Type {
		return
	}
	second, _ := tb.DiceToMessage(m, dice)
	if first.Dice.Type != telebot.Cube.Type {
		return
	}
	three, _ := tb.DiceToMessage(m, dice)
	if first.Dice.Type != telebot.Cube.Type {
		return
	}

	table := tb.Games.GetTable(games.GAME_DICE, int64(m.ID), serialno)
	table.SettleGame(first.Dice.Value, second.Dice.Value, three.Dice.Value)

}

// /百家乐回调下注
func Baccarat_BetCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_BACCARAT, c.Message.Chat.ID, c.Message.ID)
		if table.GetStatus() > games.GS_TK_BET {
			reply := telebot.CallbackResponse{Text: "已经开局，请等待结束！", ShowAlert: true}
			tb.Bot.Respond(c, &reply)
		}
		// floatvar, _ := strconv.ParseFloat(c.Data, 64)
		// //fmt.Println(floatvar)

		// player := games.PlayInfo{
		// 	Name:   fmt.Sprintf("%s %s", c.Sender.FirstName, c.Sender.LastName),
		// 	UserID: int64(c.Sender.ID),
		// }

		// totalscore, err := tb.Games.AddScore(c.ID, table, player, floatvar)
		// totalscore, err := 0, 0

		// if err != nil {
		// 	reply := telebot.CallbackResponse{Text: "余额不足，请通过签到获取资金后下注", ShowAlert: true}
		// 	tb.Bot.Respond(c, &reply)
		// } else {
		// 	bets, _ := tb.Games.BetInfos(c.Message.Chat.ID, c.Message.ID)
		// 	//下注成功
		// 	SendBetMessage(tb, c, totalscore)
		// 	players := TemplateBaccarat_BetText(bets)
		// 	tb.EditHtmlMessage(c.Message, players, nil)
		// }

		// tb.EditHtmlMessage(c.Message, "update text")
		//fmt.Println(score, totalscore)
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
		// floatvar, _ := strconv.ParseFloat(c.Data, 64)
		//fmt.Println(floatvar)

		// player := games.PlayInfo{
		// 	Name:   fmt.Sprintf("%s %s", c.Sender.FirstName, c.Sender.LastName),
		// 	UserID: int64(c.Sender.ID),
		// }

		// totalscore, err := tb.Games.AddScore(c.ID, table, player, floatvar)

		// if err != nil {
		// 	reply := telebot.CallbackResponse{Text: "余额不足，请通过签到获取资金后下注", ShowAlert: true}
		// 	tb.Bot.Respond(c, &reply)
		// } else {
		// 	bets, _ := tb.Games.BetInfos(c.Message.Chat.ID, c.Message.ID)
		// 	//下注成功
		// 	SendBetMessage(tb, c, totalscore)
		// 	players := TemplateNiuniu_BetText(bets)
		// 	tb.EditHtmlMessage(c.Message, players, nil)
		// }

		// tb.EditHtmlMessage(c.Message, "update text")
		//fmt.Println(score, totalscore)
	}
}

// /百家乐开始游戏
func Baccarat_StartCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_BACCARAT, c.Message.Chat.ID, c.Message.ID)
		start, err := table.StartGame(int64(c.Sender.ID))
		if !start {
			reply := telebot.CallbackResponse{Text: err.Error(), ShowAlert: true}
			tb.Bot.Respond(c, &reply)
			return
		}
		betsinfo, _ := table.GetStartInfos()
		//fmt.Println(betsinfo)

		msg := TemplateBaccarat_SelectText(betsinfo.(*logic.BaccaratSelect))
		reply := TemplateBaccarat_Select(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

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

		msg := TemplateNiuniu_SelectText(betsinfo.(*logic.NiuNiuSelect))
		reply := TemplateNiuniu_Select(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

// /签到
func Games_BalanceCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {
		// ac := accounting.Accounting{Symbol: "$"}
		// name := c.Sender.FirstName + " " + c.Sender.LastName

		// board, _ := tb.Controller.Balance(int64(c.Sender.ID), c.Message.Chat.ID)
		// str := fmt.Sprintf("%s\n\t\t当前余额:%s，本周累计胜场，财富榜：", name, ac.FormatMoney(board.Score))

		// reply := telebot.CallbackResponse{Text: str, ShowAlert: true}
		// tb.Bot.Respond(c, &reply)

		// score, err := tb.Controller.Sign(int64(c.Sender.ID), sign)

	}
}

// /签到表
func Games_SignCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		//调用方法
		n := rands(float32(1), float32(2))

		sign := 350000 + (100000 * n)

		score, err := tb.Controller.Sign(c.Sender.ID, c.Message.Chat.ID, int(sign))
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

// 选择庄闲和
func Baccarat_SelectCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_BACCARAT, c.Message.Chat.ID, c.Message.ID)
		data, _ := strconv.Atoi(c.Data)

		_, err := tb.Games.Bet(table, int64(c.Sender.ID), data)
		if err != nil {
			reply := telebot.CallbackResponse{Text: err.Error(), ShowAlert: true}
			tb.Bot.Respond(c, &reply)
			return
		}
		//fmt.Println(success, err)

		betsinfo, _ := table.GetSelectInfos()
		//fmt.Println(betsinfo)

		msg := TemplateBaccarat_SelectText(betsinfo.(*logic.BaccaratSelect))
		reply := TemplateBaccarat_Select(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

// 选择青龙白虎朱雀玄武
func Niuniu_SelectCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID)
		data, _ := strconv.Atoi(c.Data)

		_, err := tb.Games.Bet(table, int64(c.Sender.ID), data)
		if err != nil {
			reply := telebot.CallbackResponse{Text: err.Error(), ShowAlert: true}
			tb.Bot.Respond(c, &reply)
			return
		}
		//fmt.Println(success, err)

		betsinfo, _ := table.GetSelectInfos()
		//fmt.Println(betsinfo)

		msg := TemplateNiuniu_SelectText(betsinfo.(*logic.NiuNiuSelect))
		reply := TemplateNiuniu_Select(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

// 百家乐结算游戏
func Baccarat_EndGameCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		table := tb.Games.GetTable(games.GAME_BACCARAT, c.Message.Chat.ID, c.Message.ID)
		// playid := fmt.Sprintf("%d%d", c.Message.Chat.ID, c.Message.ID)

		//写分
		// _, err := table.SettleGame(int64(c.Sender.ID))
		// if err != nil {
		// 	reply := telebot.CallbackResponse{Text: err.Error(), ShowAlert: true}
		// 	tb.Bot.Respond(c, &reply)
		// 	return
		// }

		//回写数据库
		// fmt.Println(betsinfo)
		//获取游戏记录
		records, _ := table.GetSettleInfos()

		tb.Games.GameEnd(games.GAME_BACCARAT, c.Message.Chat.ID, c.Message.ID) //结算游戏

		msg := TemplateBaccarat_EndGameText(records.(*logic.BaccaratRecords))
		// // fmt.Println(msg)
		reply := TemplateNiuniu_EndGameReplyMarkUp(tb)

		tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

// 结束游戏
func Niuniu_EndGameCallBack(tb *TgBot) func(c *telebot.Callback) {
	return func(c *telebot.Callback) {

		// table := tb.Games.GetTable(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID)
		// playid := fmt.Sprintf("%d%d", c.Message.Chat.ID, c.Message.ID)

		//写分
		// _, err := table.SettleGame(int64(c.Sender.ID))
		// if err != nil {
		// 	reply := telebot.CallbackResponse{Text: err.Error(), ShowAlert: true}
		// 	tb.Bot.Respond(c, &reply)
		// 	return
		// }

		//回写数据库
		// fmt.Println(betsinfo)
		//获取游戏记录
		// records, _ := table.GetSettleInfos()
		// way, waycount := tb.Games.GetRecords(games.GAME_NIUNIU, c.Message.Chat.ID)

		// records.Ways = way
		// records.WaysCount = waycount

		// tb.Games.GameEnd(games.GAME_NIUNIU, c.Message.Chat.ID, c.Message.ID) //结算游戏

		// msg := TemplateNiuniu_EndGameText(records)
		// // fmt.Println(msg)
		// reply := TemplateNiuniu_EndGameReplyMarkUp(tb)

		// tb.EditHtmlMessage(c.Message, msg, reply)

	}
}

func SendBetMessage(tb *TgBot, c *telebot.Callback, score int64) {
	ac := accounting.Accounting{Symbol: "$"}

	str := fmt.Sprintf("下注成功\n您当前下注总额:%s\n请在无人跟注6s后点击开始游戏！", ac.FormatMoney(score))
	reply := telebot.CallbackResponse{Text: str, ShowAlert: true}
	tb.Bot.Respond(c, &reply)

}

func rands(min, max float32) float64 {
	max = max - min
	rand.Seed(time.Now().UnixNano()) //设置随机种子，使每次结果不一样
	res := Round2(float64(min+max*rand.Float32()), 2)
	fmt.Println(res)
	return res
}

func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}
