package entities

type ForwardRelation struct {
	ID           int64  `gorm:"column:id;primary_key;autoIncrement:true"`
	PeerID       int64  `gorm:"column:peer_id;not null;uniqueIndex:pid"`    //源会话id
	ToPeerID     int64  `gorm:"column:to_peer_id;not null;uniqueIndex:pid"` //目标会话id
	OnlyBot      bool   `gorm:"column:only_bot;default:true"`               //仅转发机器人消息
	ShowOrigin   bool   `gorm:"column:show_origin;default:true"`            //显示消息来源
	RelatedReply bool   `gorm:"column:related_reply;default:true"`          //关联转发回复
	Regex        string `gorm:"column:regex;default:''"`                    //触发转发消息的正则
	MustMedia    bool   `gorm:"column:must_media;default:true"`             //必须是带媒体内容的消息
}

func (ForwardRelation) TableName() string {
	return "forward_relation"
}