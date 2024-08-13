package auth

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/views/setting"
)

var sessions = dao.Sessions{}
var peers = dao.Peers{}

func LoginWindow(myApp fyne.App) {
	loginWindow := myApp.NewWindow("个人号机器人")
	pNum := config.Get("phoneNumber")
	config := dao.Config{}
	phoneNumber := binding.NewString()
	phoneNumber.Set(pNum)
	phoneNumberLabel := widget.NewLabel("手机号")
	phoneNumberEntry := newPhoneNumEntry(phoneNumber)

	phone := container.New(layout.NewFormLayout(), phoneNumberLabel, phoneNumberEntry)
	var button *widget.Button
	activity := widget.NewActivity()
	loginButton := func() {
		if e := phoneNumberEntry.Validate(); e != nil {
			dialog.ShowError(e, loginWindow)
			return
		}
		phoneNum, _ := phoneNumber.Get()
		if pNum != "" && phoneNum != pNum {
			dialog.ShowConfirm("警告", fmt.Sprintf("将清除手机号：%s 的数据", pNum), func(b bool) {
				if b {
					button.Disable()
					activity.Start()
					activity.Show()
					_ = sessions.DeleteAll()
					_ = peers.DeleteAll()
					config.Set("phoneNumber", phoneNum)
					ExpireWindow(loginWindow, myApp)
					return
				} else {
					activity.Stop()
					activity.Hide()
					button.Enable()
				}
			}, loginWindow)
		} else {
			button.Disable()
			activity.Start()
			activity.Show()
			config.Set("phoneNumber", phoneNum)
			ExpireWindow(loginWindow, myApp)
		}
	}
	button = widget.NewButton("启动", loginButton)
	showSetting := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		setting.ShowSettingModal(loginWindow)
	})
	button.Importance = widget.HighImportance
	loginContainer := container.NewVBox(phone, container.NewStack(button, activity))
	loginContainer.Resize(fyne.NewSize(250, 120))
	loginContainer.Move(fyne.NewPos(75, 100))

	showSetting.Resize(fyne.NewSize(40, 40))
	//固定窗口大小
	loginWindow.SetFixedSize(true)
	/*window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		winSize := window.Canvas().Size()
		// Calculate new position based on window size
		newX := float32(winSize.Width) * 0.25
		newY := float32(winSize.Height) * 0.25
		loginContainer.Move(fyne.NewPos(newX, newY))
	})*/
	loginWindow.SetContent(container.NewWithoutLayout(loginContainer, showSetting))
	loginWindow.Resize(fyne.NewSize(400, 400))
	loginWindow.CenterOnScreen()
	loginWindow.Show()

	return
}

type phoneNumEntry struct {
	widget.Entry
}

func newPhoneNumEntry(phoneNum binding.String) *phoneNumEntry {
	p := &phoneNumEntry{}
	p.ExtendBaseWidget(p)
	p.Bind(phoneNum)
	p.Validator = validation.NewRegexp(`\+\d+`, "手机号必须以+号加上区号开头")
	return p
}
