package dashboard

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"github.com/midnightsong/telegram-assistant/gotgproto/storage"
	"github.com/midnightsong/telegram-assistant/views/icon"
	"os"
	"reflect"
	"time"
)

var config = dao.Config{}
var cli any
var tgClient *gotgproto.Client
var openedDialogs []*msg.DialogsInfo

func MsgNewWindow(jumpInWindow fyne.Window, myApp fyne.App) {
	cli = <-assistant.NewClient
	switch c := cli.(type) {
	case *gotgproto.Client:
		tgClient = c
	case string:
		logFail := myApp.NewWindow("登录失败")
		errContent := widget.NewRichTextWithText(c)
		errContent.Wrapping = fyne.TextWrapWord
		errContent.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignCenter

		shutDownButton := widget.NewButtonWithIcon("注销", icon.ShutDown, func() {
			dialog.NewConfirm("注销？", "注销将删除当前账号所有数据，不可恢复。", func(b bool) {
				if b {
					_ = sessions.DeleteAll()
					_ = peers.DeleteAll()
					os.Exit(0)
				}
			}, logFail).Show()
		})

		shutDownButton.Importance = widget.WarningImportance
		exitButton := widget.NewButtonWithIcon("退出", theme.CancelIcon(), func() {
			os.Exit(0)
		})
		exitButton.Importance = widget.HighImportance
		box := container.NewVBox(errContent, container.NewHBox(exitButton, layout.NewSpacer(), shutDownButton))
		logFail.Resize(fyne.NewSize(300, 400))
		logFail.SetContent(box)
		logFail.Show()
		logFail.CenterOnScreen()
		logFail.SetFixedSize(true)
		jumpInWindow.Close()
		return
	}

	gotgproto.Logged = true
	dashboardWindow := myApp.NewWindow(fmt.Sprintf("欢迎：%s %s", tgClient.Self.FirstName, tgClient.Self.LastName))

	leftTabs := container.NewAppTabs(
		getMsgView(dashboardWindow),            //消息栏
		getSettingView(dashboardWindow, myApp), //设置栏
		GetLogOutView(dashboardWindow),         //注销栏
		getForwardView(dashboardWindow),        //搬运
	)
	leftTabs.SetTabLocation(container.TabLocationLeading)
	openedDialogs = refreshOpenedDialogs()
	go func() {
		for {
			time.Sleep(10 * time.Second)
			openedDialogs = refreshOpenedDialogs()
		}
	}()
	dashboardWindow.Resize(fyne.NewSize(1024, 576))
	dashboardWindow.SetContent(leftTabs)
	dashboardWindow.CenterOnScreen()
	dashboardWindow.Show()
	jumpInWindow.Close()
	dashboardWindow.SetCloseIntercept(func() {
		dashboardWindow.Close()
		myApp.Quit()
		os.Exit(0)
	})
	msg.AddLog("当前缓存文件路径为：" + myApp.Storage().RootURI().String())
}

func refreshOpenedDialogs() []*msg.DialogsInfo {
	d, e := tgClient.API().MessagesGetDialogs(context.Background(), &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{}})
	if e != nil {
		msg.AddLog("更新已打开的会话列表失败：" + e.Error())
	}
	apiDialogs := reflect.ValueOf(d)
	allChats := apiDialogs.Elem().FieldByName("Chats").Interface().([]tg.ChatClass)
	allUsers := apiDialogs.Elem().FieldByName("Users").Interface().([]tg.UserClass)
	Dialogs := apiDialogs.Elem().FieldByName("Dialogs").Interface().([]tg.DialogClass)
	var dialogsInfos []*msg.DialogsInfo
	for _, i := range Dialogs {
		peerClass := i.GetPeer()
		switch peer := peerClass.(type) {
		case *tg.PeerUser:
			for _, user := range allUsers {
				if u, ok := user.(*tg.User); ok {
					if u.ID == peer.UserID {
						info := &msg.DialogsInfo{
							Title:      u.FirstName + u.LastName,
							PeerId:     u.ID,
							EntityType: storage.TypeUser,
							Bot:        u.Bot,
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
						info := &msg.DialogsInfo{
							Title:      c.Title,
							PeerId:     c.ID,
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
						info := &msg.DialogsInfo{
							Title:      c.Title,
							PeerId:     c.ID,
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
