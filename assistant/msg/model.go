package msg

import "github.com/midnightsong/telegram-assistant/gotgproto/storage"

type DialogsInfo struct {
	Title string
	storage.EntityType
	PeerId int64
	Bot    bool
}
