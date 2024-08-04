package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/views/setting"
	"strconv"
)

func getSettingView(window fyne.Window, app fyne.App) *container.TabItem {
	/*configTab := container.NewAppTabs(
		container.NewTabItem("私聊", privateSettingView(app)),
		container.NewTabItem("群聊/频道", groupSettingView(app)),
		container.NewTabItem("配置", setting.GetSettingView(window)),
	)*/
	return container.NewTabItemWithIcon("", theme.SettingsIcon(), setting.GetSettingView(window))
}

// privateSettingView 私聊 标签栏
func privateSettingView(myApp fyne.App) *fyne.Container {
	label := widget.NewLabel("自动给消息点赞")
	radioAlign := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		_ = config.Set("emoticon assistant msg", strconv.FormatBool(s == "是"))
		msg.PrivateRepeatMsg = s == "是"
		label.Refresh()
	})
	radioAlign.Horizontal = true
	parseBool, _ := strconv.ParseBool(config.Get("emoticon assistant msg"))
	msg.PrivateRepeatMsg = parseBool
	if parseBool {
		radioAlign.SetSelected("是")
	} else {
		radioAlign.SetSelected("否")
	}
	return container.NewVBox(label, radioAlign)
}
