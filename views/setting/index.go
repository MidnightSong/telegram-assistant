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
	/*dir, _ := os.Getwd()
	information := dialog.NewInformation("path", dir, window)
	information.Resize(fyne.NewSize(300, 300))
	information.Show()*/
	authCode := binding.NewString()
	_ = authCode.Set(config.Get("authCode"))
	authCodeLabel := widget.NewLabel("激活码")
	authCodeEntry := widget.NewEntryWithData(authCode)
	authCodeEntry.Password = true
	authCodeEntry.Validator = validation.NewRegexp("\\w+", "激活码不能为空！")

	appId := binding.NewString()
	_ = appId.Set(config.Get("appId"))
	appIdLabel := widget.NewLabel("appId")
	appIdEntry := newAppIDEntry(appId)

	appHash := binding.NewString()
	_ = appHash.Set(config.Get("apiHash"))
	appHashLabel := widget.NewLabel("apiHash")
	appHashEntry := widget.NewEntryWithData(appHash)
	appHashEntry.Password = true

	socksIP := binding.NewString()
	_ = socksIP.Set(config.Get("socksAddr"))
	socksIPLabel := widget.NewLabel("IP:")
	socksIPEntry := widget.NewEntry()
	socksIPEntry.PlaceHolder = "例: 127.0.0.1"
	socksIPEntry.Bind(socksIP)
	socksIPEntry.Validator = validation.NewRegexp(ipAddressRegex, "请输入正确的ip地址！")

	socksPort := binding.NewString()
	_ = socksPort.Set(config.Get("socksPort"))
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
		_ = config.Set("socksOpen", fmt.Sprint(b))
	}
	socksOpen := widget.NewCheck("socks5 代理", checkSocks)

	parseBool, _ := strconv.ParseBool(config.Get("socksOpen"))
	socksOpen.SetChecked(parseBool)
	checkSocks(parseBool)
	socksOpen.Refresh()

	configContainer := container.New(layout.NewFormLayout(),
		authCodeLabel, authCodeEntry,
		appIdLabel, appIdEntry,
		appHashLabel, appHashEntry,
		layout.NewSpacer(), socksOpen,
		socksIPLabel, socksIPEntry,
		socksPortLabel, socksPortEntry)

	confirmButton := &widget.Button{
		Text:       "保存",
		Importance: widget.HighImportance,
		OnTapped: func() {
			if err := authCodeEntry.Validate(); err != nil {
				dialog.NewError(err, window).Show()
				return
			}
			if err := appIdEntry.Validate(); err != nil {
				dialog.NewError(err, window).Show()
				return
			}
			if err := appHashEntry.Validate(); err != nil {
				dialog.NewError(err, window).Show()
				return
			}

			aid, _ := appId.Get()
			ah, _ := appHash.Get()
			sKey, _ := authCode.Get()
			config.Set("appId", aid)
			config.Set("apiHash", ah)
			config.Set("authCode", sKey)
			//打开代理
			if socksOpen.Checked {
				if socksIPEntry.Validate() != nil || socksPortEntry.Validate() != nil {
					dialog.NewError(errors.New("代理配置输入有误"), window).Show()
					return
				}
				sIP, _ := socksIP.Get()
				sPort, _ := socksPort.Get()
				config.Set("socksAddr", sIP)
				config.Set("socksPort", sPort)
			}
			dialog.ShowInformation("成功", "保存配置成功,配置将在重启客户端后生效", window)
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
		layout.NewSpacer(), confirmButton, layout.NewSpacer())
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
