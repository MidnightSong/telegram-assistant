package setting

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/dao"
	"strconv"
)

var config = dao.Config{}
var ipAddressRegex = "((25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]?\\d)\\.){3}(25[0-5]|2[0-4]\\d|1\\d{2}|[1-9]?\\d)"
var numberRegex = "(\\d{1,5})"

func GetSettingView(window fyne.Window) *fyne.Container {

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

	confirmButton := &widget.Button{
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

			dialog.NewInformation("成功", "保存配置成功", window).Show()
			//codeE.FocusLost() //显示输入框后面的感叹号
			//e := codeE.Validate()
			//if e != nil {
			//	time.Sleep(time.Second * 2)
			//	codeE.FocusGained()
			//	return
			//}

		},
	}
	//取消和确认按钮
	cancelAndConfirmButton := container.NewHBox(
		layout.NewSpacer(), layout.NewSpacer(), confirmButton)
	return container.New(layout.NewVBoxLayout(), configContainer, cancelAndConfirmButton)
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
