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

//选择按钮
func TemplateNiuniu_SelectText(plays *logic.Select) string {
	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/select.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	a := tmpl.Execute(&b, plays)
	fmt.Println(a)
	fmt.Println(b.String())
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

	a := tmpl.Execute(&b, plays)
	fmt.Println(a)
	fmt.Println(b.String())
	return b.String()
}

//选择按钮
func TemplateNiuniu_Select(tb *TgBot) *telebot.ReplyMarkup {
	menu := telebot.ReplyMarkup{}
	menu.ResizeReplyKeyboard = true
	menu.Selective = true

	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	viper.ReadInConfig()
	fmt.Println(viper.AllKeys())
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

func TemplateNiuniu_BetText(score []logic.Bets) string {

	var b bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/niuniu/bet.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	err = tmpl.Execute(&b, score)
	fmt.Println(err)
	fmt.Println(b.String())

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

//下注按钮
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
		tb.Bot.Handle(&btn, Niuniu_BalanceCallBack(tb))
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
		tb.Bot.Handle(&btn, Niuniu_SignCallBack(tb))
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
