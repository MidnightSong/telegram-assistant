package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/views/icon"
	"time"
)

var sessions = dao.Sessions{}
var peers = dao.Peers{}
var ExpireTime int

func GetLogOutView(window fyne.Window) *container.TabItem {
	ac := widget.NewLabel("激活有效期至:")
	ac1 := widget.NewLabel(time.Unix(int64(ExpireTime), 0).Format(time.DateTime))

	var logOutButton *widget.Button
	confirm := func(b bool) {
		if b {
			sessions.DeleteAll()
			peers.DeleteAll()
		}
	}
	logOutButton = widget.NewButton("注销当前账号", func() {
		dialog.NewConfirm("注销", "注销账号将会清除当前账号内所有记录，并退出程序，确认执行吗？", confirm, window).Show()
	})
	logOutButton.Importance = widget.DangerImportance
	testButton := widget.NewButton("test", func() {
		//ShowSendMsgModal(window)
		msg.Refresh2()
	})
	testButton.Hide()
	logOutBox := container.NewVBox(ac, ac1, logOutButton, testButton)

	return container.NewTabItemWithIcon("", icon.GetIcon(icon.ShutDown), logOutBox)
}
