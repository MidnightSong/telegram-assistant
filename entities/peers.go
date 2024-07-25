package entities

type Peers struct {
	ID         int64  `json:"id" gorm:"column:id;primary_key;autoIncrement:true"`
	AccessHash int64  `json:"access_hash" gorm:"column:access_hash"`
	Type       int64  `json:"type" gorm:"column:type"`
	Username   string `json:"username" gorm:"column:username;"`
}

func (Peers) TableName() string {
	return "peers"
}
