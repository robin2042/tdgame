package games

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

func TestGameDesk_CalculateScore(t *testing.T) {

	g := &GameDesk{}

	g.InitTable("1", 1, 2)

	g.m_cbTableCardArray[0] = [5]byte{9, 12, 59, 24, 25}
	g.m_cbTableCardArray[1] = [5]byte{5, 53, 60, 13, 27}
	g.m_cbTableCardArray[2] = [5]byte{39, 33, 44, 29, 18}
	g.m_cbTableCardArray[3] = [5]byte{10, 21, 4, 58, 52}
	g.m_cbTableCardArray[4] = [5]byte{20, 28, 26, 23, 40}
	g.CalculateScore()

	// play := PlayInfo{UserID: 12345}

	// g.Bets[play] = 10000
	// g.Areas[play] = 1

	// g.m_cbTableCardArray = [5][5]byte{{39, 25, 29, 53, 58}, //27 19
	// 	{17, 23, 5, 28, 26},

	// 	{50, 59, 27, 41, 61}, //2

	// 	{49, 54, 11, 10, 45},

	// 	{56, 2, 53, 12, 20},
	// }

	// //发牌
	// // g.DispatchTableCard()

	// GetUnicode(g.m_cbTableCardArray[0], MAX_COUNT)
	// GetUnicode(g.m_cbTableCardArray[1], MAX_COUNT)
	// GetUnicode(g.m_cbTableCardArray[2], MAX_COUNT)
	// GetUnicode(g.m_cbTableCardArray[3], MAX_COUNT)
	// GetUnicode(g.m_cbTableCardArray[4], MAX_COUNT)

	// //结算

}

func TestGameDesk_GetSettleInfos(t *testing.T) {
	g := &GameDesk{}

	g.InitTable("1", 1, 2)
	g.m_cbTableCardArray = [5][5]byte{{39, 25, 29, 53, 58}, //27 19
		{17, 23, 5, 28, 26},

		{50, 59, 27, 41, 61}, //2

		{49, 54, 11, 10, 45},

		{56, 2, 53, 12, 20},
	}
	record, _ := g.GetSettleInfos()
	fmt.Println(record)

}

func TestGameDesk_SettleGame(t *testing.T) {
	type fields struct {
		GameTable          GameTable
		Rdb                *storage.CloudStore
		MsgID              int
		PlayID             string
		ChatID             int64
		NameID             int
		GameStation        int
		LastBetTime        time.Time
		BetCountDownTime   time.Time
		BeginTime          time.Time
		StartTime          time.Time
		NextStartTime      time.Time
		m_cbTableCardArray [5][5]byte
		Players            map[int64]PlayInfo
		Bets               map[int64]int64
		Areas              map[int64]int
		Historys           map[PlayInfo]int64
		m_cbTimers         [5]int
		m_lUserWinScore    map[int64]int64
		m_lUserReturnScore map[int64]int64
		m_GameRecordArrary []byte
	}
	type args struct {
		userid int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []logic.Scorelogs
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GameDesk{
				GameTable:          tt.fields.GameTable,
				Rdb:                tt.fields.Rdb,
				MsgID:              tt.fields.MsgID,
				PlayID:             tt.fields.PlayID,
				ChatID:             tt.fields.ChatID,
				NameID:             tt.fields.NameID,
				GameStation:        tt.fields.GameStation,
				LastBetTime:        tt.fields.LastBetTime,
				BetCountDownTime:   tt.fields.BetCountDownTime,
				BeginTime:          tt.fields.BeginTime,
				StartTime:          tt.fields.StartTime,
				NextStartTime:      tt.fields.NextStartTime,
				m_cbTableCardArray: tt.fields.m_cbTableCardArray,
				Players:            tt.fields.Players,
				Bets:               tt.fields.Bets,
				Areas:              tt.fields.Areas,
				Historys:           tt.fields.Historys,
				m_cbTimers:         tt.fields.m_cbTimers,
				m_lUserWinScore:    tt.fields.m_lUserWinScore,
				m_lUserReturnScore: tt.fields.m_lUserReturnScore,
				m_GameRecordArrary: tt.fields.m_GameRecordArrary,
			}
			got, err := g.SettleGame(tt.args.userid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GameDesk.SettleGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameDesk.SettleGame() = %v, want %v", got, tt.want)
			}
		})
	}
}
