package telegram

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/spf13/viper"
	telebot "gopkg.in/tucnak/telebot.v2"
)

// main.go
type UserInfo struct {
	Name   string
	Gender string
	Age    int
}

func TemplateBaccarat_Text() string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/Baccarat/start.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, nil)
	return b.String()
}

func TemplateNiuniu_Text() string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/start.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, nil)
	return b.String()
}

//百家乐选择按钮
func TemplateBaccarat_SelectText(plays *logic.BaccaratSelect) string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/Baccarat/select.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, plays)

	return b.String()
}

//选择按钮
func TemplateNiuniu_SelectText(plays *logic.NiuNiuSelect) string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/select.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, plays)

	return b.String()
}

//百家乐按钮
func TemplateBaccarat_EndGameText(plays *logic.BaccaratRecords) string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/Baccarat/endgame.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, plays)

	return b.String()
}

//选择按钮
func TemplateNiuniu_EndGameText(plays *logic.Records) string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/endgame.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, plays)

	return b.String()
}

//百家乐选择按钮
func TemplateBaccarat_Select(tb *TgBot) *telebot.ReplyMarkup {
	menu := telebot.ReplyMarkup{}
	menu.ResizeReplyKeyboard = true
	menu.Selective = true

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	viper.ReadInConfig()

	btnarray := make([][]telebot.InlineButton, 0)
	jettons := TemplateBaccarat_SelectJetton(tb, viper.GetViper())

	btnarray = append(btnarray, jettons...)

	buttons := make([]telebot.InlineButton, 0)
	start := TemplateBaccarat_SettlementButton(tb, viper.GetViper()) //余额，签到

	buttons = append(buttons, start)

	btnarray = append(btnarray, buttons)

	menu.InlineKeyboard = btnarray
	return &menu
}

//选择按钮
func TemplateNiuniu_Select(tb *TgBot) *telebot.ReplyMarkup {
	menu := telebot.ReplyMarkup{}
	menu.ResizeReplyKeyboard = true
	menu.Selective = true

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	viper.ReadInConfig()

	btnarray := make([][]telebot.InlineButton, 0)
	jettons := TemplateNiuniu_SelectJetton(tb, viper.GetViper())

	btnarray = append(btnarray, jettons...)

	buttons := make([]telebot.InlineButton, 0)
	start := TemplateNiuniu_SettlementButton(tb, viper.GetViper()) //开始

	buttons = append(buttons, start)

	btnarray = append(btnarray, buttons)

	menu.InlineKeyboard = btnarray
	return &menu
}

//百家乐下注文本
func TemplateBaccarat_BetText(score []logic.Bets) string {

	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/Baccarat/bet.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, score)

	return b.String()
}

func TemplateNiuniu_BetText(score []logic.Bets) string {

	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/bet.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, score)

	return b.String()
}

func TemplateNiuniu_limit() string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/block.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, nil)
	return b.String()
}

func TemplateBaccarat_limit() string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/block.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, nil)
	return b.String()
}

//百家乐下注
func TemplateBaccarat_Bet(tb *TgBot) *telebot.ReplyMarkup {
	menu := telebot.ReplyMarkup{}
	menu.ResizeReplyKeyboard = true
	menu.Selective = true

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	viper.ReadInConfig()
	fmt.Println(viper.AllKeys())
	btnarray := make([][]telebot.InlineButton, 0)
	jettons := TemplateBaccarat_Jetton(tb, viper.GetViper())

	btnarray = append(btnarray, jettons...)

	starts := make([][]telebot.InlineButton, 0)

	// buttons := make([]telebot.InlineButton, 0)
	// start := TemplateBaccarat_Start(tb, viper.GetViper()) //开始
	// buttons = append(buttons, start)
	// balance := TemplateNiuniu_Balance(tb, viper.GetViper()) //余额
	// buttons = append(buttons, balance)
	// sign := TemplateNiuniu_Sign(tb, viper.GetViper()) //签到
	// buttons = append(buttons, sign)
	// starts = append(starts, balance...)

	// starts = append(starts, buttons)

	btnarray = append(btnarray, starts...)

	menu.InlineKeyboard = btnarray
	return &menu
}

