package msg

import (
	"fmt"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/gotgproto/ext"
)

var PrivateRepeatMsg bool

func HandlerPrivate(ctx *ext.Context, update *ext.Update) error {
	user := update.EffectiveUser()
	go func() {
		GroupLog <- fmt.Sprintf(
			"======收到私聊消息======\n%s: %s", user.FirstName+user.LastName, update.EffectiveMessage.Text)
	}()
	if PrivateRepeatMsg {
		//给消息点赞
		cool := &tg.MessagesSendReactionRequest{Reaction: []tg.ReactionClass{&tg.ReactionEmoji{Emoticon: "👍"}}}
		cool.MsgID = update.EffectiveMessage.ID
		ctx.SendReaction(user.GetID(), cool)
	}

	//查看所有已打开的聊天窗口
	/*dialogs, err := ctx.Raw.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{},
	})
	if err != nil {
		return err
	}
	utils.LogInfo(ctx.Context, "dialogs", zap.Any("dialogs", dialogs))*/
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
