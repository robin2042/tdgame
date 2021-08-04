package niuniu

import (
	"fmt"
	"strconv"

	"github.com/aoyako/telegram_2ch_res_bot/games"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

func GetNiuniu_Record(s *storage.CloudStore, nameid, chatid int64) (*logic.Way, int) {

	key := fmt.Sprintf("%d%d", chatid, nameid)

	scores, err := s.LRange(key, 0, 10)
	if err != nil {
		return nil, 0
	}

	if len(scores) >= 10 {
		s.Del(key)
	}

	betinfo := &logic.Way{}
	// //天地玄黄
	for _, v := range scores {
		//fmt.Println(v)
		k, _ := strconv.Atoi(v)
		// //天地玄黄
		if (games.ID_TIAN_MARK & byte(k)) > 0 {
			betinfo.Tian += "● "
		} else {
			betinfo.Tian += "○ "
		}
		if (games.ID_DI_MARK & byte(k)) > 0 {
			betinfo.Di += "● "
		} else {
			betinfo.Di += "○ "
		}
		if (games.ID_XUAN_MARK & byte(k)) > 0 {
			betinfo.Xuan += "● "
		} else {
			betinfo.Xuan += "○ "
		}
		if (games.ID_HUANG_MARK & byte(k)) > 0 {
			betinfo.Huang += "● "
		} else {
			betinfo.Huang += "○ "
		}

	}

	// for i := 0; i < MAX_COUNT; i++ {
	// 	var str string
	// 	if i == INDEX_BANKER {
	// 		str += "🎴庄家 "
	// 		str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
	// 		str += " "
	// 		str += GetCardValueEmoj(g.m_cbTableCardArray[i])

	// 	} else if i == INDEX_PLAYER1 {
	// 		str += "🐲青龙 "
	// 		str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
	// 		str += " "
	// 		str += GetCardValueEmoj(g.m_cbTableCardArray[i])

	// 	} else if i == INDEX_PLAYER2 {
	// 		str += "🐯白虎 "
	// 		str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
	// 		str += " "
	// 		str += GetCardValueEmoj(g.m_cbTableCardArray[i])

	// 	} else if i == INDEX_PLAYER3 {
	// 		str += "🦚朱雀 "
	// 		str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
	// 		str += " "
	// 		str += GetCardValueEmoj(g.m_cbTableCardArray[i])

	// 	} else if i == INDEX_PLAYER4 {
	// 		str += "🐢玄武 "
	// 		str += GetCardTimesEmoj(g.m_cbTableCardArray[i])
	// 		str += " "
	// 		str += GetCardValueEmoj(g.m_cbTableCardArray[i])

	// 	}
	// 	betinfo.Detail = append(betinfo.Detail, str)
	// }
	// for k := range g.Players {
	// 	change := logic.ChangeScore{}
	// 	change.UserName = g.Players[k].Name

	// 	change.FmtArea = betsinfo[g.Areas[k]]

	// 	if v, ok := g.m_lUserWinScore[k]; ok {
	// 		if g.m_lUserWinScore[k] > 0 { //赢钱了

	// 			str := fmt.Sprintf("*赢* \\+%s", ac.FormatMoney(v))
	// 			change.FmtChangescore = str
	// 		} else {
	// 			str := fmt.Sprintf("*输* ~\\%s~", ac.FormatMoney(v))
	// 			change.FmtChangescore = str
	// 		}
	// 	} else {
	// 		str := fmt.Sprintf("*返回* \\+%s", ac.FormatMoney(g.m_lUserReturnScore[k]))
	// 		change.FmtChangescore = str
	// 	}

	// 	betinfo.Change = append(betinfo.Change, change)
	// }

	// betinfo.WaysCount = len(g.m_GameRecordArrary) //路子

	// fmt.Println(betinfo)
	return betinfo, len(scores)
}
