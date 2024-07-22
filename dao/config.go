package dao

import (
	"github.com/midnightsong/telegram-assistant/entities"
	"gorm.io/gorm/clause"
)

type Config struct {
}

func (Config) Set(key string, value string) error {
	config := entities.Config{
		Key:   key,
		Value: value,
	}
	return GetDb().Clauses(clause.Insert{Modifier: "OR REPLACE"}).Create(config).Error
}

func (Config) Get(key string) (value string) {
	c := entities.Config{}
	GetDb().Select("value").Where("key = ?", key).Take(&c)
	return c.Value
}

func (Config) All() (r []entities.Config, e error) {
	e = GetDb().Find(&r).Error
	return
}

func (Config) Delete(fm entities.Config) error {
	return GetDb().Delete(fm).Error
}
