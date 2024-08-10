package dashboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/midnightsong/telegram-assistant/assistant/msg"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/gotgproto/storage"
	"github.com/midnightsong/telegram-assistant/utils"
	"github.com/midnightsong/telegram-assistant/views/icon"
	"reflect"
	"time"
	"unsafe"
)

var fr = dao.ForwardRelation{}

func getForwardView(window fyne.Window) *container.TabItem {
	var frsCache []*entities.ForwardRelation
	var unBindChatTitle = make([]string, 0)
	var dialogsMapById = map[int64]*msg.DialogsInfo{}     //方便通过peerId找到对应的会话信息,每次点击左侧时更新
	var dialogsMapByTitle = map[string]*msg.DialogsInfo{} //方便通过title找到对应的会话信息,每次点击左侧时更新
	var listIndex = -1                                    //选中会话来源view某行的索引
	var ac *widget.Accordion                              //绑定会话关系的树形view（右侧）
	var selectUnBinding *widget.Select                    //选中会话尚未绑定会话的select（右侧）
	var rightBox *fyne.Container                          //右侧view的整体布局
	var addBindButton *widget.Button                      //添加绑定按钮（右侧）
	var selectIndex string                                //选中的绑定对象的索引（右侧）
	var clickOrigin func(id int)

	topTitle := widget.NewRichTextFromMarkdown("## **消息转发**")
	topTitle.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignCenter //标题居中
	originText := widget.NewRichTextFromMarkdown("## 源会话")
	targetText := widget.NewRichTextFromMarkdown("## 目标会话")
	tmp := container.NewHBox(originText, layout.NewSpacer(), targetText)
	topBox := container.NewVBox(topTitle, tmp)

	listLen := func() int { return len(msg.OpenedDialogs) }
	createItem := func() fyne.CanvasObject {
		title := widget.NewLabel("尚无已打开的会话")
		title.Truncation = fyne.TextTruncateClip
		titleType := widget.NewIcon(theme.AccountIcon())
		return container.NewHBox(titleType, title)
	}
	updateItem := func(i widget.ListItemID, o fyne.CanvasObject) {
		item := msg.OpenedDialogs[i]
		var itemType *widget.Icon
		if item.Bot {
			itemType = widget.NewIcon(icon.GetIcon(icon.Bot))
		} else if item.EntityType == storage.TypeUser {
			itemType = widget.NewIcon(icon.GetIcon(icon.People))
		} else {
			itemType = widget.NewIcon(icon.GetIcon(icon.Group))
		}
		o.(*fyne.Container).Objects[1].(*widget.Label).SetText(msg.OpenedDialogs[i].Title)
		o.(*fyne.Container).Objects[0] = itemType
	}
	originList := widget.NewList(listLen, createItem, updateItem) //会话来源的view（左侧）
	clickOrigin = func(id widget.ListItemID) {
		dialogsMapById = map[int64]*msg.DialogsInfo{}     //方便通过peerId找到对应的会话信息,每次点击左侧时更新
		dialogsMapByTitle = map[string]*msg.DialogsInfo{} //方便通过会话名称找到对应的会话信息,每次点击左侧时更新
		utils.Select(msg.OpenedDialogs, func(p *msg.DialogsInfo, index int) error {
			dialogsMapById[p.PeerId] = p
			dialogsMapByTitle[p.Title] = p
			return nil
		})
		//更新索引
		listIndex = id
		//显示当前选中会话所有绑定的会话
		info := msg.OpenedDialogs[id]      //当前选中的会话
		frsCache, _ = fr.Find(info.PeerId) //已存库的绑定关系表
		//先清空绑定会话目标view
		ac.Items = make([]*widget.AccordionItem, 0)
		time.Sleep(time.Millisecond * 15)
		//展示绑定会话
		for _, item := range frsCache {
			to := dialogsMapById[item.ToPeerID]
			if to == nil {
				//如果数据库中有绑定记录，但打开的会话中没有，说明已关闭会话，则删除对应绑定关系
				fr.DeleteById(item.ID)
				continue
			}
			// 创建局部变量，避免闭包捕获问题
			itemCopy := item
			//叶子节点内容
			onlyBot := widget.NewCheck("仅转发机器人的消息", func(b bool) {
				if b != itemCopy.OnlyBot {
					itemCopy.OnlyBot = b
					_ = fr.Add(itemCopy)
					msg.CacheRelationsMap.Delete(itemCopy.PeerID) //删除处理接收消息处理模块的缓存
				}
			})
			onlyBot.SetChecked(itemCopy.OnlyBot)
			showOrigin := widget.NewCheck("显示消息来源", func(b bool) {
				if b != itemCopy.ShowOrigin {
					itemCopy.ShowOrigin = b
					_ = fr.Add(itemCopy)
					msg.CacheRelationsMap.Delete(itemCopy.PeerID) //删除处理接收消息处理模块的缓存
				}
			})
			showOrigin.SetChecked(itemCopy.ShowOrigin)
			relatedReply := widget.NewCheck("关联回复原始消息", func(b bool) {
				if b != itemCopy.RelatedReply {
					itemCopy.RelatedReply = b
					_ = fr.Add(itemCopy)
					msg.CacheRelationsMap.Delete(itemCopy.PeerID) //删除处理接收消息处理模块的缓存
				}
			})
			relatedReply.SetChecked(itemCopy.RelatedReply)
			conditionLabel := widget.NewRichTextFromMarkdown("## 触发条件")
			conditionLabel.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignTrailing //标题居右
			regexLabel := widget.NewLabel("文字(正则表达式)")
			regexEntry := widget.NewEntry()
			v := reflect.ValueOf(regexEntry).Elem().FieldByName("onFocusChanged")
			va := func(b bool) {
				if regexEntry.Text != itemCopy.Regex {
					itemCopy.Regex = regexEntry.Text
					_ = fr.Add(itemCopy)
					msg.CacheRelationsMap.Delete(itemCopy.PeerID) //删除处理接收消息处理模块的缓存
				}
			}
			reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(va))

			regex := container.New(layout.NewFormLayout(), regexLabel, regexEntry)
			regexEntry.Text = itemCopy.Regex
			regexEntry.PlaceHolder = "\\w{8,}"
			mustMedia := widget.NewCheck("转发消息中必须带图片", func(b bool) {
				if b != itemCopy.MustMedia {
					itemCopy.MustMedia = b
					_ = fr.Add(itemCopy)
					msg.CacheRelationsMap.Delete(itemCopy.PeerID) //删除处理接收消息处理模块的缓存
				}
			})
			mustMedia.SetChecked(itemCopy.MustMedia)
			deleteBindButton := widget.NewButton("删除", func() {
				dialog.ShowConfirm("确定？", "将删除该绑定对象", func(b bool) {
					if b {
						fr.DeleteById(itemCopy.ID)
						msg.CacheRelationsMap.Delete(itemCopy.PeerID) //删除处理接收消息处理模块的缓存
						clickOrigin(listIndex)
					}
				}, window)
			})
			deleteBindButton.Importance = widget.WarningImportance
			deleteBox := container.NewHBox(layout.NewSpacer(), deleteBindButton)
			inTreeBox := container.NewAdaptiveGrid(2,
				onlyBot, showOrigin, relatedReply, layout.NewSpacer(),
				conditionLabel, layout.NewSpacer(), regex, mustMedia, deleteBox)

			//追加到绑定会话列表下
			ac.Append(widget.NewAccordionItem(to.Title, inTreeBox))
		}
		//缓存当前会话尚未绑定的会话
		unBindChatTitle = make([]string, 0)
		for i, v := range dialogsMapById {
			flag := true
			for _, item := range frsCache {
				if item.ToPeerID == i {
					flag = false
					break
				}

			}
			if flag {
				unBindChatTitle = append(unBindChatTitle, v.Title)
			}
		}
		selectUnBinding.Options = unBindChatTitle
		ac.CloseAll()
		//ac.Items[0].Detail.(*fyne.Container).Objects[0].(*widget.Check).Checked = true
		rightBox.Show()
	}
	originList.OnSelected = clickOrigin
	ac = widget.NewAccordion()
	//在绑定会话列表的最后一行增加添加绑定关系的按钮
	addBindButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		if listIndex == -1 || selectIndex == "" { //默认值，尚未更新
			return
		}
		origin := msg.OpenedDialogs[listIndex]   //源会话
		target := dialogsMapByTitle[selectIndex] //绑定目标
		bind := &entities.ForwardRelation{
			PeerID:       origin.PeerId,
			ToPeerID:     target.PeerId,
			OnlyBot:      true,
			ShowOrigin:   true,
			RelatedReply: true,
			Regex:        "\\w{8,}",
			MustMedia:    true,
			PeerTitle:    origin.Title,
			ToPeerTitle:  target.Title,
		}
		_ = fr.Add(bind)
		selectUnBinding.ClearSelected()
		time.Sleep(time.Millisecond * 15)         //不休眠的话，两个组件刷新后可能重叠在一起
		clickOrigin(listIndex)                    //刷新右侧view
		msg.CacheRelationsMap.Delete(bind.PeerID) //删除处理接收消息处理模块的缓存
	})
	selectUnBinding = widget.NewSelect(unBindChatTitle, func(s string) {
		selectIndex = s
	})
	selectUnBinding.PlaceHolder = "--- 添加绑定 ---"
	bottomBox := container.NewBorder(nil, nil, nil, addBindButton, selectUnBinding)
	rightBox = container.NewVBox(ac, bottomBox)
	rightBox.Hide()
	splitBox := container.NewHSplit(originList, rightBox)
	splitBox.Offset = 0.3
	border := container.NewBorder(topBox, nil, nil, nil, splitBox)
	return container.NewTabItemWithIcon("", icon.GetIcon(icon.Forward), border)
}
