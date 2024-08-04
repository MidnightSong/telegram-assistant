package msg

import (
	"errors"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	mtp_errors "github.com/midnightsong/telegram-assistant/gotgproto/errors"
	"github.com/midnightsong/telegram-assistant/gotgproto/types"
	"sync"
	"time"
)

var Client *gotgproto.Client
var lock sync.Mutex

func SendMessage(chatId int64, msg string) (int, error) {
	time.Sleep(time.Millisecond * 30)
	lock.Lock()
	result, err := Client.CreateContext().SendMessage(chatId, &tg.MessagesSendMessageRequest{Message: msg})
	lock.Unlock()
	if err != nil {
		return -1, err
	}
	return result.ID, nil
}

func DeleteMessage(chatId int64, msgId int) error {
	msgIds := []int{msgId}
	ctx := Client.CreateContext()
	lock.Lock()
	err := ctx.DeleteMessages(chatId, msgIds)
	lock.Unlock()
	if errors.Is(err, mtp_errors.ErrNotChat) {
		_, err = ctx.Raw.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
			Revoke: true,
			ID:     msgIds,
		})
	}
	return err
}

func SendReplyMessage(chatId int64, msgId int) (*types.Message, error) {
	request := &tg.MessagesSendMessageRequest{}
	//request.SetSendAs(Client.Self.AsInputPeer())
	request.SetReplyTo(&tg.InputReplyToMessage{
		ReplyToMsgID: msgId,
	})
	lock.Lock()
	message, err := Client.CreateContext().SendMessage(chatId, request)
	lock.Unlock()
	return message, err
}

func SendReplyMessageWhitPhoto(chatId int64, msgId int, photoId int64) (*types.Message, error) {
	photo := &tg.InputPhoto{
		ID: photoId,
	}
	request := &tg.MessagesSendMediaRequest{}
	request.Media = &tg.InputMediaPhoto{ID: photo}
	request.SetReplyTo(&tg.InputReplyToMessage{
		ReplyToMsgID: msgId,
	})
	lock.Lock()
	message, err := Client.CreateContext().SendMedia(chatId, request)
	lock.Unlock()
	return message, err
}

func ForwardMessage(fromChatId, toChatId int64) (tg.UpdatesClass, error) {
	ctx := Client.CreateContext()
	request := &tg.MessagesForwardMessagesRequest{
		SendAs: ctx.Self.AsInputPeer(),
	}
	lock.Lock()
	message, err := ctx.ForwardMessages(fromChatId, toChatId, request)
	lock.Unlock()
	return message, err
}
