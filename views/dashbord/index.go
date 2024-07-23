package dashbord

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"os"
)

var config = dao.Config{}
var tgClient *gotgproto.Client

func MsgNewWindow(window fyne.Window, myApp fyne.App) {
	tgClient = <-assistant.NewClient
	gotgproto.Logged = true
	dashboardWindow := myApp.NewWindow(fmt.Sprintf("欢迎：%s %s", tgClient.Self.FirstName, tgClient.Self.LastName))

	leftTabs := container.NewAppTabs(
		getMsgView(),          //消息栏
		getSettingView(myApp), //设置栏
	)

	leftTabs.SetTabLocation(container.TabLocationLeading)

	dashboardWindow.Resize(fyne.NewSize(1024, 576))
	dashboardWindow.SetContent(leftTabs)
	dashboardWindow.Show()
	dashboardWindow.SetCloseIntercept(func() {
		dashboardWindow.Close()
		myApp.Quit()
		os.Exit(0)
	})
	window.Close()
}
