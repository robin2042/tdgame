package configs

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
	telebot "gopkg.in/tucnak/telebot.v2"
)

func TestInfoController_GetLastTimestamp(t *testing.T) {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")

	viper.ReadInConfig()
	// fmt.Println(viper.AllKeys())

	// groups := viper.Get("niuniu_start_button.bet")
	// fmt.Println(groups)
	groups1 := viper.Get("niuniu_start_button.bet")

	for _, v := range groups1.([]interface{}) {
		// fmt.Println(v)
		restlt := v.(map[interface{}]interface{})
		fmt.Println(restlt["url"])
	}

	// viper.AllKeys()
	// fmt.Println(groups1)

	// 	viper.GetUint64("disk.size")),
	// menu := &telebot.ReplyMarkup{}
	// loadMenu()

	// fmt.Println(menu)

}
func loadMenu() *telebot.ReplyMarkup {

	fmt.Println(viper.Get("niuniu_start_button.bet"))

	return nil

}
