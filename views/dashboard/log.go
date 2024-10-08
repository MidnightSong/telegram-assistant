package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/views/icon"
	"strconv"
)

var logEntry *widget.RichText
var title = `
## 日志
---
`
var logStr string
var logBox fyne.CanvasObject

func getLogView(window fyne.Window) *container.TabItem {
	bottomBox := container.NewHBox(clearLogButton(), layout.NewSpacer(), switchLog())
	box := container.NewBorder(nil, bottomBox, nil, nil, logView())
	return container.NewTabItemWithIcon("", icon.GetIcon(icon.Log), box)
}

func logView() fyne.CanvasObject {
	logEntry = widget.NewRichTextFromMarkdown(title)
	logEntry.Wrapping = fyne.TextWrapWord
	go func() {
		for {
			log := <-msg.Log
			logStr += log
			if len(logStr) > 66666 {
				logStr = logStr[len(logStr)-66666:]
			}
			logEntry.ParseMarkdown(title + logStr)
			logBox.(*container.Scroll).ScrollToBottom()
			logEntry.Refresh()
		}
	}()
	logBox = container.NewVScroll(logEntry)
	return logBox
}

func clearLogButton() *widget.Button {
	button := widget.NewButtonWithIcon("", icon.GetIcon(icon.Delete), func() {
		logEntry.ParseMarkdown(title)
		logEntry.Refresh()
		logStr = ""
	})
	button.Importance = widget.LowImportance
	return button
}
func switchLog() *widget.Button {
	value := config.Get("logSwitch")
	v, err := strconv.ParseBool(value)
	var buttonName string
	if err != nil || v == false {
		msg.LogSwitch = false
		buttonName = "打开日志"
	} else {
		msg.LogSwitch = v
		buttonName = "关闭日志"
	}
	var button *widget.Button
	button = widget.NewButton(buttonName, func() {
		if msg.LogSwitch {
			msg.LogSwitch = false
			_ = config.Set("logSwitch", "false")
			button.SetText("打开日志")
		} else {
			msg.LogSwitch = true
			_ = config.Set("logSwitch", "true")
			button.SetText("关闭日志")
		}
	})
	button.Importance = widget.LowImportance
	return button
}
