package msg

import (
	"fmt"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/gotgproto/ext"
	"github.com/midnightsong/telegram-assistant/utils"
	"regexp"
	"strings"
	"time"
)

var oMatch = regexp.MustCompile(`\w{8,}`)
var GroupLog = make(chan string, 10)
var GroupRepeatMsg bool              //重复机器人消息
var GroupRepeatMsgReplyTo bool       //关联回复重复过的机器人消息
var GroupHideSourceRepeatBotMsg bool //当重复消息时，是否隐藏来源

func HandlerGroups(ctx *ext.Context, update *ext.Update) error {
	now := time.Now().Unix()
	if now-int64(update.EffectiveMessage.Date) > 60 {
		//utils.LogInfo(ctx, "读取到1分钟之前的消息，忽略")
		go func() { GroupLog <- "读取到1分钟之前的消息，忽略" }()
		return nil
	}
	//如果没有打开重复机器人消息的设置，则不进行后续处理
	if !GroupRepeatMsg {
		return nil
	}
	//关联回复重复过的机器人消息
	if groupRepeatMsgReplyToFunc(ctx, update) {
		return nil
	}
	//重复发送订单号类型的消息
	if GroupRepeatMsgFunc(ctx, update) {
		return nil
	}

	return nil
}

// groupRepeatMsgReplyToFunc 关联回复重复过的机器人消息
func groupRepeatMsgReplyToFunc(ctx *ext.Context, update *ext.Update) bool {
	if GroupRepeatMsgReplyTo {
		replyTo := update.EffectiveMessage.ReplyToMessage
		//回复自己的消息
		if replyTo != nil && ctx.Self.ID == replyTo.FromID.(*tg.PeerUser).UserID {
			//仅处理自己转发的消息
			if replyTo.FwdFrom.FromID != nil {
				go func() { GroupLog <- fmt.Sprint("关联回复机器人消息:\n", update.EffectiveMessage.Text) }()
				//TODO 释放
				/*answer := ctx.Sender.Answer(*update.Entities, update.UpdateClass.(message.AnswerableMessageUpdate))
				f := entities.FwdMsg{
					ChatID: update.EffectiveChat().GetID(),
					MsgID:  replyTo.ID,
				}
				fwd, e := dao.FwdMsg{}.GetFwd(f)
				go func() { GroupLog <- fmt.Sprint("回复原始消息", fwd.FwdMsgID) }()
				if e != nil {
					//utils.LogInfo(ctx, "没有找到转发消息的原始id：可能已经删除")
					go func() { GroupLog <- "没有找到转发消息的原始id：可能已经删除" }()
					return true
				}
				//查找自己转发消息的原始消息id
				answer.Reply(fwd.FwdMsgID).Text(ctx, update.EffectiveMessage.Text)*/
				return true
			}
			return true
		}
	}
	return false
}

// GroupRepeatMsgFunc 重复发送订单号类型的消息
func GroupRepeatMsgFunc(ctx *ext.Context, update *ext.Update) bool {
	msg := update.EffectiveMessage.Text
	if oMatch.MatchString(msg) && update.EffectiveMessage.Media != nil {
		if strings.Contains(msg, "http") {
			utils.LogInfo(ctx, "出现http，忽略")
			return true
		}
		for _, u := range update.Entities.Users {
			if u.Bot { //仅重复机器人的查单消息
				/*var u tg.UpdatesClass
				var e error*/
				//当重复消息时，是否隐藏来源
				if GroupHideSourceRepeatBotMsg {
					//发送带图片的消息
					//TODO 释放
					go func() {
						var username string
						for _, user := range update.Entities.Users {
							username += fmt.Sprintf("<@%s>+\n", user.Username)
						}
						GroupLog <- fmt.Sprintf("重复发送订单号类型的消息,隐藏来源:\n 消息：%s \n发送者：%s", msg, username)
					}()
					/*p := update.EffectiveMessage.Media.(*tg.MessageMediaPhoto)
					photo := &tg.InputPhoto{
						ID:            p.Photo.GetID(),
						AccessHash:    p.Photo.(*tg.Photo).AccessHash,
						FileReference: p.Photo.(*tg.Photo).FileReference,
					}
					b := ctx.Sender.Answer(*update.Entities, update.UpdateClass.(message.AnswerableMessageUpdate))
					u, e = b.Photo(ctx, photo, styling.Unknown(update.EffectiveMessage.Text))*/
				} else {
					//TODO 释放
					/*u, e = ctx.ForwardMessages(update.EffectiveChat().GetID(), update.EffectiveChat().GetID(),
					&tg.MessagesForwardMessagesRequest{ID: []int{update.EffectiveMessage.ID},
						SendAs: ctx.Self.AsInputPeer(),
					})*/
					go func() {
						var username string
						for _, user := range update.Entities.Users {
							username += fmt.Sprintf("<@%s>+\n", user.Username)
						}
						GroupLog <- fmt.Sprintf("重复发送订单号类型的消息,显示来源:\n 消息：%s \n发送者：%s", msg, username)
					}()
				}

				//utils.LogInfo(ctx, "获取到机器人发送订单号类型消息："+msg)
				/*if e != nil {
					//utils.LogError(ctx, "重复机器人消息失败："+e.Error())
					go func() { GroupLog <- "重复机器人消息失败：" + e.Error() }()
					return true
				}
				f := &entities.FwdMsg{
					ChatID:   update.EffectiveChat().GetID(),
					FwdMsgID: update.EffectiveMessage.ID,
					MsgID:    u.(*tg.Updates).Updates[0].(*tg.UpdateMessageID).ID,
					FwdTime:  time.Now().Unix(),
				}
				e = dao.FwdMsg{}.Insert(f)
				if e != nil {
					//utils.LogError(ctx, "保存重复发送的消息id失败"+e.Error())
					go func() { GroupLog <- "保存重复发送的消息id失败" + e.Error() }()
				}*/
				return true
			}
		}
	}
	return false
}
