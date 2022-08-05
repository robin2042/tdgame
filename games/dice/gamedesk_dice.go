package dice

import (
	"fmt"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

// 大单", "小双", "大双", "小单", "小", "大", "单", "双"

// "大","小","单","双","大单","大双","小单","小双"

//骰子
type Dice struct {
	games.GameDesk
	WinPoint int //点数
	WinArea  int //赢点	牌值大小单双
}

func (g *Dice) AddScore(player games.PlayInfo, area, score int) (int64, error) {
	return g.GameDesk.AddScore(player, area, score)

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
func (g *Dice) GetBetInfo(userid int64) ([]string, int) {
	return g.GameDesk.GetBetInfo(userid)
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

//根据牌值类型 单双,返回大小单双
func RetTypes(value, types int) int {
	values := value & types
	return values
}

//根据牌值类型 单双
func GetCardTypes(sums int) int {
	if sums%2 == 0 {
		return games.ID_SHUANG_MARK
	} else {
		return games.ID_DAN_MARK
	}

}

//根据牌值计算大小值
func GetCardValue(sums int) int {
	if sums <= 10 {
		return games.ID_XIAO_MARK //小
	}
	if sums > 10 {
		return games.ID_DA_MARK //大
	}
	return 0
}

//计算牌点
func CalcPoint(first, second, three int) int {
	return first + second + three
}

//结算用户
func (g *Dice) CalculateScore() {
	//推算赢家
}

//回写数据库
func (g *Dice) SettleGame(first, second, three int) ([]logic.Scorelogs, error) {

	g.BetMux.Lock()
	defer g.BetMux.Unlock()
	g.WinPoint = CalcPoint(first, second, three)
	g.WinArea = RetTypes(GetCardValue(g.WinPoint), GetCardTypes(g.WinPoint))

	// var bfind bool
	// for i := range g.Bets {
	// 	if i == userid {
	// 		bfind = true
	// 		break
	// 	}
	// }
	// if !bfind {
	// 	return nil, errors.New("您没有参与此游戏，无权更改游戏状态")
	// }
	// if time.Now().Before(g.LastBetTime.Add(time.Second * 6)) {
	// 	return nil, errors.New("所有用户无操作6s后才能开始游戏")
	// }

	// ncountdown := time.Until(g.BetCountDownTime)
	// if int(ncountdown.Seconds()) > 0 {

	// }

	return nil, nil
}
