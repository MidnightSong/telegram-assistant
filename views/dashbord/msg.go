package dashbord

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/icon"
	"github.com/midnightsong/telegram-assistant/utils"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

func getMsgView() *container.TabItem {
	msgTab := container.NewAppTabs(
		container.NewTabItem("处理日志", msgRtView()),
		container.NewTabItem("已打开的会话", openedDialogs()),
	)
	return container.NewTabItemWithIcon("", theme.MailSendIcon(), msgTab)
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

func openedDialogs() fyne.CanvasObject {
	chats := getDialogs()
	chatChecked := map[int]interface{}{}
	var table *widget.Table
	// 创建保存 Check widget 的二维数组
	checks := make([][]*widget.Check, len(chats)+1)
	for i := range checks {
		checks[i] = make([]*widget.Check, 4)
		for j := range checks[i] {
			ck := widget.NewCheck("", nil)
			ck.Hide()
			checks[i][j] = ck
		}
	}
	//表格行列数量
	tableLength := func() (int, int) { return len(chats), 4 }
	//初始化单元格
	initCell := func() fyne.CanvasObject {
		label := widget.NewLabel("")
		label.Wrapping = fyne.TextWrapWord
		return container.NewStack(label, widget.NewCheck("", nil))
	}
	//标题栏更新 标题栏的id.Row索引是-1
	updateHeader := func(id widget.TableCellID, cell fyne.CanvasObject) {
		c := cell.(*fyne.Container)
		switch id.Col {
		case 0:
			allCheck := checks[0][0]
			allCheck.OnChanged = func(b bool) {
				if b {
					for i := range checks {
						/*for j := range checks[i] {
							checks[i][j].Checked = b
							checks[i][j].Refresh()
						}*/
						checks[i][0].Checked = b
						checks[i][0].Refresh()
						//第0行check是标题栏
						if i == 0 {
							continue
						}
						chatChecked[i-1] = chats[i-1]
					}
					utils.LogInfo(context.Background(), "选中的chat：", zap.Any("map", chatChecked))
					return
				}
				for i := range chatChecked {
					delete(chatChecked, i)
				}
				utils.LogInfo(context.Background(), "选中的chat：", zap.Any("map", chatChecked))
			}
			allCheck.Show()
			c.Objects[1] = allCheck
		case 1:
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).SetText("标题")
			c.Objects[0].(*widget.Label).Show()
		case 2:
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).SetText("ID")
			c.Objects[0].(*widget.Label).Show()
		case 3:
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).SetText("类型")
			c.Objects[0].(*widget.Label).Show()
		}
		return
	}
	//单元格更新
	updateCell := func(id widget.TableCellID, cell fyne.CanvasObject) {
		c := cell.(*fyne.Container)
		t := reflect.TypeOf(chats[id.Row]).String()
		v := reflect.ValueOf(chats[id.Row]).Elem()
		switch id.Col {
		case 0:
			c.Objects[0].(*widget.Label).Hide()
			//因为标题栏的索引是-1，check组件的第1行作为表格的第0行
			check := checks[id.Row+1][id.Col]
			c.Objects[1] = check
			check.OnChanged = func(b bool) {
				if b {
					chatChecked[id.Row] = chats[id.Row]
					utils.LogInfo(context.Background(), "选中的chat：", zap.Any("map", chatChecked))
					return
				}
				delete(chatChecked, id.Row)
				//fmt.Printf("Check Row:%d Col: %d \n", id.Row+1, id.Col)
				utils.LogInfo(context.Background(), "选中的chat：", zap.Any("map", chatChecked))
			}
			check.Show()
		case 1:
			var title string
			if strings.Contains(t, "User") {
				title = v.FieldByName("FirstName").String() + " " + v.FieldByName("LastName").String()
			} else {
				title = v.FieldByName("Title").String()
			}
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).SetText(title)
			c.Objects[0].(*widget.Label).Show()
		case 2:
			ID := v.FieldByName("ID")
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).SetText(fmt.Sprint(ID.Int()))
			c.Objects[0].(*widget.Label).Show()
		case 3:
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).Show()
			if strings.Contains(t, "User") {
				isBot := v.FieldByName("Bot").Bool()
				if isBot {
					c.Objects[0].(*widget.Label).SetText("机器人")
				} else {
					c.Objects[0].(*widget.Label).SetText("用户")
				}
			} else {
				c.Objects[0].(*widget.Label).SetText("群组/频道")
			}
		}
	}
	table = widget.NewTable(tableLength, initCell, updateCell)
	table.CreateHeader = initCell
	table.UpdateHeader = updateHeader
	table.ShowHeaderRow = true
	table.SetColumnWidth(1, 250)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 80)
	//名字太长，把行高弄高一点
	for index, chat := range chats {
		v := reflect.ValueOf(chat).Elem()
		title := v.FieldByName("Title").String()
		if len(title) > 46 {
			table.SetRowHeight(index, 50)
		}
	}
	var refresh *widget.Button
	refresh = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		chats = getDialogs()
		table.Refresh()
		go func() {
			refresh.Disable()
			time.Sleep(time.Second * 5)
			refresh.Enable()
		}()
	})
	var loudSpeaker *widget.Button
	loudSpeaker = widget.NewButtonWithIcon("", icon.GetIcon(icon.LoudSpeaker), func() {
		fmt.Println("点了一下广播按钮")
	})
	buttonsBox := container.NewHBox(loudSpeaker, layout.NewSpacer(), refresh)
	return container.NewBorder(buttonsBox, nil, nil, nil, table)
}

func getDialogs() []interface{} {
	d, _ := tgClient.API().MessagesGetDialogs(context.Background(), &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{}})
	apiDialogs := reflect.ValueOf(d)
	allChats := apiDialogs.Elem().FieldByName("Chats").Interface().([]tg.ChatClass)
	//allChats := d.(*tg.MessagesDialogsSlice).Chats
	allUsers := apiDialogs.Elem().FieldByName("Users").Interface().([]tg.UserClass)
	//allUsers := d.(*tg.MessagesDialogsSlice).Users
	var allDialogs []interface{}
	for _, i := range allChats {
		allDialogs = append(allDialogs, i)
	}
	for _, i := range allUsers {
		allDialogs = append(allDialogs, i)
	}
	//排除已被禁止的群聊或频道
	chats := allDialogs[:0]
	for i := 0; i < len(allDialogs); i++ {
		typeName := reflect.TypeOf(allDialogs[i]).String()
		if !strings.Contains(typeName, "Forbidden") {
			chats = append(chats, allDialogs[i])
		}
	}
	return chats
}
