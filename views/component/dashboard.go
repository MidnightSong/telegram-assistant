package component

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/bot"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"github.com/midnightsong/telegram-assistant/msg"
	"github.com/midnightsong/telegram-assistant/utils"
	"go.uber.org/zap"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var config = dao.Config{}
var tgClient *gotgproto.Client

func MsgNewWindow(window fyne.Window, myApp fyne.App) {
	tgClient = <-bot.NewClient
	gotgproto.Logged = true
	dashboardWindow := myApp.NewWindow(fmt.Sprintf("欢迎：%s %s", tgClient.Self.FirstName, tgClient.Self.LastName))

	msgTab := container.NewAppTabs(
		container.NewTabItem("处理日志", msgRtView()),
		container.NewTabItem("已打开的会话", openedDialogs()),
	)

	configTab := container.NewAppTabs(
		container.NewTabItem("私聊", configPrivateView(myApp)),
		container.NewTabItem("群聊/频道", configGroupView(myApp)),
	)

	leftTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("消息", theme.MailSendIcon(), msgTab),
		container.NewTabItemWithIcon("配置", theme.SettingsIcon(), container.NewVBox(configTab)),
	)

	leftTabs.SetTabLocation(container.TabLocationLeading)

	dashboardWindow.Resize(fyne.NewSize(1024, 576))
	dashboardWindow.SetContent(leftTabs)
	dashboardWindow.Show()
	dashboardWindow.SetCloseIntercept(func() {
		dashboardWindow.Close()
		myApp.Quit()
		os.Exit(0)
	})
	window.Close()
}

func configGroupView(myApp fyne.App) *fyne.Container {
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
	groupRepeatMsgReplyToAlign := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		_ = config.Set("GroupRepeatMsgReplyTo", strconv.FormatBool(s == "是"))
		msg.GroupRepeatMsgReplyTo = s == "是"
		groupRepeatMsgReplyTo.Refresh()
	})
	groupRepeatMsgReplyToAlign.Horizontal = true
	parseBool, _ = strconv.ParseBool(config.Get("GroupRepeatMsgReplyTo"))
	msg.GroupRepeatMsgReplyTo = parseBool
	if parseBool {
		groupRepeatMsgReplyToAlign.SetSelected("是")
	} else {
		groupRepeatMsgReplyToAlign.SetSelected("否")
	}
	return container.NewVBox(repeatBotMsg, repeatBotMsgAlign, hideRepeatBotMsg, hideRepeatBotMsgAlign, groupRepeatMsgReplyTo, groupRepeatMsgReplyToAlign)
}

func configPrivateView(myApp fyne.App) *fyne.Container {
	label := widget.NewLabel("自动给消息点赞")
	radioAlign := widget.NewRadioGroup([]string{"是", "否"}, func(s string) {
		_ = config.Set("emoticon bot msg", strconv.FormatBool(s == "是"))
		msg.PrivateRepeatMsg = s == "是"
		label.Refresh()
	})
	radioAlign.Horizontal = true
	parseBool, _ := strconv.ParseBool(config.Get("emoticon bot msg"))
	msg.PrivateRepeatMsg = parseBool
	if parseBool {
		radioAlign.SetSelected("是")
	} else {
		radioAlign.SetSelected("否")
	}
	return container.NewVBox(label, radioAlign)
}
func msgRtView() fyne.CanvasObject {
	RtMsg := widget.NewMultiLineEntry()
	RtMsg.Wrapping = fyne.TextWrapWord
	go func() {
		for {
			log := <-msg.GroupLog
			RtMsg.Text += "\n" + log
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

	refresh := widget.NewButton("刷新", func() {
		chats = getDialogs()
		table.Refresh()
	})
	return container.NewBorder(refresh, nil, nil, nil, table)
}

func getDialogs() []interface{} {
	fmt.Println("=======")
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
