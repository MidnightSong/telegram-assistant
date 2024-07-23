package views

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"time"
)

// 要求输入验证码
func showVerify(app fyne.App) {
	for {
		time.Sleep(1 * time.Second)
		if gotgproto.Logged {
			break
		}
		if !gotgproto.AskVerifyCode {
			continue
		}
		gotgproto.AskVerifyCode = false
		window := app.NewWindow("验证码")
		verifyCode := newVerifyEntry()
		c := container.New(layout.NewFormLayout(), widget.NewLabel("请输入验证码："), verifyCode)
		button := widget.NewButton("确认", func() {
			err := verifyCode.Validate()
			if err != nil {
				dialog.NewError(err, window).Show()
				return
			}
			window.Close()
			gotgproto.AuthCode <- verifyCode.Text
		})
		window.SetContent(container.NewVBox(c, button))
		window.Resize(fyne.NewSize(300, 200))
		window.Show()
	}
}

// 要求输入两步验证密码
func showPassword(app fyne.App) {
	for {
		if gotgproto.Logged {
			break
		}
		if !gotgproto.AskPassword {
			time.Sleep(1 * time.Second)
			continue
		}
		gotgproto.AskPassword = false
		window := app.NewWindow("两步验证密码")
		fmt.Println("两步验证密码")
		verifyCode := widget.NewEntry()
		verifyCode.Password = true
		c := container.New(layout.NewFormLayout(), widget.NewLabel("请输入两步验证密码："), verifyCode)
		button := widget.NewButton("确认", func() {
			window.Close()
			gotgproto.AuthCode <- verifyCode.Text
		})
		window.SetContent(container.NewVBox(c, button))
		window.Resize(fyne.NewSize(300, 200))
		window.Show()
	}
}

type verifyEntry struct {
	widget.Entry
}

func newVerifyEntry() *verifyEntry {
	p := &verifyEntry{}
	p.ExtendBaseWidget(p)
	p.Validator = validation.NewRegexp(`^\d+$`, "必须是纯数字")
	return p
}
