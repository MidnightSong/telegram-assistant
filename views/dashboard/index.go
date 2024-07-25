package dashboard

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"os"
)

var config = dao.Config{}
var cli any
var tgClient *gotgproto.Client

func MsgNewWindow(jumpInWindow fyne.Window, myApp fyne.App) {
	cli = <-assistant.NewClient
	switch c := cli.(type) {
	case *gotgproto.Client:
		tgClient = c
	case string:
		logErr := dialog.NewError(fmt.Errorf("登录失败：%s", c), jumpInWindow)
		logErr.SetOnClosed(func() {
			os.Exit(0)
		})
		logErr.Show()
		assistant.NewClient <- nil
	}

	gotgproto.Logged = true
	dashboardWindow := myApp.NewWindow(fmt.Sprintf("欢迎：%s %s", tgClient.Self.FirstName, tgClient.Self.LastName))

	leftTabs := container.NewAppTabs(
		getMsgView(),                           //消息栏
		getSettingView(dashboardWindow, myApp), //设置栏
		GetLogOutView(dashboardWindow),         //注销栏
	)

	leftTabs.SetTabLocation(container.TabLocationLeading)

	dashboardWindow.Resize(fyne.NewSize(1024, 576))
	dashboardWindow.SetContent(leftTabs)
	dashboardWindow.CenterOnScreen()
	dashboardWindow.Show()
	jumpInWindow.Close()
	dashboardWindow.SetCloseIntercept(func() {
		dashboardWindow.Close()
		myApp.Quit()
		os.Exit(0)
	})
}
