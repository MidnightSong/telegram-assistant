package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/views/icon"
)

// TODO 清空日志信息
func getLogView(window fyne.Window) *container.TabItem {
	return container.NewTabItemWithIcon("", icon.GetIcon(icon.Log), msgRtView())
}

func msgRtView() fyne.CanvasObject {
	RtMsg := widget.NewMultiLineEntry()
	RtMsg.Wrapping = fyne.TextWrapWord
	go func() {
		for {
			log := <-msg.Log
			RtMsg.Text += "\n" + log
			if len(RtMsg.Text) > 6666 {
				RtMsg.Text = RtMsg.Text[len(RtMsg.Text)-6666:]
			}
			RtMsg.Refresh()
		}
	}()
	return RtMsg
}
