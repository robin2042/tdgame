package telegram

import (
	"bytes"
	"fmt"
	"text/template"

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

func TemplateNiuniu_BetText() string {
	sweaters := []UserInfo{
		{
			Name:   "tom",
			Gender: "男人",
			Age:    1,
		},
		{
			Name:   "john",
			Gender: "女人人",
			Age:    1,
		},
	}
	fmt.Println(sweaters)

	var b bytes.Buffer
	tmpl, err := template.ParseFiles("hello.tmpl")
	if err != nil {
		fmt.Println("create template failed,err:", err)
		return "无效"
	}

	err = tmpl.Execute(&b, sweaters)
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

	// btnarray := make([][]telebot.InlineButton, 0)
	jettons := TemplateNiuniu_Jetton(tb, viper.GetViper())
	starts := TemplateNiuniu_Start(tb, viper.GetViper()) //开始
	btnarray := updateSlice(jettons, starts)

	balance := TemplateNiuniu_Balance(tb, viper.GetViper()) //余额
	btnarray = updateSlice(btnarray, balance)
	signs := TemplateNiuniu_Sign(tb, viper.GetViper()) //签到

	btnarray = updateSlice(btnarray, signs)

	// TemplateNiuniu_Jetton(tb, viper.GetViper(), &btnarray)

	// // groups := viper.Get("niuniu_start_button.bet")
	// // fmt.Println(groups)
	// groups1 := viper.Get("niuniu_start_button.bet")

	// a := groups1.([]interface{})

	// for _, row := range a {
	// 	keys := make([]telebot.InlineButton, 0, len(row.([]interface{})))
	// 	fmt.Println(row)
	// 	for _, v := range row.([]interface{}) {
	// 		var btn telebot.InlineButton

	// 		restlt := v.(map[interface{}]interface{})
	// 		btn.Text = restlt["text"].(string)
	// 		switch restlt["data"].(type) {
	// 		case float64:
	// 			btn.Data = fmt.Sprintf("%f", restlt["data"].(float64))
	// 		case int:
	// 			btn.Data = fmt.Sprintf("%d", restlt["data"].(int))
	// 		}

	// 		btn.Unique = restlt["unique"].(string)
	// 		tb.Bot.Handle(&btn, Niuniu_BetCallBack(tb))
	// 		keys = append(keys, btn)
	// 	}
	// 	btnarray = append(btnarray, keys)
	// }

	menu.InlineKeyboard = btnarray
	return &menu
}

//下注按钮
func TemplateNiuniu_Jetton(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("niuniu_start_button.bet")

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
func TemplateNiuniu_Start(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("niuniu_start_button.start")

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
			tb.Bot.Handle(&btn, Niuniu_StartCallBack(tb))
			keys = append(keys, btn)
		}
		btnarray = append(btnarray, keys)
	}
	return btnarray
}

//余额
func TemplateNiuniu_Balance(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("niuniu_start_button.balance")

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
			tb.Bot.Handle(&btn, Niuniu_StartCallBack(tb))
			keys = append(keys, btn)
		}
		btnarray = append(btnarray, keys)
	}
	return btnarray
}

//签到
func TemplateNiuniu_Sign(tb *TgBot, viper *viper.Viper) [][]telebot.InlineButton {
	jetton := viper.Get("niuniu_start_button.sign")

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
			tb.Bot.Handle(&btn, Niuniu_StartCallBack(tb))
			keys = append(keys, btn)
		}
		btnarray = append(btnarray, keys)
	}
	return btnarray
}

func updateSlice(first [][]telebot.InlineButton, second [][]telebot.InlineButton) [][]telebot.InlineButton {
	btnarray := make([][]telebot.InlineButton, 0)

	for _, row := range first {
		btnarray = append(btnarray, row)
	}
	for _, row := range second {
		btnarray = append(btnarray, row)
	}
	return btnarray
}
