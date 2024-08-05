package msg

import (
	"context"
	"fmt"
	"github.com/gotd/td/tg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/gotgproto"
	"github.com/midnightsong/telegram-assistant/gotgproto/storage"
	"github.com/midnightsong/telegram-assistant/utils"
	"reflect"
	"regexp"
	"sync"
	"time"
)

var CliChan = make(chan *gotgproto.Client)
var Client *gotgproto.Client
var defaultReg = regexp.MustCompile(`\w{8,}`)
var fr = &dao.ForwardRelation{}
var CacheRelationsMap = sync.Map{}
var fwdMsgDao = dao.FwdMsg{}
var OpenedDialogs = make([]*DialogsInfo, 0) //定时刷新已打开的会话

// Init 使用init函数初始化变量会导致App启动异常
func Init() {
	// 查出所有的绑定关系集合，然后通过源id进行分类后缓存
	refreshCache := func() {
		all := fr.All()
		CacheRelationsMap = sync.Map{}
		for _, d := range all {
			key := d.PeerID
			var relations []*entities.ForwardRelation
			value, _ := CacheRelationsMap.Load(key)
			if value != nil { //TODO bug fix
				relations = value.([]*entities.ForwardRelation)
			} else {
				relations = []*entities.ForwardRelation{}
			}
			relations = append(relations, d)
			CacheRelationsMap.Store(key, relations)
		}
	}

	//获取当前账号已打开的会话
	go func() {
		for {
			if Client == nil {
				Client = <-CliChan
			}
			refreshCache()
			OpenedDialogs = refreshOpenedDialogs()
			//如果绑定关系中有记录，但已打开会话中没有
			//说明该会话已被关闭，那么对应删除绑定关系
			if len(OpenedDialogs) != 0 {
				CacheRelationsMap.Range(func(k, v interface{}) bool {
					flag := true
					for _, dialog := range OpenedDialogs {
						if dialog.PeerId == k {
							flag = false
							break
						}
					}
					if flag {
						CacheRelationsMap.Delete(k)
						fr.Delete(k.(int64))
						CacheRelationsMap.Delete(k)
					}
					return true
				})
			}
			time.Sleep(30 * time.Second)
		}
	}()
}
func refreshOpenedDialogs() []*DialogsInfo {
	defer func() {
		if err := recover(); err != nil {
			log := fmt.Sprintf("更新已打开的会话列表异常：%v", err)
			AddLog(log)
		}
	}()
	d, e := Client.API().MessagesGetDialogs(context.Background(), &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{}})
	if e != nil {
		AddLog("更新已打开的会话列表失败：" + e.Error())
		utils.LogError(context.Background(), "更新已打开的会话列表失败："+e.Error())
		return nil
	}
	apiDialogs := reflect.ValueOf(d)
	allChats := apiDialogs.Elem().FieldByName("Chats").Interface().([]tg.ChatClass)
	allUsers := apiDialogs.Elem().FieldByName("Users").Interface().([]tg.UserClass)
	Dialogs := apiDialogs.Elem().FieldByName("Dialogs").Interface().([]tg.DialogClass)
	var dialogsInfos []*DialogsInfo
	for _, i := range Dialogs {
		peerClass := i.GetPeer()
		switch peer := peerClass.(type) {
		case *tg.PeerUser:
			for _, user := range allUsers {
				if u, ok := user.(*tg.User); ok {
					if u.ID == peer.UserID {
						info := &DialogsInfo{
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
						info := &DialogsInfo{
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
						info := &DialogsInfo{
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

func getRelationsByPeer(peerID int64) (relations []*entities.ForwardRelation) {
	value, ok := CacheRelationsMap.Load(peerID)
	if !ok { //缓存中不存在，则从数据库查询
		find, _ := fr.Find(peerID)
		if len(find) == 0 {
			AddLog("收到消息，未绑定转发关系")
			empty := make([]*entities.ForwardRelation, 0)
			CacheRelationsMap.Store(peerID, empty)
			return nil
		}
		CacheRelationsMap.Store(peerID, find)
		value = find
	}
	relations = value.([]*entities.ForwardRelation)
	return
}

func saveForwardMsg(fromChatId, toChatId int64, fromMsgId, toMsgId int) {
	go func() {
		_ = fwdMsgDao.Insert(&entities.FwdMsg{
			OriginChatID: fromChatId,
			OriginMsgID:  fromMsgId,
			TargetChatID: toChatId,
			TargetMsgID:  toMsgId,
		})
	}()
}
func getReg(pattern string) (r *regexp.Regexp, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	if pattern == "\\w{8,}" {
		return defaultReg, nil
	}
	r = regexp.MustCompile(pattern)
	return r, nil
}
