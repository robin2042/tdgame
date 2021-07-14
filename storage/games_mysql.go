package storage

import (
	"errors"

	logger "github.com/aoyako/telegram_2ch_res_bot/logger"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// UserPostgres is an implementation of storage.User
type GamesMysql struct {
	db *gorm.DB
}

// NewUserPostgres constructor of UserPostgres struct
func NewGamesMysql(db *gorm.DB) *GamesMysql {
	return &GamesMysql{
		db: db,
	}
}

// Register adds user in databse
func (groupStorage *GamesMysql) SaveGameRound(game *logic.Gamerounds) error {

	result := groupStorage.db.Create(game)
	return result.Error
}

//下注金额，携带金额，错误
func (groupStorage *GamesMysql) AddScore(addscore *logic.AddScore) (int64, int64, error) {

	var user logic.User
	// var score int64
	var floatscore float64

	groupStorage.db.Where("userid = ?", addscore.Userid).First(&user)

	if user.ChatID == 0 {
		return 0, 0, errors.New("找不到用户!")
	}
	if user.Wallmoney <= 0 {
		return 0, 0, errors.New("金额不足")
	}
	addscore.Score = user.Wallmoney

	if addscore.Bet < 99.0 {
		floatscore = float64(user.Wallmoney) * addscore.Bet
		if user.Wallmoney < int64(floatscore) {
			return 0, 0, errors.New("金额不足!")
		}
		user.Wallmoney = user.Wallmoney - int64(floatscore)
		result := groupStorage.db.Model(&logic.User{}).Where("userid = ?", addscore.Userid).Update("wallmoney", gorm.Expr("wallmoney-?", int64(floatscore)))
		// result := groupStorage.db.Update(&user)
		if result.Error != nil {
			return 0, 0, errors.New("金额不足!")
		}
		addscore.Bet = floatscore
	} else {

		user.Wallmoney = user.Wallmoney - int64(addscore.Bet)
		result := groupStorage.db.Save(&user)
		if result.Error != nil {
			return 0, 0, errors.New("金额不足!")
		}

	}
	// user := groupStorage.db.get
	result := groupStorage.db.Create(addscore)
	return int64(addscore.Bet), user.Wallmoney, result.Error
}

//获取所有投注人
func (groupStorage *GamesMysql) BetInfos(playid string) ([]logic.Scorelogs, error) {

	var score []logic.Scorelogs

	result := groupStorage.db.Model(&logic.Scorelogs{}).Where("playid = ? order by createtime asc", playid).Find(&score)

	return score, result.Error
}

//获取所有投注人
func (groupStorage *GamesMysql) WriteChangeScore(scores []logic.Scorelogs) error {

	for _, v := range scores {
		var user logic.User
		user.Userid = v.Userid
		user.Wallmoney += v.Changescore
		result := groupStorage.db.Model(&logic.User{}).Where("userid = ?", v.Userid).Update("wallmoney", gorm.Expr("wallmoney-?", v.Changescore))
		if result.Error != nil {
			logger.Errorf("更新用户金额失败!")
			return errors.New("更新用户金额失败")
		}

		result = groupStorage.db.Model(&logic.Scorelogs{}).Where("userid = ? and playid =?", v.Userid, v.Playid).Updates(logic.Scorelogs{
			Userid:      v.Userid,
			Playid:      v.Playid,
			Chatid:      v.Chatid,
			Nameid:      v.Nameid,
			Bet:         v.Bet,
			Changescore: v.Changescore,
			Score:       v.Score,
			Status:      2,
		})

		if result.Error != nil {
			logger.Errorf("更新用户金额失败!")
			return errors.New("更新用户金额失败")
		}
	}

	return nil
}
