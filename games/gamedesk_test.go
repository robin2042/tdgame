package games

import (
	"testing"
	"time"
)

func TestGameDesk_CalculateScore(t *testing.T) {

	g := &GameDesk{}

	g.InitTable("1", 1, 2)

	g.m_cbTableCardArray = [5][5]byte{{39, 25, 29, 53, 58},
		{17, 23, 5, 28, 26},

		{50, 59, 27, 41, 61},

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

	type fields struct {
		GameTable          GameTable
		MsgID              int
		PlayID             string
		ChatID             int64
		NameID             int
		GameStation        int
		LastBetTime        time.Time
		BeginTime          time.Time
		StartTime          time.Time
		NextStartTime      time.Time
		m_cbTableCardArray [5][5]byte
		Players            map[int64]PlayInfo
		Bets               map[PlayInfo]int64
		Areas              map[PlayInfo]int
		Changes            map[PlayInfo]int64
		Historys           map[PlayInfo]int64
		m_cbTimers         [5]int
		m_lUserWinScore    map[int64]int64
		m_lUserReturnScore map[int64]int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameDesk{
				GameTable:          tt.fields.GameTable,
				MsgID:              tt.fields.MsgID,
				PlayID:             tt.fields.PlayID,
				ChatID:             tt.fields.ChatID,
				NameID:             tt.fields.NameID,
				GameStation:        tt.fields.GameStation,
				LastBetTime:        tt.fields.LastBetTime,
				BeginTime:          tt.fields.BeginTime,
				StartTime:          tt.fields.StartTime,
				NextStartTime:      tt.fields.NextStartTime,
				m_cbTableCardArray: tt.fields.m_cbTableCardArray,
				Players:            tt.fields.Players,
				Bets:               tt.fields.Bets,
				Areas:              tt.fields.Areas,
				Changes:            tt.fields.Changes,
				Historys:           tt.fields.Historys,
				m_cbTimers:         tt.fields.m_cbTimers,
				m_lUserWinScore:    tt.fields.m_lUserWinScore,
				m_lUserReturnScore: tt.fields.m_lUserReturnScore,
			}
			g.CalculateScore()
		})
	}
}
