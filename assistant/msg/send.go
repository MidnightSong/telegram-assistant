package msg

import (
	"errors"
	"github.com/gotd/td/tg"
	mtp_errors "github.com/midnightsong/telegram-assistant/gotgproto/errors"
	"github.com/midnightsong/telegram-assistant/gotgproto/types"
	"sync"
	"time"
)

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

func SendMessageWithPhoto(chatId int64, photo tg.PhotoClass) (*types.Message, error) {
	photoID := &tg.InputPhoto{
		ID:            photo.GetID(),
		AccessHash:    photo.(*tg.Photo).AccessHash,
		FileReference: photo.(*tg.Photo).FileReference,
	}
	request := &tg.MessagesSendMediaRequest{}
	request.Media = &tg.InputMediaPhoto{ID: photoID}
	lock.Lock()
	message, err := Client.CreateContext().SendMedia(chatId, request)
	lock.Unlock()
	return message, err
}

func SendReplyMessage(chatId int64, msgId int, msg string) (*types.Message, error) {
	ctx := Client.CreateContext()
	request := &tg.MessagesSendMessageRequest{Message: msg}
	request.SetReplyTo(&tg.InputReplyToMessage{
		ReplyToMsgID: msgId,
	})
	lock.Lock()
	message, err := ctx.SendMessage(chatId, request)
	lock.Unlock()
	return message, err
}

func SendReplyMessageWhitPhoto(chatId int64, msgId int, msg string, photo tg.PhotoClass) (*types.Message, error) {
	photoID := &tg.InputPhoto{
		ID:            photo.GetID(),
		AccessHash:    photo.(*tg.Photo).AccessHash,
		FileReference: photo.(*tg.Photo).FileReference,
	}
	request := &tg.MessagesSendMediaRequest{Message: msg}
	request.Media = &tg.InputMediaPhoto{ID: photoID}
	request.SetReplyTo(&tg.InputReplyToMessage{
		ReplyToMsgID: msgId,
	})
	lock.Lock()
	message, err := Client.CreateContext().SendMedia(chatId, request)
	lock.Unlock()
	return message, err
}

func ForwardMessage(fromChatId, toChatId int64, msgIDs []int) (tg.UpdatesClass, error) {
	ctx := Client.CreateContext()
	request := &tg.MessagesForwardMessagesRequest{
		ID:     msgIDs,
		SendAs: ctx.Self.AsInputPeer(),
	}
	lock.Lock()
	message, err := ctx.ForwardMessages(fromChatId, toChatId, request)
	lock.Unlock()
	return message, err
}
