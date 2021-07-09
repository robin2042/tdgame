package telegram

import (
	"testing"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

func TestTemplateNiuniu_Select(t *testing.T) {
	TemplateNiuniu_Select(nil)
}

func TestTemplateNiuniu_SelectText(t *testing.T) {

	play := make([]logic.Bets, 2)

	play[0].BetArea = 0
	play[0].FmtBetArea = "青龙"

	plays := &logic.Select{
		Countdown: 60,
		Players:   play,
	}
	TemplateNiuniu_SelectText(plays)
}
