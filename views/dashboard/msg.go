package dashboard

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
	"github.com/midnightsong/telegram-assistant/gotgproto/storage"
	"github.com/midnightsong/telegram-assistant/utils"
	"github.com/midnightsong/telegram-assistant/views/icon"
	"go.uber.org/zap"
	"reflect"
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
	//chats := getDialogs()
	openDialogs := getOpenDialogs()
	chatChecked := map[int]*dialogsInfo{}
	var table *widget.Table
	// 创建保存 Check widget 的二维数组
	checks := make([][]*widget.Check, len(openDialogs)+1)
	col := 3
	for i := range checks {
		checks[i] = make([]*widget.Check, col)
		for j := range checks[i] {
			ck := widget.NewCheck("", nil)
			ck.Hide()
			checks[i][j] = ck
		}
	}
	//表格行列数量
	tableLength := func() (int, int) { return len(openDialogs), col }
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
						checks[i][0].Checked = b
						checks[i][0].Refresh()
						//第0行check是标题栏
						if i == 0 {
							continue
						}
						chatChecked[i-1] = openDialogs[i-1]
					}
					utils.LogInfo(context.Background(), "选中的chat：", zap.Any("map", chatChecked))
					return
				}
				for i := range chatChecked {
					delete(chatChecked, i)
					checks[i][0].Checked = b
					checks[i][0].Refresh()
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
			//因为标题栏的索引是-1，check组件的第1行作为表格的第0行
			check := checks[id.Row+1][id.Col]
			c.Objects[1] = check
			check.OnChanged = func(b bool) {
				if b {
					chatChecked[id.Row] = openDialogs[id.Row]
					utils.LogInfo(context.Background(), "选中的chat：", zap.Any("map", chatChecked))
					return
				}
				delete(chatChecked, id.Row)
				utils.LogInfo(context.Background(), "选中的chat：", zap.Any("map", chatChecked))
			}
			check.Show()
		case 1:
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).SetText(openDialogs[id.Row].title)
			c.Objects[0].(*widget.Label).Show()
		case 2:
			c.Objects[1] = checks[id.Row+1][id.Col]
			c.Objects[0].(*widget.Label).Show()
			if openDialogs[id.Row].EntityType == storage.TypeUser {
				if openDialogs[id.Row].bot {
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
	table.SetColumnWidth(1, 250)
	table.SetColumnWidth(2, 80)
	//名字太长，把行高弄高一点
	/*for index, chat := range chats {
		v := reflect.ValueOf(chat).Elem()
		title := v.FieldByName("Title").String()
		if len(title) > 46 {
			table.SetRowHeight(index, 50)
		}
	}*/
	var refresh *widget.Button
	refresh = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		openDialogs = getOpenDialogs()
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

type dialogsInfo struct {
	title string
	storage.EntityType
	peerId int64
	bot    bool
}

func getOpenDialogs() []*dialogsInfo {
	d, _ := tgClient.API().MessagesGetDialogs(context.Background(), &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{}})
	apiDialogs := reflect.ValueOf(d)
	allChats := apiDialogs.Elem().FieldByName("Chats").Interface().([]tg.ChatClass)
	allUsers := apiDialogs.Elem().FieldByName("Users").Interface().([]tg.UserClass)
	Dialogs := apiDialogs.Elem().FieldByName("Dialogs").Interface().([]tg.DialogClass)
	var dialogsInfos []*dialogsInfo
	for _, i := range Dialogs {
		peerClass := i.GetPeer()
		switch peer := peerClass.(type) {
		case *tg.PeerUser:
			for _, user := range allUsers {
				if u, ok := user.(*tg.User); ok {
					if u.ID == peer.UserID {
						info := &dialogsInfo{
							title:      u.FirstName + u.LastName,
							peerId:     u.ID,
							EntityType: storage.TypeUser,
							bot:        u.Bot,
						}
						dialogsInfos = append(dialogsInfos, info)
						break
					}
				}
			}
		case *tg.PeerChat:
			for _, chat := range allChats {
				if c, ok := chat.(*tg.Chat); ok {
					if c.ID == peer.ChatID {
						info := &dialogsInfo{
							title:      c.Title,
							peerId:     c.ID,
							EntityType: storage.TypeChat,
						}
						dialogsInfos = append(dialogsInfos, info)
						break
					}
				}
			}
		case *tg.PeerChannel:
			for _, chat := range allChats {
				if c, ok := chat.(*tg.Channel); ok {
					if c.ID == peer.ChannelID {
						info := &dialogsInfo{
							title:      c.Title,
							peerId:     c.ID,
							EntityType: storage.TypeChannel,
						}
						dialogsInfos = append(dialogsInfos, info)
						break
					}
				}
			}
		}
	}
	return dialogsInfos
}
