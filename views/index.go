package views

import (
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/midnightsong/telegram-assistant/views/component"
)

var phoneNum string

func Run() {
	flag.StringVar(&phoneNum, "p", "", "手机号")
	/*myApp := app.NewWithID("com.tg-bot.preferences")
	loginWindown := myApp.NewWindow("登录")

	tabs := container.NewAppTabs(
		container.NewTabItem("登录", component.LoginWidget(loginWindown, myApp)),
		container.NewTabItem("配置", component.ConfigWidget(loginWindown, myApp)),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationLeading)
	loginWindown.Resize(fyne.NewSize(500, 600))
	loginWindown.SetContent(tabs)
	loginWindown.Show()
	myApp.Run()*/
	myApp := app.New()
	configWindow := myApp.NewWindow("个人号机器人")

	configWindow.SetContent(component.ConfigWidget(configWindow, myApp))
	configWindow.Resize(fyne.NewSize(400, 400))
	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("Telegram Bot",
			fyne.NewMenuItem("显示", func() {
				configWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}
	configWindow.Show()
	myApp.Run()
}
