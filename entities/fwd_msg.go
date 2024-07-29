package entities

type FwdMsg struct {
	ID           uint64 `gorm:"column:id;primary_key;autoIncrement:true"`
	OriginChatID int64  `gorm:"column:origin_chat_id;not null"`
	TargetChatID int64  `gorm:"column:target_chat_id;not null"`
	FwdMsgID     int    `gorm:"column:origin_msg_id;not null"`
	TargetMsgID  int    `gorm:"column:target_msg_id;not null"`
	FwdTime      int64  `gorm:"column:fwd_time;not null"`
}

func (FwdMsg) TableName() string {
	return "fwd_msg"
}
