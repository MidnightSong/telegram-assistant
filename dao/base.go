package dao

import (
	"context"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/utils"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/midnightsong/telegram-assistant/gotgproto/sessionMaker"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dial = sqlite.Open("database")
var SqlSession = sessionMaker.SqlSession(dial)
var (
	_once sync.Once
	_db   *gorm.DB
)

func GetDb() *gorm.DB {
	_once.Do(func() {
		var err error
		_db, err = gorm.Open(dial, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent), PrepareStmt: true, SkipDefaultTransaction: true})
		if err != nil {
			utils.LogError(context.TODO(), err.Error())
		}
		_db = _db.Debug()
	})
	return _db
}

// 建表和索引
func Migrate() error {
	return GetDb().AutoMigrate(&entities.FwdMsg{},
		&entities.Config{},
	)
}
