package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/views/dashbord"
	"github.com/midnightsong/telegram-assistant/views/setting"
)

var config = dao.Config{}

func ConfigWidget(window fyne.Window, myApp fyne.App) *fyne.Container {
	config := dao.Config{}
	phoneNumber := binding.NewString()
	phoneNumber.Set(config.Get("phoneNumber"))
	phoneNumberLabel := widget.NewLabel("手机号")
	phoneNumberEntry := newPhoneNumEntry(phoneNumber)

	phone := container.New(layout.NewFormLayout(), phoneNumberLabel, phoneNumberEntry)
	var button *widget.Button
	activity := widget.NewActivity()
	loginButton := func() {
		button.Disable()
		activity.Start()
		activity.Show()
		phoneNum, _ := phoneNumber.Get()
		config.Set("phoneNumber", phoneNum)

		go func() {
			go assistant.Run(myApp)
			go showVerify(myApp)
			go showPassword(myApp)
			dashbord.MsgNewWindow(window, myApp)
		}()
	}
	button = widget.NewButton("启动", loginButton)
	showSetting := widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
		showConfigDialog(window, myApp)
	})
	button.Importance = widget.HighImportance
	loginContainer := container.NewVBox(phone, container.NewStack(button, activity))
	loginContainer.Resize(fyne.NewSize(250, 120))
	loginContainer.Move(fyne.NewPos(75, 100))

	showSetting.Resize(fyne.NewSize(40, 40))
	//固定窗口大小
	window.SetFixedSize(true)
	/*window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		winSize := window.Canvas().Size()
		// Calculate new position based on window size
		newX := float32(winSize.Width) * 0.25
		newY := float32(winSize.Height) * 0.25
		loginContainer.Move(fyne.NewPos(newX, newY))
	})*/

	return container.NewWithoutLayout(loginContainer, showSetting)
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

func showConfigDialog(window fyne.Window, myApp fyne.App) {
	//var d *dialog.CustomDialog
	var up *widget.PopUp
	/*d = dialog.NewCustomWithoutButtons("配置", settingView, window)
	d.Resize(fyne.NewSize(300, 300))
	d.Show()*/
	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		up.Hide()
	})
	closeButton.Importance = widget.DangerImportance
	box := container.NewVBox(
		container.NewHBox(widget.NewLabel("设置"), layout.NewSpacer(), closeButton),
		setting.GetSettingView(window))
	up = widget.NewModalPopUp(box, window.Canvas())

	up.Show()
	up.Resize(fyne.NewSize(300, 300))
}
