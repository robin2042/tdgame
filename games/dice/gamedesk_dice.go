package dice

import (
	"fmt"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

// 大单", "小双", "大双", "小单", "小", "大", "单", "双"

// "大","小","单","双","大单","大双","小单","小双"
var (
	ID_DADAN_MARK      int = 0x10
	ID_XIAOSHUANG_MARK int = 0x13
	ID_DASHUANG_MARK   int = 0x11
	ID_XIAODAN_MARK    int = 0x12
	ID_XIAO_MARK       int = 0x01
	ID_DA_MARK         int = 0x02
	ID_DAN_MARK        int = 0x04
	ID_SHUANG_MARK     int = 0x08
)

//骰子
type Dice struct {
	games.GameDesk
}

func (g *Dice) AddScore(player games.PlayInfo, area, score int) (int64, error) {
	// return g.GameDesk.AddScore(player, score)
	return 0, nil

}
func (g *Dice) Bet(userid int64, area int) (bool, error) {
	return g.GameDesk.Bet(userid, area)
}

func (g *Dice) EndGame() error {

	g.UnInitTable()
	g.GameStation = games.GS_TK_FREE

	return nil
}

//下注信息
func (g *Dice) GetBetInfos() (bets []logic.Bets, err error) {
	return g.GameDesk.GetBetInfos()

}

//获取个人的下注信息
func (g *Dice) GetBetInfo(int64) ([]logic.Bets, error) {
	return g.GameDesk.GetBetInfos()
}

//获取ID
func (g *Dice) GetPeriodID() string {
	t1 := time.Now().Year()
	t2 := time.Now().Month()
	t3 := time.Now().Day()
	date := fmt.Sprintf("%d%02d%02d", t1, t2, t3)
	fmt.Println(date)
	isexist, _, err := g.Rdb.GetValue(date)
	fmt.Println(isexist, err)

	return ""
}
