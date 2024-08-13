package auth

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"github.com/midnightsong/telegram-assistant/views/dashboard"
	"github.com/midnightsong/telegram-assistant/views/setting"
	"os"
	"time"
)

var config = dao.Config{}
var pass = make(chan bool)

func ExpireWindow(jumpInWindow fyne.Window, app fyne.App) {
	expireWindow := app.NewWindow("登录中")

	text := "正在验证信息中"
	label := widget.NewLabel(text)
	go func() {
		for i := 0; i <= 6; i++ {
			points := []string{"", ".", "..", "...", "....", ".....", "......"}
			label.SetText(text + points[i])
			label.Refresh()
			time.Sleep(500 * time.Millisecond)
			if i == 6 {
				i = 0
			}
		}
	}()
	label.Move(fyne.NewPos(20, 50))
	// 创建一个进度条以模拟活动图标
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Start()
	progressBar.Resize(fyne.NewSize(180, 20))
	progressBar.Move(fyne.NewPos(5, 90))

	content := container.NewWithoutLayout(
		label,
		progressBar,
	)
	expireWindow.SetContent(content)
	expireWindow.Resize(fyne.NewSize(200, 300))
	expireWindow.CenterOnScreen()
	expireWindow.SetFixedSize(true)
	expireWindow.Show()
	if jumpInWindow != nil {
		jumpInWindow.Close()
	}
	go func() {
		for {
			go checkAut(expireWindow, app)
			if <-pass {
				go assistant.Run()
				go verifyWindow(app)
				go passwordWindow(app)
				dashboard.MsgNewWindow(expireWindow, app)
				break
			}

		}
	}()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			if gotgproto.PhoneNumberErr != "" {
				dialog.ShowConfirm("错误", gotgproto.PhoneNumberErr, func(b bool) {
					os.Exit(0)
				}, expireWindow)
			}
		}
	}()
}

func showSettingWindow(app fyne.App) fyne.Window {
	settingWindow := app.NewWindow("配置")
	settingWindow.SetFixedSize(true)
	settingWindow.SetContent(setting.GetSettingView(settingWindow))
	settingWindow.Resize(fyne.NewSize(400, 400))
	settingWindow.Show()
	settingWindow.CenterOnScreen()
	return settingWindow
}

func waitingInput(app fyne.App) {
	//输入激活码的窗口
	inputted := make(chan bool)
	settingWindow := showSettingWindow(app)
	go func() {
		for {
			authCode := config.Get("authCode")
			appId := config.Get("appId")
			apiHash := config.Get("apiHash")
			if authCode != "" && appId != "" && apiHash != "" {
				inputted <- true
				return
			}
			time.Sleep(time.Second * 2)
		}
	}()
	<-inputted
	settingWindow.Close()
}
func checkAut(window fyne.Window, app fyne.App) {
	authCode := config.Get("authCode")
	appId := config.Get("appId")
	apiHash := config.Get("apiHash")
	//如果检测到配置为空，弹出输配置的窗口
	if authCode == "" || appId == "" || apiHash == "" {
		waitingInput(app)
	}
	result, err := assistant.Auth()
	if err != nil {
		errorDialog := dialog.NewError(err, window)
		errorDialog.Show()
		errorDialog.SetOnClosed(func() {
			waitingInput(app)
		})
		pass <- false
		return
	}
	if result.Code == 2000 {
		pass <- true
		dashboard.ExpireTime = result.Data.Exp
		return
	}
	errorDialog := dialog.NewError(errors.New("错误：\n"+result.Msg), window)
	errorDialog.Show()
	errorDialog.SetOnClosed(func() {
		_ = config.Set("authCode", "")
		pass <- false
		return
	})
}
