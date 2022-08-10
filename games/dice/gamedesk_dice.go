package dice

import (
	"fmt"
	"math"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

// 大单", "小双", "大双", "小单", "小", "大", "单", "双"

// "大","小","单","双","大单","大双","小单","小双"

func Change(slice []int64) []int64 {
	slice = append(slice, 1)
	return slice
}

//结果
type History struct {
	Last    string   `json:"last"`
	History []string `json:"history"`
}

//结果
type Lottery struct {
	PeriodID string   `json:"periodID"`
	Last     string   `json:"last"`
	Users    []string `json:"users"`
}

//骰子
type Dice struct {
	games.GameDesk
	PeriodID     string                    //第几期
	WinPoint     int                       //点数
	WinArea      int                       //赢点	牌值大小单双
	WinAreaIndex int                       //赢点	牌值大小单双
	WinAreaBets  map[int64]([]games.Areas) //赢钱区域
}

func (g *Dice) InitTable(playid string, nameid int, chatid int64) {

	g.WinAreaBets = make(map[int64][]games.Areas)

	g.GameDesk.InitTable(playid, nameid, chatid)

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
	values := value | types
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

//获取值
func (g *Dice) GetWinareaIndex(winarea int) int {
	for key, value := range games.JET_MARK {
		if value == winarea {
			return key
		}
	}
	return -1
}

//结算用户
//根据结果,比对是否选中.
func (g *Dice) CalculateScore() {
	g.BetMux.Lock()
	for userid, arrs := range g.Bets {
		for key, value := range arrs {
			if value <= 0 {
				continue
			} else {
				_, v := g.WinAreaBets[userid]
				if !v {
					g.WinAreaBets[userid] = make([]games.Areas, 0)
				}
			}
			if g.WinArea&games.JET_MARK[key] != 0 { //中奖了

				//地板取整
				score := int64(math.Floor(games.Bet_SPEED[key] * float64(value)))
				fmt.Println(score)
				wins := g.WinAreaBets[userid]
				area := games.Areas{
					Area:  key,
					Score: score,
				}
				wins = append(wins, area)
				g.WinAreaBets[userid] = wins
				g.M_lUserWinScore[userid] += score //赢钱累加
			}

		}
	}
	defer g.BetMux.Unlock()
}

//回写数据库
func (g *Dice) SettleGame(first, second, three int) ([]logic.Scorelogs, error) {
	fmt.Println(first, second, three)

	g.BetMux.Lock()

	g.WinPoint = CalcPoint(first, second, three)
	g.WinArea = RetTypes(GetCardValue(g.WinPoint), GetCardTypes(g.WinPoint))
	g.WinAreaIndex = g.GetWinareaIndex(g.WinArea)

	g.BetMux.Unlock()

	g.CalculateScore()
	g.GameDesk.WriteChangeScore(g.PlayID, g.ChatID, g.M_lUserWinScore) //回写数据库

	return nil, nil
}
