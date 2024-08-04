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

var SqlSession *sessionMaker.SqlSessionConstructor
var (
	_once sync.Once
	_db   *gorm.DB
)
var DbPath string

func GetDb() *gorm.DB {
	_once.Do(func() {
		var err error
		var dial = sqlite.Open(DbPath)
		SqlSession = sessionMaker.SqlSession(dial)
		_db, err = gorm.Open(dial, &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 logger.Default.LogMode(logger.Silent),
		})
		d, _ := _db.DB()
		d.SetMaxOpenConns(100)
		_ = _db.AutoMigrate(&entities.FwdMsg{}, &entities.Config{}, &entities.ForwardRelation{})
		if err != nil {
			utils.LogError(context.TODO(), err.Error())
		}
		_db = _db.Debug()
	})
	return _db
}
