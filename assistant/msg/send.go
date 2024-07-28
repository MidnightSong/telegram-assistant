package msg

import (
	"errors"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	mtp_errors "github.com/midnightsong/telegram-assistant/gotgproto/errors"
	"time"
)

var Client *gotgproto.Client

func SendMessage(chatId int64, msg string) (int, error) {
	time.Sleep(time.Millisecond * 30)
	result, err := Client.CreateContext().SendMessage(chatId, &tg.MessagesSendMessageRequest{Message: msg})
	if err != nil {
		return -1, err
	}
	return result.ID, nil
}

func DeleteMessage(chatId int64, msgId int) error {
	msgIds := []int{msgId}
	ctx := Client.CreateContext()
	err := ctx.DeleteMessages(chatId, msgIds)
	if errors.Is(err, mtp_errors.ErrNotChat) {
		_, err = ctx.Raw.MessagesDeleteMessages(ctx, &tg.MessagesDeleteMessagesRequest{
			Revoke: true,
			ID:     msgIds,
		})
	}
	return err
}
