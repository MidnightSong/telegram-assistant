package views

import (
	"fyne.io/fyne/v2/app"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/views/auth"
)

func Run() {
	myApp := app.New()
	/*if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("Telegram Bot",
			fyne.NewMenuItem("显示", func() {
				loginWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}*/
	d, err := dao.Sessions{}.GetSession(entities.Sessions{Version: 1})
	if err != nil {
		auth.LoginWindow(myApp)
	} else {
		auth.ExpireWindow(nil, myApp, d)
	}

	myApp.Run()
}
