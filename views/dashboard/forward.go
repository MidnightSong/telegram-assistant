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
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/entities"
	"github.com/midnightsong/telegram-assistant/gotgproto/storage"
	"github.com/midnightsong/telegram-assistant/utils"
	"github.com/midnightsong/telegram-assistant/views/icon"
	"time"
)

var fr = dao.ForwardRelation{}

func getForwardView(window fyne.Window) *container.TabItem {
	var frsCache []*entities.ForwardRelation
	var unBindChatTitle = make([]string, 0)
	var dialogsMapById = map[int64]*assistant.DialogsInfo{}     //方便通过peerId找到对应的会话信息,每次点击左侧时更新
	var dialogsMapByTitle = map[string]*assistant.DialogsInfo{} //方便通过title找到对应的会话信息,每次点击左侧时更新
	var listIndex = -1                                          //选中会话来源view某行的索引
	var ac *widget.Accordion                                    //绑定会话关系的树形view（右侧）
	var selectUnBinding *widget.Select                          //选中会话尚未绑定会话的select（右侧）
	var rightBox *fyne.Container                                //右侧view的整体布局
	var addBindButton *widget.Button                            //添加绑定按钮（右侧）
	var selectIndex string                                      //选中的绑定对象的索引（右侧）
	var clickOrigin func(id int)

	topTitle := widget.NewRichTextFromMarkdown("## **消息转发**")
	topTitle.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignCenter //标题居中
	originText := widget.NewRichTextFromMarkdown("## 源会话")
	//originText.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignLeading
	targetText := widget.NewRichTextFromMarkdown("## 目标会话")
	//targetText.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignTrailing
	tmp := container.NewHBox(originText, layout.NewSpacer(), targetText)
	topBox := container.NewVBox(topTitle, tmp)

	listLen := func() int { return len(openedDialogs) }
	createItem := func() fyne.CanvasObject {
		title := widget.NewLabel("尚无已打开的会话")
		//title.Alignment = fyne.TextAlignLeading
		titleType := widget.NewLabel("未知会话类型")
		//titleType.Alignment = fyne.TextAlignTrailing
		return container.NewHBox(title, layout.NewSpacer(), titleType)
	}
	updateItem := func(i widget.ListItemID, o fyne.CanvasObject) {
		item := openedDialogs[i]
		itemType := ""
		if item.Bot {
			itemType = "机器人"
		} else if item.EntityType == storage.TypeUser {
			itemType = "用户"
		} else {
			itemType = "群组|频道"
		}
		o.(*fyne.Container).Objects[0].(*widget.Label).SetText(openedDialogs[i].Title)
		o.(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("【%s】", itemType))
	}
	originList := widget.NewList(listLen, createItem, updateItem) //会话来源的view（左侧）
	clickOrigin = func(id widget.ListItemID) {
		dialogsMapById = map[int64]*assistant.DialogsInfo{}     //方便通过peerId找到对应的会话信息,每次点击左侧时更新
		dialogsMapByTitle = map[string]*assistant.DialogsInfo{} //方便通过会话名称找到对应的会话信息,每次点击左侧时更新
		utils.Select(openedDialogs, func(p *assistant.DialogsInfo, index int) error {
			dialogsMapById[p.PeerId] = p
			dialogsMapByTitle[p.Title] = p
			return nil
		})
		//更新索引
		listIndex = id
		//显示当前选中会话所有绑定的会话
		info := openedDialogs[id]          //当前选中的会话
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
					fmt.Printf("点了一下：%s,下面的Id %d \n", info.Title, itemCopy.ID)
					_ = fr.Add(itemCopy)
				}
			})
			onlyBot.SetChecked(itemCopy.OnlyBot)
			showOrigin := widget.NewCheck("显示消息来源", func(b bool) {
				if b != itemCopy.ShowOrigin {
					itemCopy.ShowOrigin = b
					fmt.Printf("点了一下：%s,下面的 %s 显示消息来源开关\n", info.Title, dialogsMapById[itemCopy.ToPeerID].Title)
					_ = fr.Add(itemCopy)
				}
			})
			showOrigin.SetChecked(itemCopy.ShowOrigin)
			relatedReply := widget.NewCheck("关联转发回复(当显示来源)", func(b bool) {
				if b != itemCopy.RelatedReply {
					itemCopy.RelatedReply = b
					fmt.Printf("点了一下：%s,下面的 %s 关联转发回复开关\n", info.Title, dialogsMapById[itemCopy.ToPeerID].Title)
					_ = fr.Add(itemCopy)
				}
			})
			relatedReply.SetChecked(itemCopy.RelatedReply)
			conditionLabel := widget.NewRichTextFromMarkdown("## 触发条件")
			conditionLabel.Segments[0].(*widget.TextSegment).Style.Alignment = fyne.TextAlignTrailing //标题居右
			regexLabel := widget.NewLabel("文字(正则表达式)")
			regexEntry := widget.NewEntry()
			regexEntry.OnChanged = func(s string) {
				fmt.Printf("正则：%s\n", s)
			}
			regex := container.New(layout.NewFormLayout(), regexLabel, regexEntry)
			regexEntry.Text = itemCopy.Regex
			regexEntry.PlaceHolder = "\\w{8,}"
			mustMedia := widget.NewCheck("转发消息中必须带图片", func(b bool) {
				if b != itemCopy.MustMedia {
					itemCopy.MustMedia = b
					fmt.Printf("点了一下：%s,下面的 %s 转发消息中必须带图片开关\n", info.Title, dialogsMapById[itemCopy.ToPeerID].Title)
					_ = fr.Add(itemCopy)
				}
			})
			mustMedia.SetChecked(itemCopy.MustMedia)
			deleteBindButton := widget.NewButton("删除", func() {
				dialog.ShowConfirm("确定？", "将删除该绑定对象", func(b bool) {
					if b {
						fmt.Printf("删除绑定对象：%d\n", itemCopy.ID)
						fr.DeleteById(itemCopy.ID)
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
		origin := openedDialogs[listIndex]       //源会话
		target := dialogsMapByTitle[selectIndex] //绑定目标
		bind := &entities.ForwardRelation{
			PeerID:       origin.PeerId,
			ToPeerID:     target.PeerId,
			OnlyBot:      true,
			ShowOrigin:   true,
			RelatedReply: true,
			Regex:        "\\w{8,}",
			MustMedia:    true,
		}
		_ = fr.Add(bind)
		selectUnBinding.ClearSelected()
		time.Sleep(time.Millisecond * 15) //不休眠的话，两个组件刷新后可能重叠在一起
		clickOrigin(listIndex)            //刷新右侧view
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
