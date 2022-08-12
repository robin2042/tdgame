package games

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"tdgames/logic"
)

var (
	ID_DADAN_MARK      int = 0x18
	ID_XIAOSHUANG_MARK int = 0x22
	ID_DASHUANG_MARK   int = 0x28
	ID_XIAODAN_MARK    int = 0x12
	ID_XIAO_MARK       int = 0x02
	ID_DA_MARK         int = 0x08
	ID_DAN_MARK        int = 0x10
	ID_SHUANG_MARK     int = 0x20
)

var JETTON_STR = []string{"大单", "小双", "大双", "小单", "小", "大", "单", "双"}
var JET_MARK = []int{ID_DADAN_MARK, ID_XIAOSHUANG_MARK,
	ID_DASHUANG_MARK,
	ID_XIAODAN_MARK,
	ID_XIAO_MARK,
	ID_DA_MARK,
	ID_DAN_MARK,
	ID_SHUANG_MARK}

var override = []string{"（赔率3.4倍）", "（赔率4.4倍）",
	"（赔率4.4倍）",
	"（赔率4.4倍）",
	"（赔率1.99倍）",
	"（赔率1.99倍）",
	"（赔率1.99倍）",
	"（赔率1.99倍）"}

//中奖倍率
var Bet_SPEED = []float64{3.4, 3.4, 4.4, 4.4, 1.99, 1.99, 1.99, 1.99}

//赢钱区域
type Areas struct {
	Area  int
	Score int64
}

// 日子富裕【5516166760】小 50
func GetAddScoreStr(username string, userid int64, n string, score int) string {
	return fmt.Sprintf("%s【%d】%s %d", username, userid, n, score)
}

//获取中将中文
func GetJettonStr(n int) string {
	return JETTON_STR[n]
}

func SplitBet(str []string) []logic.DiceBetInfo {
	// var arbet []logic.DiceBetInfo = make([]logic.DiceBetInfo, 0) //最大同时20格下注
	// return arbet
	// 	fmt.Println(len(str))
	var arbet []logic.DiceBetInfo = make([]logic.DiceBetInfo, 0) //最大同时20格下注
	for len(str) > 0 {
		for x := 0; x < len(JETTON_STR); x++ {
			find := strings.Index(str[0], JETTON_STR[x])
			if find >= 0 { //找到了
				re := regexp.MustCompile("[0-9]+")
				strbet := re.FindString(str[0])
				if strbet == "" {
					break
				}
				score, err := strconv.Atoi(strbet)
				if err != nil {
					break
				}
				var bet logic.DiceBetInfo
				bet.Bet = x
				bet.Score = score
				arbet = append(arbet, bet)
				fmt.Println(str[0])
				str = str[1:]
				if len(str) == 0 {
					break
				}
				fmt.Println(len(str))
				break
			}

		}
	}
	return arbet
}
