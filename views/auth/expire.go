package auth

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/utils"
	"github.com/midnightsong/telegram-assistant/views/dashboard"
	"github.com/midnightsong/telegram-assistant/views/setting"
	"os"
	"time"
)

func ExpireWindow(jumpInWindow fyne.Window, app fyne.App) {
	pass := make(chan bool)
	config := dao.Config{}

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
	expireWindow.Resize(fyne.NewSize(200, 200))
	expireWindow.CenterOnScreen()
	expireWindow.SetFixedSize(true)
	expireWindow.Show()
	if jumpInWindow != nil {
		jumpInWindow.Close()
	}
	go func() {
		authCode := config.Get("authCode")
		appId := config.Get("appId")
		apiHash := config.Get("apiHash")
		//如果检测到配置为空，弹出输配置的窗口
		if authCode == "" || appId == "" || apiHash == "" {
			//输入激活码的窗口
			inputted := make(chan bool)
			settingWindow := showSettingWindow(app)
			go func() {
				for {
					authCode = config.Get("authCode")
					appId = config.Get("appId")
					apiHash = config.Get("apiHash")
					if authCode != "" && appId != "" && apiHash != "" {
						inputted <- true
						return
					}
					time.Sleep(time.Second)
				}
			}()
			<-inputted
			settingWindow.Close()
		}
		name, err := os.Hostname()
		if err != nil {
			errorDialog := dialog.NewError(errors.New("获取设备信息失败，请联系客服处理"+err.Error()), expireWindow)
			errorDialog.Show()
			errorDialog.SetOnClosed(func() {
				os.Exit(0)
			})
			return
		}
		params := make(map[string]interface{})
		params["device_id"] = name
		params["uuid"] = authCode
		params["timestamp"] = time.Now().Unix()
		result := &AuthResponse{}
		err = utils.HttpClient.Post("https://auth.seven-d76.workers.dev/acv", params, result)
		if err != nil {
			errorDialog := dialog.NewError(errors.New("内部错误："+err.Error()), expireWindow)
			errorDialog.Show()
			errorDialog.SetOnClosed(func() {
				os.Exit(0)
			})
			return
		}
		if result.Code == 2000 {
			pass <- true
			return
		}
		errorDialog := dialog.NewError(errors.New("错误："+result.Msg), expireWindow)
		errorDialog.Show()
		errorDialog.SetOnClosed(func() {
			_ = config.Set("authCode", "")
			os.Exit(0)
		})
		return
	}()
	go func() {
		if <-pass {
			go assistant.Run()
			go verifyWindow(app)
			go passwordWindow(app)
			dashboard.MsgNewWindow(expireWindow, app)
		}
	}()
}

func showSettingWindow(app fyne.App) fyne.Window {
	settingWindow := app.NewWindow("配置")
	settingWindow.SetFixedSize(true)
	settingWindow.SetContent(setting.GetSettingView(settingWindow))
	settingWindow.Resize(fyne.NewSize(300, 400))
	settingWindow.Show()
	return settingWindow
}

type AuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UUID      string `json:"uuid"`
		Exp       int    `json:"exp"`
		Duration  int    `json:"duration"`
		Timestamp int    `json:"timestamp"`
	} `json:"data"`
}