//牛牛下注
func TemplateNiuniu_Bet(tb *TgBot) *telebot.ReplyMarkup {
	menu := telebot.ReplyMarkup{}
	menu.ResizeReplyKeyboard = true
	menu.Selective = true

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	viper.ReadInConfig()
	fmt.Println(viper.AllKeys())
	btnarray := make([][]telebot.InlineButton, 0)
	jettons := TemplateNiuniu_Jetton(tb, viper.GetViper())

	btnarray = append(btnarray, jettons...)

	starts := make([][]telebot.InlineButton, 0)

	buttons := make([]telebot.InlineButton, 0)
	start := TemplateNiuniu_Start(tb, viper.GetViper()) //开始
	buttons = append(buttons, start)
	balance := TemplateNiuniu_Balance(tb, viper.GetViper()) //余额
	buttons = append(buttons, balance)
	sign := TemplateNiuniu_Sign(tb, viper.GetViper()) //签到
	buttons = append(buttons, sign)
	// starts = append(starts, balance...)

	starts = append(starts, buttons)

	btnarray = append(btnarray, starts...)

	menu.InlineKeyboard = btnarray
	return &menu
}

//余额
func TemplateBaccarat_SettlementButton(tb *TgBot, viper *viper.Viper) telebot.InlineButton {
	jetton := viper.Get("niuniu_start_button.settle")

	arr := jetton.([]interface{})
	// btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		// keys := make(telebot.InlineButton, 0)
		// for _, v := range row.([]interface{}) {
		var btn telebot.InlineButton

		restlt := row.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Baccarat_EndGameCallBack(tb))
		return btn

	}

	return telebot.InlineButton{}
}

//余额
func TemplateNiuniu_SettlementButton(tb *TgBot, viper *viper.Viper) telebot.InlineButton {
	jetton := viper.Get("niuniu_start_button.settle")

	arr := jetton.([]interface{})
	// btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		// keys := make(telebot.InlineButton, 0)
		// for _, v := range row.([]interface{}) {
		var btn telebot.InlineButton

		restlt := row.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Niuniu_EndGameCallBack(tb))
		return btn

		// keys = append(keys, btn)
		// btnarray = append(btnarray, keys)
	}

	return telebot.InlineButton{}
}

//百家乐选择区域
func TemplateBaccarat_SelectJetton(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("baccarat_start_button.select")
	arr := jetton.([]interface{})
	btnarray := make([][]telebot.InlineButton, 0)

	for _, row := range arr {
		keys := make([]telebot.InlineButton, 0, len(row.([]interface{})))
		for _, v := range row.([]interface{}) {
			var btn telebot.InlineButton

			restlt := v.(map[interface{}]interface{})

			btn.Text = restlt["text"].(string)
			switch restlt["data"].(type) {
			case float64:
				btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
			case int:
				btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
			}

			btn.Unique = restlt["unique"].(string)
			tb.Bot.Handle(&btn, Baccarat_SelectCallBack(tb))
			keys = append(keys, btn)
		}
		btnarray = append(btnarray, keys)
	}

	// }
	return btnarray
}

//青龙白虎
func TemplateNiuniu_SelectJetton(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("niuniu_start_button.select")
	arr := jetton.([]interface{})
	btnarray := make([][]telebot.InlineButton, 0)
	// for _, row := range arr {
	keys := make([]telebot.InlineButton, 0)
	for _, v := range arr {
		var btn telebot.InlineButton

		restlt := v.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Niuniu_SelectCallBack(tb))
		keys = append(keys, btn)
	}
	btnarray = append(btnarray, keys)
	// }
	return btnarray
}

