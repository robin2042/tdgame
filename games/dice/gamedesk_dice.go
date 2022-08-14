package dice

import (
	"fmt"
	"math"
	"time"

	"tdgames/games"
	"tdgames/logic"
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

	Periodinfo   logic.PeriodInfo
	WinPoint     int                       //点数
	WinArea      int                       //赢点	牌值大小单双
	WinAreaIndex int                       //赢点	牌值大小单双
	WinAreaBets  map[int64]([]games.Areas) //赢钱区域
	GameTimer    *time.Timer               //定时器
}

//获取分钟
func GetFormatHourMinute(minute, second int) string {
	t4 := time.Now().Hour() //小时

	t5 := fmt.Sprintf("%02d:%02d:%02d", t4, minute, second)

	return t5
}

//获取分钟
func GetMinute() int {
	t5 := time.Now().Minute() //分钟
	return t5
}

//获取秒
func GetSecond() int {
	t5 := time.Now().Second() //秒
	return t5
}

//设置开局信息
func (g *Dice) InitPeriodInfo() (logic.PeriodInfo, int, error) {

	g.BetMux.Lock()
	defer g.BetMux.Unlock()

	t1 := time.Now().Year()
	t2 := time.Now().Month()
	t3 := time.Now().Day()
	date := fmt.Sprintf("%d%02d%02d", t1, t2, t3)

	values, err := g.Rdb.GetValue(date)
	if err == nil {
		g.Rdb.Incr(date)
	}

	Period := fmt.Sprintf("%s%03s", date, values)

	durationsec := 1
	//开盘时间\封盘时间
	var turnontime, closetime string

	if GetMinute()%2 == 0 {
		durationsec = 2*60 - GetSecond()
		turnontime = GetFormatHourMinute(GetMinute()+2, 0)
		closetime = GetFormatHourMinute(GetMinute()+3, 50)
	} else {
		durationsec = 1*60 - GetSecond()
		turnontime = GetFormatHourMinute(GetMinute()+1, 0)
		closetime = GetFormatHourMinute(GetMinute()+2, 50)
	}

	periondInfo := logic.PeriodInfo{
		PeriodID:   Period,
		Turnontime: turnontime,
		Closetime:  closetime,
	}
	g.Periodinfo = periondInfo
	g.LastOpentime = durationsec
	g.GameTimer = time.NewTimer(time.Duration(1) * time.Second) //定时器

	return periondInfo, durationsec, nil
}

//获取期号
func (g *Dice) GetPeriodInfo() logic.PeriodInfo {
	g.BetMux.Lock()
	defer g.BetMux.Unlock()
	return g.Periodinfo
}

func (g *Dice) SetPeriodInfo(info logic.PeriodInfo) {
	g.Periodinfo = info
}

func (g *Dice) InitTable(playid string, nameid int, chatid int64) {

	g.WinAreaBets = make(map[int64][]games.Areas)

	g.GameDesk.InitTable(playid, nameid, chatid)

}

func (g *Dice) GetBetPrex() string {
	bet := fmt.Sprintf("%s_bet", g.Periodinfo.PeriodID)
	return bet
}

func (g *Dice) GetWinPrex() string {
	bet := fmt.Sprintf("%s_Win", g.Periodinfo.PeriodID)
	return bet
}

//下注
func (g *Dice) AddScore(player games.PlayInfo, area, score int) (int64, error) {

	return g.GameDesk.AddScore(player, area, score)

}

func (g *Dice) InserRedisBetList(betpre, betstring string) {
	g.Rdb.RPush(betpre, betstring)
}

//获取赢钱字符串
// 赢点吧【5586650684】双 46(1.99倍率)
func (g *Dice) GetWinString(player *games.PlayInfo, area, score int) string {
	bet := games.GetJettonStr(area)
	odds := games.GetOddsStr(area)
	return fmt.Sprintf("%s【%d】%s %d%s", player.Name, player.UserID, bet, score, odds)

}

func (g *Dice) GetBetString(player *games.PlayInfo, area, score int) string {
	bet := games.GetJettonStr(area)
	userbet := games.GetAddScoreStr(player.Name, player.UserID, bet, score)
	return userbet
}

//
func (g *Dice) Bet(userid int64, area int) (bool, error) {
	return g.GameDesk.Bet(userid, area)
}

func (g *Dice) EndGame() error {

	g.UnInitTable()
	g.GameStation = games.GS_TK_FREE

	return nil
}

//下注信息
//江湖人【5344882004】小单 500
// GetAddScoreStr
func (g *Dice) GetBetInfos() (bets []string, err error) {

	s := make([]string, 0)
	for userid, arrs := range g.Bets {
		player := g.Players[userid]
		if player == nil {
			continue
		}
		for key, value := range arrs {
			if value <= 0 {
				continue
			}
			jet := games.GetJettonStr(key)
			bet := games.GetAddScoreStr(player.Name, player.UserID, jet, int(value))

			s = append(s, bet)
		}
	}

	return s, nil

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

	values, _ := g.Rdb.GetValue(date)

	return values
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

//结算信息
// 20220814601期开奖结果
// 2+5+5=12 大双
// 赢点吧【5586650684】双 46(1.99倍率)
func (g *Dice) GetSettleInfos() (logic.Records, error) {
	winsinfo := logic.Records{}
	for userid, arrs := range g.WinAreaBets {
		if len(arrs) == 0 {
			continue
		}
		player := g.Players[userid]
		if player == nil {
			continue
		}
		for _, value := range arrs {
			winstr := g.GetWinString(player, value.Area, int(value.Score))
			fmt.Println(winstr)
		}

	}
	return winsinfo, nil
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

	g.BetMux.Lock()

	g.WinPoint = CalcPoint(first, second, three)
	g.WinArea = RetTypes(GetCardValue(g.WinPoint), GetCardTypes(g.WinPoint))
	g.WinAreaIndex = g.GetWinareaIndex(g.WinArea)

	g.BetMux.Unlock()

	g.CalculateScore()
	g.GameDesk.WriteChangeScore(g.PlayID, g.ChatID, g.M_lUserWinScore) //回写数据库

	return nil, nil
}

//停止下注
func (g *Dice) CloseGameBet() {
	values := g.GetPeriodID()
	fmt.Println(values)

}
