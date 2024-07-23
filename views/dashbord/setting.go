package dashbord

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"image/color"
	"strconv"
)

func getSettingView(app fyne.App) *container.TabItem {
	configTab := container.NewAppTabs(
		container.NewTabItem("私聊", privateSettingView(app)),
		container.NewTabItem("群聊/频道", groupSettingView(app)),
	)
	return container.NewTabItemWithIcon("", theme.SettingsIcon(), container.NewVBox(configTab))
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

// groupSettingView 群聊/频道 标签栏
func groupSettingView(myApp fyne.App) *fyne.Container {
	//重复机器人的消息
	repeatBotMsg := widget.NewLabel("重复机器人的消息")
	repeatBotMsgAlign := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		_ = config.Set("GroupRepeatMsg", strconv.FormatBool(s == "是"))
		msg.GroupRepeatMsg = s == "是"
		repeatBotMsg.Refresh()
	})
	repeatBotMsgAlign.Horizontal = true
	parseBool, _ := strconv.ParseBool(config.Get("GroupRepeatMsg"))
	msg.GroupRepeatMsg = parseBool
	if parseBool {
		repeatBotMsgAlign.SetSelected("是")
	} else {
		repeatBotMsgAlign.SetSelected("否")
	}

	//是否隐藏重复消息的来源
	hideRepeatBotMsg := widget.NewLabel("当重复消息时，是否隐藏来源")
	hideRepeatBotMsgAlign := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		_ = config.Set("GroupHideSourceRepeatBotMsg", strconv.FormatBool(s == "是"))
		msg.GroupHideSourceRepeatBotMsg = s == "是"
		hideRepeatBotMsg.Refresh()
	})
	hideRepeatBotMsgAlign.Horizontal = true
	parseBool, _ = strconv.ParseBool(config.Get("GroupHideSourceRepeatBotMsg"))
	msg.GroupHideSourceRepeatBotMsg = parseBool
	if parseBool {
		hideRepeatBotMsgAlign.SetSelected("是")
	} else {
		hideRepeatBotMsgAlign.SetSelected("否")
	}

	//关联回复重复过的机器人消息
	groupRepeatMsgReplyTo := widget.NewLabel("当有回复自己的消息时，关联回复重复过的机器人消息")
	groupRepeatMsgReplyTo2 := widget.NewLabel("(仅当显示消息来源时生效)")
	groupRepeatMsgReplyToAlign := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		_ = config.Set("GroupRepeatMsgReplyTo", strconv.FormatBool(s == "是"))
		msg.GroupRepeatMsgReplyTo = s == "是"
	})
	groupRepeatMsgReplyToAlign.Horizontal = true
	parseBool, _ = strconv.ParseBool(config.Get("GroupRepeatMsgReplyTo"))
	msg.GroupRepeatMsgReplyTo = parseBool
	if parseBool {
		groupRepeatMsgReplyToAlign.SetSelected("是")
	} else {
		groupRepeatMsgReplyToAlign.SetSelected("否")
	}
	lineH := canvas.NewLine(color.Black)
	repeatBotBox := container.NewVBox(repeatBotMsg, repeatBotMsgAlign)
	hideRepeatBotBox := container.NewVBox(lineH, hideRepeatBotMsg, hideRepeatBotMsgAlign)
	groupRepeatMsgReplyToBox := container.NewVBox(lineH, groupRepeatMsgReplyTo, groupRepeatMsgReplyTo2, groupRepeatMsgReplyToAlign)
	return container.NewVBox(repeatBotBox, hideRepeatBotBox, groupRepeatMsgReplyToBox)
}