//百家乐按钮
func TemplateBaccarat_Jetton(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("baccarat_jetton_button.bet")

	arr := jetton.([]interface{})
	btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		keys := make([]telebot.InlineButton, 0, len(row.([]interface{})))
		for _, v := range row.([]interface{}) {
			var btn telebot.InlineButton

			restlt := v.(map[interface{}]interface{})
			btn.Text = restlt["text"].(string)
			switch restlt["data"].(type) {
			case float64:
				btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
			case int:
				btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
			}

			btn.Unique = restlt["unique"].(string)
			tb.Bot.Handle(&btn, Baccarat_BetCallBack(tb))
			keys = append(keys, btn)
		}
		btnarray = append(btnarray, keys)
	}
	return btnarray
}

//牛牛按钮
func TemplateNiuniu_Jetton(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("niuniu_jetton_button.bet")

	arr := jetton.([]interface{})
	btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		keys := make([]telebot.InlineButton, 0, len(row.([]interface{})))
		for _, v := range row.([]interface{}) {
			var btn telebot.InlineButton

			restlt := v.(map[interface{}]interface{})
			btn.Text = restlt["text"].(string)
			switch restlt["data"].(type) {
			case float64:
				btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
			case int:
				btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
			}

			btn.Unique = restlt["unique"].(string)
			tb.Bot.Handle(&btn, Niuniu_BetCallBack(tb))
			keys = append(keys, btn)
		}
		btnarray = append(btnarray, keys)
	}
	return btnarray
}

//百家乐开始
func TemplateBaccarat_Start(tb *TgBot, viper *viper.Viper) telebot.InlineButton {
	jetton := viper.Get("baccarat_jetton_button.start")

	restlt := jetton.([]interface{})
	fmt.Println(restlt)

	// btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range restlt {
		// keys := make([]telebot.InlineButton, 0)
		var btn telebot.InlineButton

		restlt := row.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Baccarat_StartCallBack(tb)) //百家乐开始游戏
		// keys = append(keys, btn)
		return btn
	}

	// }
	return telebot.InlineButton{}
}

//开始
func TemplateNiuniu_Start(tb *TgBot, viper *viper.Viper) telebot.InlineButton {
	jetton := viper.Get("niuniu_jetton_button.start")

	restlt := jetton.([]interface{})
	fmt.Println(restlt)

	// btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range restlt {
		// keys := make([]telebot.InlineButton, 0)
		var btn telebot.InlineButton

		restlt := row.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Niuniu_StartCallBack(tb))
		// keys = append(keys, btn)
		return btn
	}

	// }
	return telebot.InlineButton{}
}

//百家乐余额
func TemplateBaccarat_Balance(tb *TgBot, viper *viper.Viper) telebot.InlineButton {
	jetton := viper.Get("baccarat_jetton_button.balance")

	arr := jetton.([]interface{})
	// btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		// keys := make(telebot.InlineButton, 0)
		// for _, v := range row.([]interface{}) {
		var btn telebot.InlineButton

		restlt := row.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Games_BalanceCallBack(tb))
		return btn

		// keys = append(keys, btn)
		// btnarray = append(btnarray, keys)
	}

	return telebot.InlineButton{}
}

//余额
func TemplateNiuniu_Balance(tb *TgBot, viper *viper.Viper) telebot.InlineButton {
	jetton := viper.Get("niuniu_jetton_button.balance")

	arr := jetton.([]interface{})
	// btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		// keys := make(telebot.InlineButton, 0)
		// for _, v := range row.([]interface{}) {
		var btn telebot.InlineButton

		restlt := row.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Games_BalanceCallBack(tb))
		return btn

		// keys = append(keys, btn)
		// btnarray = append(btnarray, keys)
	}

	return telebot.InlineButton{}
}

