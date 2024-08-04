package msg

import (
	"fmt"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/gotgproto/ext"
	"github.com/midnightsong/telegram-assistant/gotgproto/types"
	"github.com/midnightsong/telegram-assistant/utils"
	"time"
)

func ProcessMsg(ctx *ext.Context, update *ext.Update) error {
	if update.EffectiveMessage.IsService {
		return nil
	}
	if update.EffectiveUser().Self {
		return nil
	}
	now := time.Now().Unix()
	if now-int64(update.EffectiveMessage.Date) > 60 {
		//utils.LogInfo(ctx, "读取到1分钟之前的消息，忽略")
		AddLog("读取到1分钟之前的消息，忽略")
		return nil
	}
	if update.EffectiveMessage.ReplyToMessage != nil {
		processReply(ctx, update)
		return nil
	}

	//检查绑定关系
	messageFromPeerId := update.EffectiveChat().GetID()
	relations := getRelationsByPeer(messageFromPeerId)
	//可能绑定了多个需要转发的目标
	for _, relation := range relations {
		//仅转发机器人消息开关打开
		if relation.OnlyBot && !update.EffectiveUser().Bot {
			continue
		}
		//只转发带图片的消息开关打开
		if relation.MustMedia && update.EffectiveMessage.Media == nil {
			continue
		}
		//满足文字条件
		reg, err := getReg(relation.Regex)
		if err != nil {
			log := fmt.Sprintf("转发消息失败：\n%s:%s\n原因: %s", update.EffectiveUser().FirstName+update.EffectiveUser().LastName, update.EffectiveMessage.Text, err.Error())
			AddLog(log)
			continue
		}
		if reg.MatchString(update.EffectiveMessage.Text) {

			//显示消息来源
			if relation.ShowOrigin {
				updatesClass, err := ForwardMessage(messageFromPeerId, relation.ToPeerID, []int{update.EffectiveMessage.ID})
				if err != nil {
					AddLog("转发消息错误（显示消息来源）：" + err.Error())
					utils.LogWarn(ctx.Context, "转发消息错误（显示消息来源）："+err.Error())
					continue
				}
				targetMsgID := updatesClass.(*tg.Updates).Updates[0].(*tg.UpdateMessageID).ID
				saveForwardMsg(messageFromPeerId, relation.ToPeerID, update.EffectiveMessage.ID, targetMsgID)
				continue
			}
			//隐藏消息来源（生成新的消息,不带图片）
			if update.EffectiveMessage.Media == nil {
				targetMsgID, err := SendMessage(relation.ToPeerID, update.EffectiveMessage.Text)
				if err != nil {
					AddLog("转发消息错误（隐藏消息来源,无图）：" + err.Error())
					utils.LogWarn(ctx.Context, "转发消息错误（隐藏消息来源,无图）"+err.Error())
					continue
				}
				saveForwardMsg(messageFromPeerId, relation.ToPeerID, update.EffectiveMessage.ID, targetMsgID)
				continue
			}
			//隐藏消息来源（生成新的消息,带图片）
			targetMsgID, err := SendMessageWithPhoto(relation.ToPeerID, update.EffectiveMessage.Media.(*tg.MessageMediaPhoto).Photo)
			if err != nil {
				AddLog("转发消息错误（隐藏消息来源，有图）：" + err.Error())
				utils.LogWarn(ctx.Context, "转发消息错误（隐藏消息来源，有图）:"+err.Error())
				continue
			}
			saveForwardMsg(messageFromPeerId, relation.ToPeerID, update.EffectiveMessage.ID, targetMsgID.ID)
		}
	}
	return nil
}

// 回复类型的消息
func processReply(ctx *ext.Context, update *ext.Update) {
	replyTo := update.EffectiveMessage.ReplyToMessage
	//仅处理回复自己的消息
	if ctx.Self.ID == replyTo.FromID.(*tg.PeerUser).UserID {
		beforeForward := &entities.FwdMsg{
			TargetChatID: update.EffectiveChat().GetID(),
			TargetMsgID:  replyTo.ID,
		}
		e := dao.FwdMsg{}.GetFwd(beforeForward)
		if e != nil {
			log := fmt.Sprintf("收到无原始消息回复:\n %s: %s\n", update.EffectiveUser().FirstName+update.EffectiveUser().LastName, update.EffectiveMessage.Text)
			AddLog(log)
			return
		}

		//检查原始消息会话配置
		originChatRelations := getRelationsByPeer(beforeForward.OriginChatID)
		for _, originChatRelation := range originChatRelations {
			if originChatRelation.ToPeerID != update.EffectiveChat().GetID() {
				continue
			}
			//找到原本Chat的绑定配置，检查是否打开了关联回复
			if originChatRelation.RelatedReply {
				var msg *types.Message
				var err error
				if update.EffectiveMessage.Media != nil {
					p := update.EffectiveMessage.Media.(*tg.MessageMediaPhoto)
					msg, err = SendReplyMessageWhitPhoto(beforeForward.OriginChatID, beforeForward.OriginMsgID, update.EffectiveMessage.Text, p.Photo)
				} else {
					msg, err = SendReplyMessage(beforeForward.OriginChatID, beforeForward.OriginMsgID, update.EffectiveMessage.Text)
				}

				if err != nil {
					AddLog("关联回复消息错误：" + err.Error())
					utils.LogWarn(ctx.Context, "关联回复消息错误:"+err.Error())
					continue
				}
				saveForwardMsg(update.EffectiveChat().GetID(), beforeForward.OriginChatID, update.EffectiveMessage.ID, msg.ID)
			}
		}
	}
}
