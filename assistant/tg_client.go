package assistant

import (
	"fmt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"github.com/midnightsong/telegram-assistant/gotgproto/dispatcher/handlers"
	"github.com/midnightsong/telegram-assistant/gotgproto/dispatcher/handlers/filters"
	"github.com/midnightsong/telegram-assistant/utils"
	"golang.org/x/net/proxy"
	"runtime"
	"strconv"
)

var NewClient = make(chan any)
var Cli *gotgproto.Client

func Run() error {
	config := dao.Config{}
	appid, _ := strconv.Atoi(config.Get("appId"))
	goos := runtime.GOOS
	systemVersion := "23H2"
	if goos == "darwin" {
		goos = "MacBook Pro 2023"
		systemVersion = "macOS 14.5"
	}

	opts := &gotgproto.ClientOpts{
		Session:        dao.SqlSession,
		Logger:         utils.Logger,
		Device:         &telegram.DeviceConfig{DeviceModel: goos, SystemVersion: systemVersion, AppVersion: "1.0 beta18", SystemLangCode: "en", LangPack: "gotgproto", LangCode: "golang"},
		SystemLangCode: "zh_cn",
		ClientLangCode: "zh_cn",
		AutoFetchReply: true,
	}
	//是否开启SOCK5代理
	socksOpen, _ := strconv.ParseBool(config.Get("socksOpen"))
	if socksOpen {
		socksIP := config.Get("socksIP")
		socksPort := config.Get("socksPort")
		sock5, _ := proxy.SOCKS5("tcp", socksIP+":"+socksPort, &proxy.Auth{
			User:     "",
			Password: "",
		}, proxy.Direct)
		dc := sock5.(proxy.ContextDialer)
		opts.Resolver = dcs.Plain(dcs.PlainOptions{Dial: dc.DialContext})
	}

	client, err := gotgproto.NewClient(appid, config.Get("apiHash"), gotgproto.ClientTypePhone(config.Get("phoneNumber")), opts)
	if err != nil {
		NewClient <- err.Error()
		//log.Fatalln("启动客户端失败:", err)
		return err
	}
	dispatcher := client.Dispatcher

	// Command Handler for /start
	//dispatcher.AddHandler(handlers.NewCommand("start", start))
	// Callback Query Handler with prefix filter for recieving specific query
	//dispatcher.AddHandler(handlers.NewCallbackQuery(filters.CallbackQuery.Prefix("cb_"), buttonCallback))
	// This Message Handler will call our echo function on new messages
	//dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, echo), 1)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.ChatType(filters.ChatTypeChat), msg.HandlerGroups), 1) //普通群
	//dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.ChatType(filters.ChatTypeChannel), msg.HandlerGroups), 2) //超级群
	//dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.ChatType(filters.ChatTypeUser), msg.HandlerPrivate), 3)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.ChatType(filters.ChatTypeUser), msg.HandlerGroups), 2)

	fmt.Printf("客户端 (@%s) 已启动...\n", client.Self.Username)
	go func() {
		msg.Client = client
		NewClient <- client
	}()
	client.Idle()
	return nil
}
