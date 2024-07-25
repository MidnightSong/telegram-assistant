package dao

import (
	"github.com/midnightsong/telegram-assistant/entities"
)

type Sessions struct {
}

func (Sessions) Insert(s *entities.Sessions) error {
	return GetDb().Create(s).Error
}

func (Sessions) GetSession(s entities.Sessions) (r entities.Sessions, e error) {
	e = GetDb().Where(s).Take(&r).Error
	return
}

func (Sessions) All() (r []entities.Sessions, e error) {
	e = GetDb().Find(&r).Error
	return
}

func (Sessions) DeleteAll() error {
	return GetDb().Where("1=1").Delete(&entities.Sessions{}).Error
}
