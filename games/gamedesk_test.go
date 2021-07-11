package games

import (
	"fmt"
	"testing"
)

func TestGameDesk_CalculateScore(t *testing.T) {

	g := &GameDesk{}

	g.InitTable("1", 1, 2)
	play := PlayInfo{UserID: 12345}

	g.Bets[play] = 10000
	g.Areas[play] = 1

	g.m_cbTableCardArray = [5][5]byte{{39, 25, 29, 53, 58}, //27 19
		{17, 23, 5, 28, 26},

		{50, 59, 27, 41, 61}, //2

		{49, 54, 11, 10, 45},

		{56, 2, 53, 12, 20},
	}

	//发牌
	// g.DispatchTableCard()

	GetUnicode(g.m_cbTableCardArray[0], MAX_COUNT)
	GetUnicode(g.m_cbTableCardArray[1], MAX_COUNT)
	GetUnicode(g.m_cbTableCardArray[2], MAX_COUNT)
	GetUnicode(g.m_cbTableCardArray[3], MAX_COUNT)
	GetUnicode(g.m_cbTableCardArray[4], MAX_COUNT)

	//结算
	g.CalculateScore()

}

func TestGameDesk_GetSettleInfos(t *testing.T) {
	g := &GameDesk{}

	g.InitTable("1", 1, 2)
	play := PlayInfo{UserID: 12345}

	g.Bets[play] = 10000
	g.Areas[play] = 1

	g.m_cbTableCardArray = [5][5]byte{{39, 25, 29, 53, 58}, //27 19
		{17, 23, 5, 28, 26},

		{50, 59, 27, 41, 61}, //2

		{49, 54, 11, 10, 45},

		{56, 2, 53, 12, 20},
	}
	//records := &logic.Records{}
	detail, _ := g.GetSettleInfos()

	fmt.Println(detail)

}