//签到
func TemplateNiuniu_Sign(tb *TgBot, viper *viper.Viper) telebot.InlineButton {
	jetton := viper.Get("niuniu_jetton_button.sign")

	arr := jetton.([]interface{})
	// btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		// keys := make([]telebot.InlineButton, 0)
		// for _, v := range row.([]interface{}) {
		var btn telebot.InlineButton

		restlt := row.(map[interface{}]interface{})
		btn.Text = restlt["text"].(string)
		switch restlt["data"].(type) {
		case float64:
			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
		case int:
			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
		}

		btn.Unique = restlt["unique"].(string)
		tb.Bot.Handle(&btn, Games_SignCallBack(tb))
		return btn

	}
	// btnarray = append(btnarray, keys)
	// }
	return telebot.InlineButton{}
}

//青龙白虎
func TemplateNiuniu_EndGameReplyMarkUp(tb *TgBot) *telebot.ReplyMarkup {

	menu := telebot.ReplyMarkup{}
	menu.ResizeReplyKeyboard = true
	menu.Selective = true

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	viper.ReadInConfig()

	btnarray := TemplateNiuniu_EndGameJetton(tb, viper.GetViper())

	fmt.Println(btnarray)
	menu.InlineKeyboard = btnarray
	return &menu
}

//青龙白虎
func TemplateNiuniu_EndGameJetton(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {

	btnarray := TemplateNiuniu_BuildJetton(tb)
	fmt.Println(btnarray)
	buttons := make([]telebot.InlineButton, 0)

	balance := TemplateNiuniu_Balance(tb, viper) //余额
	buttons = append(buttons, balance)
	sign := TemplateNiuniu_Sign(tb, viper) //签到
	buttons = append(buttons, sign)

	btnarray = append(btnarray, buttons)
	return btnarray
}

func TemplateNiuniu_BuildJetton(tb *TgBot) [][]telebot.InlineButton {

	btnarray := make([][]telebot.InlineButton, 0)
	return btnarray
}

func TemplateNiuniu_transerror() string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/transerror.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, nil)
	return b.String()
}

//骰子按钮
func TemplateDice_Jetton(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("dice_jetton_button.bet")

	arr := jetton.([]interface{})
	btnarray := make([][]telebot.InlineButton, 0)
	for _, row := range arr {
		keys := make([]telebot.InlineButton, 0, len(row.([]interface{})))
		for _, v := range row.([]interface{}) {
			var btn telebot.InlineButton

			restlt := v.(map[interface{}]interface{})
			btn.Text = restlt["text"].(string)
			switch restlt["data"].(type) {
			case float64:
				btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
			case int:
				btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
			}

			btn.Unique = restlt["unique"].(string)
			tb.Bot.Handle(&btn, Niuniu_BetCallBack(tb))
			keys = append(keys, btn)
		}
		btnarray = append(btnarray, keys)
	}
	return btnarray
}

func TemplateDice_Text(period logic.PeriodInfo) string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/dice/start.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, period)
	return b.String()
}

//骰子下注按钮
func TemplateDice_Bet(tb *TgBot) *telebot.ReplyMarkup {
	menu := telebot.ReplyMarkup{}
	menu.ResizeReplyKeyboard = true
	menu.Selective = true

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	viper.ReadInConfig()
	fmt.Println(viper.AllKeys())
	btnarray := make([][]telebot.InlineButton, 0)
	jettons := TemplateDice_Jetton(tb, viper.GetViper())

	btnarray = append(btnarray, jettons...)

	starts := make([][]telebot.InlineButton, 0)

	// buttons := make([]telebot.InlineButton, 0)
	// start := TemplateNiuniu_Start(tb, viper.GetViper()) //开始
	// buttons = append(buttons, start)
	// balance := TemplateNiuniu_Balance(tb, viper.GetViper()) //余额
	// buttons = append(buttons, balance)
	// sign := TemplateNiuniu_Sign(tb, viper.GetViper()) //签到
	// buttons = append(buttons, sign)
	// // starts = append(starts, balance...)

	// starts = append(starts, buttons)

	btnarray = append(btnarray, starts...)

	menu.InlineKeyboard = btnarray
	return &menu
}

//骰子下注信息
func TemplateDice_BetText(score logic.DiceJettonInfo) string {

	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/dice/bet.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	tmpl.Execute(&b, score)

	return b.String()
}
