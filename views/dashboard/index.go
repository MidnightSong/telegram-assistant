package dashboard

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"github.com/midnightsong/telegram-assistant/views/icon"
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
		logFail := myApp.NewWindow("登录失败")
		errContent := widget.NewRichTextWithText(c)
		errContent.Wrapping = fyne.TextWrapWord
		errContent.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignCenter

		shutDownButton := widget.NewButtonWithIcon("注销", icon.ShutDown, func() {
			dialog.NewConfirm("注销？", "注销将删除当前账号所有数据，不可恢复。", func(b bool) {
				if b {
					_ = sessions.DeleteAll()
					_ = peers.DeleteAll()
					os.Exit(0)
				}
			}, logFail).Show()
		})

		shutDownButton.Importance = widget.WarningImportance
		exitButton := widget.NewButtonWithIcon("退出", theme.CancelIcon(), func() {
			os.Exit(0)
		})
		exitButton.Importance = widget.HighImportance
		box := container.NewVBox(errContent, container.NewHBox(exitButton, layout.NewSpacer(), shutDownButton))
		logFail.Resize(fyne.NewSize(300, 400))
		logFail.SetContent(box)
		logFail.Show()
		logFail.CenterOnScreen()
		logFail.SetFixedSize(true)
		jumpInWindow.Close()
		return
	}

	gotgproto.Logged = true
	dashboardWindow := myApp.NewWindow(fmt.Sprintf("欢迎：%s %s", tgClient.Self.FirstName, tgClient.Self.LastName))

	leftTabs := container.NewAppTabs(
		getMsgView(dashboardWindow),            //消息栏
		getSettingView(dashboardWindow, myApp), //设置栏
		GetLogOutView(dashboardWindow),         //注销栏
		getForwardView(dashboardWindow),        //搬运
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
