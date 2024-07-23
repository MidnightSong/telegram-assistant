package msg

import (
	"github.com/midnightsong/telegram-assistant/gotgproto/ext"
	"github.com/midnightsong/telegram-assistant/utils"
)

func HandlerChannel(ctx *ext.Context, update *ext.Update) error {
	utils.LogInfo(ctx.Context, "收到频道消息")
	/*
		//给消息点赞
		cool := &tg.MessagesSendReactionRequest{Reaction: []tg.ReactionClass{&tg.ReactionEmoji{Emoticon: "👍"}}}
		cool.MsgID = update.EffectiveMessage.ID
		ctx.SendReaction(update.EffectiveUser().GetID(), cool)

	*/

	/*
		//转发带图片消息消息
		ctx.ForwardMessages(update.EffectiveUser().GetID(), update.EffectiveUser().GetID(),
			&tg.MessagesForwardMessagesRequest{ID: []int{update.EffectiveMessage.ID},
				SendAs: ctx.Self.AsInputPeer(),
			})
	*/

	/*
		//发送带图片的消息
		p := update.EffectiveMessage.Media.(*tg.MessageMediaPhoto)
		photo := &tg.InputPhoto{
			ID:            p.Photo.GetID(),
			AccessHash:    p.Photo.(*tg.Photo).AccessHash,
			FileReference: p.Photo.(*tg.Photo).FileReference,
		}

		b := ctx.Sender.Answer(*update.Entities, update.UpdateClass.(message.AnswerableMessageUpdate))
		b.Photo(ctx, photo, styling.Unknown(update.EffectiveMessage.Text))

	*/
	return nil

}
