package dao

import (
	"github.com/midnightsong/telegram-assistant/entities"
)

type Peers struct {
}

func (Peers) Insert(fm *entities.Peers) error {
	return GetDb().Create(fm).Error
}

func (Peers) All() (r []entities.Peers, e error) {
	e = GetDb().Find(&r).Error
	return
}

// DeleteAll 删除全部
func (Peers) DeleteAll() error {
	return GetDb().Where("1=1").Delete(&entities.Peers{}).Error
}
