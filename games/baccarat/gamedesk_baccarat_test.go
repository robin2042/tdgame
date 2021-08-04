package baccarat

import (
	"testing"

	"github.com/aoyako/telegram_2ch_res_bot/games"
)

func TestBaccarat_DispatchTableCard(t *testing.T) {

	g := &Baccarat{}
	g.InitTable("123", 40023000, 001)
	player := &games.PlayInfo{
		Name:      "hello",
		UserID:    1213434,
		WallMoney: 1000000,
	}

	//g.Players[1213434] = player

	g.AddScore(*player, 1000)
	g.Bet(1213434, 0) //闲家

	g.DispatchTableCard()
	g.CalculateScore()

	// }
}
