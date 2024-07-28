package dao

import (
	"github.com/midnightsong/telegram-assistant/entities"
	"gorm.io/gorm/clause"
)

type ForwardRelation struct{}

func (ForwardRelation) Add(fr *entities.ForwardRelation) error {
	return GetDb().Clauses(clause.Insert{Modifier: "OR REPLACE"}).Create(fr).Error
}
func (ForwardRelation) DeleteById(id int64) {
	GetDb().Where("id = ?", id).Delete(&entities.ForwardRelation{})
}

func (ForwardRelation) Delete(id int64) {
	GetDb().Where("peer_id = ?", id).Delete(&entities.ForwardRelation{})
}

func (ForwardRelation) Find(id int64) ([]*entities.ForwardRelation, error) {
	var frs []*entities.ForwardRelation
	err := GetDb().Where("peer_id = ?", id).Find(&frs).Error
	return frs, err
}
