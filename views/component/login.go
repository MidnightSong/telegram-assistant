package component

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/bot"
	"github.com/midnightsong/telegram-assistant/dao"
	"strconv"
)

var ipAddressRegex = "((25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]?\\d)\\.){3}(25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]?\\d)"
var numberRegex = "(\\d{1,5})"

func ConfigWidget(window fyne.Window, myApp fyne.App) *fyne.Container {
	config := dao.Config{}
	phoneNumber := binding.NewString()
	phoneNumber.Set(config.Get("phoneNumber"))
	phoneNumberLabel := widget.NewLabel("手机号")
	phoneNumberEntry := newPhoneNumEntry(phoneNumber)

	phone := container.New(layout.NewFormLayout(), phoneNumberLabel, phoneNumberEntry)
	var button *widget.Button
	activity := widget.NewActivity()
	confirmSave := func() {
		button.Disable()
		activity.Start()
		activity.Show()
		phoneNum, _ := phoneNumber.Get()
		config.Set("phoneNumber", phoneNum)

		go func() {
			go bot.Run(myApp)
			go showVerify(myApp)
			go showPassword(myApp)
			MsgNewWindow(window, myApp)
		}()
	}
	button = widget.NewButton("启动", confirmSave)
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

type appIDEntry struct {
	widget.Entry
}

func newAppIDEntry(appID binding.String) *appIDEntry {
	p := &appIDEntry{}
	p.ExtendBaseWidget(p)
	p.Bind(appID)
	p.Password = true
	p.Validator = validation.NewRegexp(`^\d+$`, "必须是纯数字")
	return p
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
	var d *dialog.CustomDialog

	secretKey := binding.NewString()
	secretKey.Set(config.Get("secretKey"))
	secretKeyLabel := widget.NewLabel("激活码")
	secretKeyEntry := widget.NewEntryWithData(secretKey)
	secretKeyEntry.Password = true

	appId := binding.NewString()
	appId.Set(config.Get("appId"))
	appIdLabel := widget.NewLabel("appId")
	appIdEntry := newAppIDEntry(appId)

	appHash := binding.NewString()
	appHash.Set(config.Get("apiHash"))
	appHashLabel := widget.NewLabel("apiHash")
	appHashEntry := widget.NewEntryWithData(appHash)
	appHashEntry.Password = true

	socksIP := binding.NewString()
	socksIP.Set(config.Get("socksIP"))
	socksIPLabel := widget.NewLabel("IP:")
	socksIPEntry := widget.NewEntry()
	socksIPEntry.PlaceHolder = "例: 127.0.0.1"
	socksIPEntry.Bind(socksIP)
	socksIPEntry.Validator = validation.NewRegexp(ipAddressRegex, "请输入正确的ip地址！")

	socksPort := binding.NewString()
	socksPort.Set(config.Get("socksPort"))
	socksPortLabel := widget.NewLabel("端口:")
	socksPortEntry := widget.NewEntry()
	socksPortEntry.PlaceHolder = "例: 8080"
	socksPortEntry.Bind(socksPort)
	socksPortEntry.Validator = validation.NewRegexp(numberRegex, "端口号有误！")

	checkSocks := func(b bool) {
		if b {
			socksIPEntry.Enable()
			socksPortEntry.Enable()
		} else {
			socksIPEntry.Disable()
			socksPortEntry.Disable()
			socksIPEntry.Text = ""
			socksPortEntry.Text = ""
		}
		config.Set("socksOpen", fmt.Sprint(b))
	}
	socksOpen := widget.NewCheck("socks5 代理", checkSocks)

	parseBool, _ := strconv.ParseBool(config.Get("socksOpen"))
	socksOpen.SetChecked(parseBool)
	checkSocks(parseBool)
	socksOpen.Refresh()

	configContainer := container.New(layout.NewFormLayout(),
		secretKeyLabel, secretKeyEntry,
		appIdLabel, appIdEntry,
		appHashLabel, appHashEntry,
		layout.NewSpacer(), socksOpen,
		socksIPLabel, socksIPEntry,
		socksPortLabel, socksPortEntry)

	//取消和确认按钮
	cancelAndConfirmButton := container.NewHBox(
		layout.NewSpacer(), layout.NewSpacer(),
		&widget.Button{
			Text:       "取消",
			Importance: widget.DangerImportance,
			OnTapped:   func() { d.Hide() },
		},

		&widget.Button{
			Text:       "确认",
			Importance: widget.HighImportance,
			OnTapped: func() {
				//打开代理
				if socksOpen.Checked {
					if socksIPEntry.Validate() != nil || socksPortEntry.Validate() != nil {
						dialog.NewError(errors.New("代理配置输入有误"), window).Show()
						return
					}
					sIP, _ := socksIP.Get()
					sPort, _ := socksPort.Get()
					config.Set("socksIP", sIP)
					config.Set("socksPort", sPort)
				}

				aid, _ := appId.Get()
				ah, _ := appHash.Get()
				sKey, _ := secretKey.Get()
				config.Set("appId", aid)
				config.Set("apiHash", ah)
				config.Set("secretKey", sKey)

				//codeE.FocusLost() //显示输入框后面的感叹号
				//e := codeE.Validate()
				//if e != nil {
				//	time.Sleep(time.Second * 2)
				//	codeE.FocusGained()
				//	return
				//}
				d.Hide()

			},
		})
	content := container.New(layout.NewVBoxLayout(),
		configContainer,
		cancelAndConfirmButton,
	)
	d = dialog.NewCustomWithoutButtons("配置", content, window)
	d.Resize(fyne.NewSize(300, 300))
	d.Show()
}
