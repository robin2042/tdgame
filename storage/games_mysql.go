package storage

import (
	"errors"
	"fmt"
	"tdgames/logger"
	"tdgames/logic"

	"gorm.io/gorm"
)

const (
	GAME_END = 2
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

//下注金额，错误
func (groupStorage *GamesMysql) AddScore(addscore *logic.AddScore) (int64, error) {

	var user logic.User
	uid := fmt.Sprintf("%d%d", addscore.Userid, addscore.Chatid)

	tx := groupStorage.db.Begin()
	defer tx.Commit()

	tx.Where("uid=?", uid).First(&user)

	if user.ChatID == 0 {
		return 0, errors.New("找不到用户")
	}
	if user.Wallmoney <= 0 {
		return 0, errors.New("金额不足")
	}
	addscore.Uid = uid //用户ID
	addscore.Score = user.Wallmoney

	result := tx.Model(&logic.User{}).Where("userid = ?", addscore.Userid).Update("wallmoney", gorm.Expr("wallmoney-?", addscore.Bet))
	// result := groupStorage.db.Update(&user)
	if result.Error != nil {
		tx.Rollback()
		return 0, errors.New("金额不足")
	}

	groupStorage.db.Create(addscore)

	return int64(addscore.Bet), result.Error
}

//获取所有投注人
func (groupStorage *GamesMysql) BetInfos(playid string) ([]logic.Scorelogs, error) {

	var score []logic.Scorelogs

	result := groupStorage.db.Model(&logic.Scorelogs{}).Where("playid = ? order by createtime asc", playid).Find(&score)

	return score, result.Error
}

//获取所有投注人
func (groupStorage *GamesMysql) WriteUserRecords(playid string, scores []logic.Scorelogs) error {

	for _, v := range scores {

		result := groupStorage.db.Model(&logic.Scorelogs{}).Where("userid = ? and playid =?", v.Userid, v.Playid).Updates(logic.Scorelogs{
			Userid:      v.Userid,
			Playid:      v.Playid,
			Chatid:      v.Chatid,
			Nameid:      v.Nameid,
			Bet:         v.Bet,
			Changescore: v.Changescore,
			Score:       v.Score,
			Area:        v.Area, //下注区域
			Status:      2,
		})

		if result.Error != nil {
			logger.Errorf("更新用户金额失败")
			return errors.New("更新用户金额失败")
		}
	}

	return nil
}

//获取所有投注人
func (groupStorage *GamesMysql) WriteChangeScore(playid string, chatid int64, users map[int64]int64) error {

	//更新本局结束
	groupStorage.db.Model(&logic.Gamerounds{}).Where("playid = ?", playid).Update("status", 2)

	for k, v := range users {
		var user logic.User
		user.Userid = k
		user.Wallmoney += v
		result := groupStorage.db.Model(&logic.User{}).Where("userid = ? ", k).Update("wallmoney", gorm.Expr("wallmoney+?", v))
		if result.Error != nil {
			logger.Errorf("更新用户金额失败")
			return errors.New("更新用户金额失败")
		}

		if result.Error != nil {
			logger.Errorf("更新用户金额失败")
			return errors.New("更新用户金额失败")
		}
	}

	return nil
}

//判断是否能开局
func (groupStorage *GamesMysql) NewGames(nameid int, chatid int64) error {

	var game logic.Gamerounds

	result := groupStorage.db.Model(&logic.Gamerounds{}).Where("nameid = ? and chatid =? ", nameid, chatid).Order("createtime desc").Limit(1).Find(&game)

	// timer, _ := time.ParseInLocation("2006-01-02 15:04:05", game.Createtime, time.Local)

	if game.Status == GAME_END {
		return nil
	}
	// if time.Since(timer).Seconds() <= 90 {
	// 	return errors.New("上局90s后才能开始游戏")
	// }

	if result.Error != nil {
		return errors.New("上局90s后才能开始游戏")
	}
	return nil
}
func (groupStorage *GamesMysql) GetRecords(nameid int, chatid int64) []logic.Records {
	var game []logic.Records

	groupStorage.db.Model(&logic.Gamerounds{}).Where("nameid = ? and chatid =?  and status =2", nameid, chatid).Order("createtime desc").Limit(10).Find(&game)

	return game
}
