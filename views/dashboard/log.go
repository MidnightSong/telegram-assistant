package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/views/icon"
)

var logEntry *widget.RichText

func getLogView(window fyne.Window) *container.TabItem {
	box := container.NewBorder(nil, clearLogButton(), nil, nil, logView())
	return container.NewTabItemWithIcon("", icon.GetIcon(icon.Log), box)
}

var title = `
## 日志
---
`
var logStr = ""

func logView() fyne.CanvasObject {

	logEntry = widget.NewRichTextFromMarkdown(title)
	logEntry.Wrapping = fyne.TextWrapWord
	logEntry.Scroll = container.ScrollBoth
	go func() {
		for {
			log := <-msg.Log
			logStr += "\n\n" + log
			if len(logStr) > 6666 {
				logStr = logStr[len(logStr)-6666:]
			}
			logEntry.ParseMarkdown(title + logStr)
			logEntry.Refresh()
		}
	}()
	return logEntry
}

func clearLogButton() *widget.Button {
	button := widget.NewButtonWithIcon("", icon.GetIcon(icon.Delete), func() {
		logEntry.ParseMarkdown(title)
		logEntry.Refresh()
	})
	button.Importance = widget.LowImportance
	return button
}
