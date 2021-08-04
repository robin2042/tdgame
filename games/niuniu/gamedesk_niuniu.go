package niuniu

import (
	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

//百家乐
type Niuniu struct {
	games.GameDesk
}

func (g *Niuniu) AddScore(player games.PlayInfo, score float64) (int64, error) {
	return g.GameDesk.AddScore(player, score)

}
func (g *Niuniu) Bet(userid int64, area int) (bool, error) {
	return g.GameDesk.Bet(userid, area)
}

func (g *Niuniu) EndGame() error {

	g.UnInitTable()
	g.GameStation = games.GS_TK_FREE

	return nil
}

//下注信息
func (g *Niuniu) GetBetInfos() (bets []logic.Bets, err error) {
	return g.GameDesk.GetBetInfos()

}
