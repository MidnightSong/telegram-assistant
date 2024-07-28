package dashboard

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/gotgproto/storage"
	"github.com/midnightsong/telegram-assistant/views/icon"
	"time"
)

func getMsgView(window fyne.Window) *container.TabItem {
	msgTab := container.NewAppTabs(
		container.NewTabItem("处理日志", msgRtView()),
		container.NewTabItem("已打开的会话", getOpenedDialogs(window)),
	)

	return container.NewTabItemWithIcon("", icon.GetIcon(icon.Telegram), msgTab)
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

// 已选中的会话
var chatChecked map[int]*assistant.DialogsInfo

// openedDialogs 返回已打开的会话视图（表格）
func getOpenedDialogs(window fyne.Window) fyne.CanvasObject {
	chatChecked = map[int]*assistant.DialogsInfo{}
	var table *widget.Table
	checks := map[int]*widget.Check{}
	col := 3
	//表格行列数量
	tableLength := func() (int, int) { return len(openedDialogs), col }
	//初始化单元格
	initCell := func() fyne.CanvasObject {
		label := widget.NewLabel("")
		//文本末尾不够展示时，加入省略号
		label.Truncation = fyne.TextTruncateEllipsis
		check := widget.NewCheck("", nil)
		check.Hide()
		return container.NewStack(label, check)
	}
	//标题栏更新 标题栏的id.Row索引是-1
	updateHeader := func(id widget.TableCellID, cell fyne.CanvasObject) {
		c := cell.(*fyne.Container)
		switch id.Col {
		case 0:
			c.Objects[0].(*widget.Label).Hide()
			allCheck := c.Objects[1].(*widget.Check)
			allCheck.Show()
			if checks[id.Row] == nil {
				checks[id.Row] = allCheck
			}
			allCheck.OnChanged = func(b bool) {
				if b {
					for i := range checks {
						if i == -1 {
							continue
						}
						checks[i].Checked = b
						checks[i].Refresh()
						chatChecked[i] = openedDialogs[i]
					}
					return
				}
				for i := range chatChecked {
					if i == -1 {
						continue
					}
					delete(chatChecked, i)
					checks[i].Checked = b
					checks[i].Refresh()
				}
			}
		case 1:
			c.Objects[0].(*widget.Label).SetText("标题")
			c.Objects[0].(*widget.Label).Show()
		case 2:
			c.Objects[0].(*widget.Label).SetText("类型")
			c.Objects[0].(*widget.Label).Show()
		}
		return
	}
	//单元格更新
	updateCell := func(id widget.TableCellID, cell fyne.CanvasObject) {
		c := cell.(*fyne.Container)
		switch id.Col {
		case 0:
			c.Objects[0].(*widget.Label).Hide()
			check := c.Objects[1].(*widget.Check)
			if checks[id.Row] == nil {
				checks[id.Row] = check
			}
			check.OnChanged = func(b bool) {
				if b {
					chatChecked[id.Row] = openedDialogs[id.Row]
					return
				}
				delete(chatChecked, id.Row)
			}
			check.Show()
		case 1:
			c.Objects[0].(*widget.Label).SetText(openedDialogs[id.Row].Title)
			c.Objects[0].(*widget.Label).Show()
		case 2:
			c.Objects[0].(*widget.Label).Show()
			if openedDialogs[id.Row].EntityType == storage.TypeUser {
				if openedDialogs[id.Row].Bot {
					c.Objects[0].(*widget.Label).SetText("机器人")
					return
				}
				c.Objects[0].(*widget.Label).SetText("用户")
				return
			}
			c.Objects[0].(*widget.Label).SetText("群组/频道")
		}
	}

	table = widget.NewTable(tableLength, initCell, updateCell)
	table.CreateHeader = initCell
	table.UpdateHeader = updateHeader
	table.ShowHeaderRow = true
	table.SetColumnWidth(0, 30)
	table.SetColumnWidth(1, 150)
	table.SetColumnWidth(2, 80)
	//名字太长，把行高弄高一点
	/*for index, chat := range chats {
		v := reflect.ValueOf(chat).Elem()
		title := v.FieldByName("Title").String()
		if len(title) > 46 {
			table.SetRowHeight(index, 50)
		}
	}*/
	//刷新已打开会话的按钮
	var refresh *widget.Button
	refresh = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		//清除当前checks除了标题栏
		table.Refresh()
		go func() {
			refresh.Disable()
			time.Sleep(time.Second * 5)
			refresh.Enable()
		}()
	})
	var loudSpeaker *widget.Button
	loudSpeaker = widget.NewButtonWithIcon("", icon.GetIcon(icon.LoudSpeaker), func() {
		ShowSendMsgModal(window)
	})
	buttonsBox := container.NewHBox(loudSpeaker, layout.NewSpacer(), refresh)
	return container.NewBorder(buttonsBox, nil, nil, nil, table)
}

func ShowSendMsgModal(window fyne.Window) {
	var up *widget.PopUp
	title := widget.NewLabel("向已选中的会话群发消息")
	closeButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
		up.Hide()
	})
	closeButton.Importance = widget.LowImportance
	inputMsg := widget.NewMultiLineEntry()
	sendMsgActivity := widget.NewActivity()
	var returnMsgButton *widget.Button
	var sendMsgButton *widget.Button

	//init buttons
	var sentMsgId map[int64]int //保存上一次发送的消息的id
	returnMsgButton = widget.NewButtonWithIcon("撤回", icon.GetIcon(icon.TelegramReturn), func() {
		dialog.ShowConfirm("撤回", "确定要撤回上一次发出的所有消息吗？", func(yes bool) {
			if yes {
				for chatId, msgId := range sentMsgId {
					err := msg.DeleteMessage(chatId, msgId)
					if err != nil {
						msg.AddLog(fmt.Sprintf("撤回消息失败：%s", err.Error()))
						continue
					}
				}
				returnMsgButton.Disable()
			}
		}, window)
	})
	returnMsgButton.Importance = widget.DangerImportance
	returnMsgButton.Disable()
	sendMsgButton = widget.NewButtonWithIcon("发送", icon.GetIcon(icon.Telegram), func() {
		if inputMsg.Text == "" {
			return
		}
		sendMsgActivity.Start()
		sendMsgActivity.Show()
		sendMsgButton.Disable()
		sentMsgId = make(map[int64]int)
		for _, chat := range chatChecked {
			msgId, err := msg.SendMessage(chat.PeerId, inputMsg.Text)
			if err != nil {
				msg.AddLog(err.Error())
				continue
			}
			sentMsgId[chat.PeerId] = msgId
		}
		sendMsgActivity.Stop()
		sendMsgActivity.Hide()
		sendMsgButton.Enable()
		returnMsgButton.Enable()
		dialog.ShowInformation("", "消息已发送完毕", window)
	})
	sendMsgButton.Importance = widget.HighImportance

	topBox := container.NewHBox(title, layout.NewSpacer(), closeButton)
	bottomBox := container.NewHBox(returnMsgButton, layout.NewSpacer(), container.NewStack(sendMsgActivity, sendMsgButton))
	layoutBox := container.NewBorder(topBox, bottomBox, nil, nil, inputMsg)
	up = widget.NewModalPopUp(layoutBox, window.Canvas())
	up.Show()
	up.Resize(fyne.NewSize(300, 300))
}
