package msg

import (
	"fmt"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/gotgproto/ext"
	"github.com/midnightsong/telegram-assistant/gotgproto/types"
	"github.com/midnightsong/telegram-assistant/utils"
	"regexp"
	"strings"
	"sync"
	"time"
)

var oMatch = regexp.MustCompile(`\w{8,}`)
var GroupRepeatMsgReplyTo bool       //关联回复重复过的机器人消息
var GroupHideSourceRepeatBotMsg bool //当重复消息时，是否隐藏来源
var fr = &dao.ForwardRelation{}
var CacheRelationsMap = sync.Map{}
var fwdMsgDao = dao.FwdMsg{}

// Init 使用init函数初始化变量会导致App启动异常
func Init() {
	// 查出所有的绑定关系集合，然后通过源id进行分类后缓存
	all := fr.All()
	for _, d := range all {
		key := d.PeerID
		var relations []*entities.ForwardRelation
		value, _ := CacheRelationsMap.Load(key)
		if value != nil {
			relations = value.([]*entities.ForwardRelation)
		} else {
			relations = []*entities.ForwardRelation{}
		}
		relations = append(relations, d)
		CacheRelationsMap.Store(key, relations)
	}
}

func getRelationsByPeer(peerID int64) (relations []*entities.ForwardRelation) {
	value, ok := CacheRelationsMap.Load(peerID)
	if !ok { //缓存中不存在，则从数据库查询
		find, _ := fr.Find(peerID)
		if len(find) == 0 {
			AddLog("收到消息，未绑定转发关系")
			empty := make([]*entities.ForwardRelation, 0)
			CacheRelationsMap.Store(peerID, empty)
			return nil
		}
		CacheRelationsMap.Store(peerID, find)
		value = find
	}
	relations = value.([]*entities.ForwardRelation)
	return
}

func saveForwardMsg(fromChatId, toChatId int64, fromMsgId, toMsgId int) {
	go func() {
		_ = fwdMsgDao.Insert(&entities.FwdMsg{
			OriginChatID: fromChatId,
			OriginMsgID:  fromMsgId,
			TargetChatID: toChatId,
			TargetMsgID:  toMsgId,
		})
	}()
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

func HandlerGroups(ctx *ext.Context, update *ext.Update) error {
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
		if oMatch.MatchString(update.EffectiveMessage.Text) {
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

	/*//如果没有打开重复机器人消息的设置，则不进行后续处理
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
	*/
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
				AddLog(fmt.Sprint("关联回复机器人消息:\n", update.EffectiveMessage.Text))
				//TODO 释放
				/*answer := ctx.Sender.Answer(*update.Entities, update.UpdateClass.(message.AnswerableMessageUpdate))
				f := entities.FwdMsg{
					ChatID: update.EffectiveChat().GetID(),
					MsgID:  replyTo.ID,
				}
				fwdMsgDao, e := dao.FwdMsg{}.GetFwd(f)
				go func() { AddLog(fmt.Sprint("回复原始消息", fwdMsgDao.FwdMsgID)) }()
				if e != nil {
					//utils.LogInfo(ctx, "没有找到转发消息的原始id：可能已经删除")
					AddLog(fmt.Sprint("没有找到转发消息的原始id：可能已经删除"))
					return true
				}
				//查找自己转发消息的原始消息id
				answer.Reply(fwdMsgDao.FwdMsgID).Text(ctx, update.EffectiveMessage.Text)
				return true*/
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
				var e error
				//当重复消息时，是否隐藏来源
				if GroupHideSourceRepeatBotMsg {
					//发送带图片的消息
					go func() {
						var username string
						for _, user := range update.Entities.Users {
							username += fmt.Sprintf("<@%s>+\n", user.Username)
						}
						AddLog(fmt.Sprintf("重复发送订单号类型的消息,隐藏来源:\n 消息：%s \n发送者：%s", msg, username))
					}()
					p := update.EffectiveMessage.Media.(*tg.MessageMediaPhoto)
					photo := &tg.InputPhoto{
						ID:            p.Photo.GetID(),
						AccessHash:    p.Photo.(*tg.Photo).AccessHash,
						FileReference: p.Photo.(*tg.Photo).FileReference,
					}
					b := ctx.Sender.Answer(*update.Entities, update.UpdateClass.(message.AnswerableMessageUpdate))
					_, e = b.Photo(ctx, photo, styling.Unknown(update.EffectiveMessage.Text))
				} else {
					_, e = ctx.ForwardMessages(update.EffectiveChat().GetID(), update.EffectiveChat().GetID(),
						&tg.MessagesForwardMessagesRequest{ID: []int{update.EffectiveMessage.ID},
							SendAs: ctx.Self.AsInputPeer(),
						})
					go func() {
						var username string
						for _, user := range update.Entities.Users {
							username += fmt.Sprintf("<@%s>+\n", user.Username)
						}
						AddLog(fmt.Sprintf("重复发送订单号类型的消息,显示来源:\n 消息：%s \n发送者：%s", msg, username))
					}()
				}
				if e != nil {
					AddLog("重复机器人消息失败：" + e.Error())
					return true
				}
				/*f := &entities.FwdMsg{
					ChatID:   update.EffectiveChat().GetID(),
					FwdMsgID: update.EffectiveMessage.ID,
					MsgID:    u.(*tg.Updates).Updates[0].(*tg.UpdateMessageID).ID,
					FwdTime:  time.Now().Unix(),
				}
				e = dao.FwdMsg{}.Insert(f)*/
				if e != nil {
					//utils.LogError(ctx, "保存重复发送的消息id失败"+e.Error())
					AddLog("保存重复发送的消息id失败" + e.Error())
				}
				return true
			}
		}
	}
	return false
}
