package games

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
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

var strjetton = []string{"大单", "小双", "大双", "小单", "小", "大", "单", "双"}
var jetmark = []int{ID_DADAN_MARK, ID_XIAOSHUANG_MARK,
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

func SplitBet(str []string) []logic.DiceBetInfo {
	// var arbet []logic.DiceBetInfo = make([]logic.DiceBetInfo, 0) //最大同时20格下注
	// return arbet
	// 	fmt.Println(len(str))
	var arbet []logic.DiceBetInfo = make([]logic.DiceBetInfo, 0) //最大同时20格下注
	for len(str) > 0 {
		for x := 0; x < len(strjetton); x++ {
			find := strings.Index(str[0], strjetton[x])
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
