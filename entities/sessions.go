package entities

type Sessions struct {
	Version int64  `json:"version" gorm:"column:version;primary_key;autoIncrement:true"`
	Data    string `json:"data" gorm:"column:data;not null"`
}

func (Sessions) TableName() string {
	return "sessions"
}
