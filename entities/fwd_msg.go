package entities

type FwdMsg struct {
	ID       uint64 `json:"id" gorm:"column:id;primary_key;autoIncrement:true"`
	ChatID   int64  `json:"chat_id" gorm:"column:chat_id;not null"`
	FwdMsgID int    `json:"fwd_msg_id" gorm:"column:fwd_msg_id;not null"`
	MsgID    int    `json:"msg_id" gorm:"column:msg_id;not null"`
	FwdTime  int64  `json:"fwd_time" gorm:"column:fwd_time;not null"`
}

func (FwdMsg) TableName() string {
	return "fwd_msg"
}
