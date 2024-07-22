package entities

type Config struct {
	Key   string `json:"key" gorm:"column:key;primary_key;not null"`
	Value string `json:"value" gorm:"column:value;not null"`
}

func (Config) TableName() string {
	return "config"
}
