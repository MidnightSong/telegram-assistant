package dao

import (
	"github.com/midnightsong/telegram-assistant/entities"
)

type FwdMsg struct {
}

func (FwdMsg) Insert(fm *entities.FwdMsg) error {
	return GetDb().Create(fm).Error
}

func (FwdMsg) GetFwd(fm entities.FwdMsg) (r entities.FwdMsg, e error) {
	e = GetDb().Where(fm).Take(&r).Error
	return
}

func (FwdMsg) All() (r []entities.FwdMsg, e error) {
	e = GetDb().Find(&r).Error
	return
}

func (FwdMsg) Delete(fm entities.FwdMsg) error {
	return GetDb().Delete(fm).Error
}
