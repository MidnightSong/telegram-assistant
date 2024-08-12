package views

import (
	"runtime"
	"strings"

	"fyne.io/fyne/v2/app"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/job"
	"github.com/midnightsong/telegram-assistant/views/auth"
)

func Run() {
	myApp := app.NewWithID("com.panda.assistant")
	//myApp.Settings().SetTheme(custom.Theme{})
	/*if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("Telegram Bot",
			fyne.NewMenuItem("显示", func() {
				loginWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}*/
	s := myApp.Storage()
	path := s.RootURI().String()
	/*	create, err := s.Create("../cache.db")
		if err == nil {
			create.Close()
		}*/
	if runtime.GOOS == "windows" {
		path = strings.ReplaceAll(path, "file://", "")
	}
	dao.DbPath = path + "/cache.db"

	_, err := dao.Sessions{}.GetSession(entities.Sessions{Version: 1})
	if err != nil {
		auth.LoginWindow(myApp)
	} else {
		auth.ExpireWindow(nil, myApp)
	}
	go func() {
		msg.Init()
	}()
	job.Init()
	myApp.Run()
}
