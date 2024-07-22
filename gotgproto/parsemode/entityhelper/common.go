package entityhelper

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gotd/td/tg"
)

// EntityRoot is used to create message entities from the input string through its various methods.
type EntityRoot struct {
	String   string
	Entities tg.MessageEntityClassArray
}

type entity rune

const (
	BoldEntity      entity = 'b'
	ItalicEntity    entity = 'i'
	UnderlineEntity entity = 'u'
	CodeEntity      entity = 'c'
	StrikeEntity    entity = '~'
	SpoilertEntity  entity = 's'
)

// Combine creates a new entity root and appends the provided string as combined entities to this entity root.
func Combine(s string, entity1, entity2 entity) *EntityRoot {
	return startParsing().Combine(s, entity1, entity2)
}

// Combine function combines the entity1 and entity2 and appends the resultant entity to the EntityRoot.
func (root *EntityRoot) Combine(s string, entity1, entity2 entity) *EntityRoot {
	root.setNormalEntity(s, entity1)
	root.setNormalEntity(s, entity2)
	root.String += s
	return root
}

// CombineToLink creates a new entity root and appends the provided string as combined link entities to this entity root.
func CombineToLink(text, link string, entity entity) *EntityRoot {
	return startParsing().CombineToLink(text, link, entity)
}

// CombineToLink function combines the given entity to the link entity of the EntityRoot.
func (root *EntityRoot) CombineToLink(text, link string, entity entity) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityTextURL{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(text), URL: link})
	root.setNormalEntity(text, entity)
	root.String += text
	return root
}

func (root *EntityRoot) setNormalEntity(s string, e entity) {
	switch e {
	case BoldEntity:
		root.Entities = append(root.Entities, &tg.MessageEntityBold{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	case ItalicEntity:
		root.Entities = append(root.Entities, &tg.MessageEntityItalic{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	case UnderlineEntity:
		root.Entities = append(root.Entities, &tg.MessageEntityUnderline{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	case CodeEntity:
		root.Entities = append(root.Entities, &tg.MessageEntityCode{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	case StrikeEntity:
		root.Entities = append(root.Entities, &tg.MessageEntityStrike{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	case SpoilertEntity:
		root.Entities = append(root.Entities, &tg.MessageEntitySpoiler{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	}
}

// StartParsing function creates an empty EntityRoot.
// DEPRECATED
func StartParsing() *EntityRoot {
	fmt.Println("GoTGProto: func StartParsing() is deprecated, please use individual entity types instead.")
	return startParsing()
}

// startParsing function creates an empty EntityRoot.
// Only for internal use
func startParsing() *EntityRoot {
	return &EntityRoot{
		String:   "",
		Entities: tg.MessageEntityClassArray{},
	}
}

// Bold creates a new entity root and appends the provided string as bold to this entity root.
func Bold(s string) *EntityRoot {
	return startParsing().Bold(s)
}

// Bold appends the provided string as bold to the entity root.
func (root *EntityRoot) Bold(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityBold{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	root.String += s
	return root
}

// Italic creates a new entity root and appends the provided string as italic to this entity root.
func Italic(s string) *EntityRoot {
	return startParsing().Italic(s)
}

// Italic appends the provided string as italic to the entity root.
func (root *EntityRoot) Italic(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityItalic{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	root.String += s
	return root
}

// Underline creates a new entity root and appends the provided string as underline to this entity root.
func Underline(s string) *EntityRoot {
	return startParsing().Underline(s)
}

// Underline appends the provided string as underline to the entity root.
func (root *EntityRoot) Underline(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityUnderline{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	root.String += s
	return root
}

// Code creates a new entity root and appends the provided string as code to this entity root.
func Code(s string) *EntityRoot {
	return startParsing().Code(s)
}

// Code appends the provided string as code/mono to the entity root.
func (root *EntityRoot) Code(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityCode{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	root.String += s
	return root
}

// Strike creates a new entity root and appends the provided string as strike to this entity root.
func Strike(s string) *EntityRoot {
	return startParsing().Strike(s)
}

// Strike appends the provided string as strike to the entity root.
func (root *EntityRoot) Strike(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityStrike{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	root.String += s
	return root
}

// Spoiler creates a new entity root and appends the provided string as spoiler to this entity root.
func Spoiler(s string) *EntityRoot {
	return startParsing().Spoiler(s)
}

// Spoiler appends the provided string as spoiler to the entity root.
func (root *EntityRoot) Spoiler(s string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntitySpoiler{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(s)})
	root.String += s
	return root
}

// Link creates a new entity root and appends the provided string as a link to this entity root.
func Link(text, url string) *EntityRoot {
	return startParsing().Link(text, url)
}

// Link appends the provided link to the entity root.
func (root *EntityRoot) Link(text, url string) *EntityRoot {
	root.Entities = append(root.Entities, &tg.MessageEntityTextURL{
		Offset: utf8.RuneCountInString(root.String),
		Length: utf8.RuneCountInString(text),
		URL:    url,
	})
	root.String += text
	return root
}

// Mention creates a new entity root and appends the provided string as a mention to this entity root.
func Mention(name string, user tg.InputUserClass) *EntityRoot {
	return startParsing().Mention(name, user)
}

// Mention appends the provided mention to the entity root.
func (root *EntityRoot) Mention(name string, user tg.InputUserClass) *EntityRoot {
	root.Entities = append(root.Entities, &tg.InputMessageEntityMentionName{
		Offset: utf8.RuneCountInString(root.String),
		Length: utf8.RuneCountInString(name),
		UserID: user,
	})
	root.String += name
	return root
}

// Plain creates a new entity root and appends the provided string as plain text to this entity root.
func Plain(s string) *EntityRoot {
	return startParsing().Plain(s)
}

// Plain appends the provided text to the entity root as it is.
func (root *EntityRoot) Plain(text string) *EntityRoot {
	// root.Entities = append(root.Entities, &tg.MessageEntityUnknown{Offset: utf8.RuneCountInString(root.String), Length: utf8.RuneCountInString(text)})
	// root.Offset = utf8.RuneCountInString(root.String)
	root.String += text
	return root
}

// BotMention creates a new entity root and appends the provided string as a mention to this entity root.
func BotMention(text string, user interface{}) *EntityRoot {
	return startParsing().BotMention(text, user)
}

// BotMention creates a telegram user mention link with the provided user and text to display.
//
// NOTE: This method works for bots only
func (root *EntityRoot) BotMention(text string, user interface{}) *EntityRoot {
	switch user := user.(type) {
	case int, int64:
		return root.Link(text, fmt.Sprintf("tg://user?id=%d", user))
	case string:
		return root.Link(text, fmt.Sprintf("tg://resolve?domain=%s", strings.TrimPrefix(user, "@")))
	}
	return root
}

func RawEntity(entity tg.MessageEntityClass) *EntityRoot {
	return startParsing().RawEntity(entity)
}

func (root *EntityRoot) RawEntity(entity tg.MessageEntityClass) *EntityRoot {
	root.Entities = append(root.Entities, entity)
	for i := 1; i <= entity.GetLength(); i++ {
		root.String += " "
	}
	return root
}

func (root *EntityRoot) GetEntities() []tg.MessageEntityClass {
	return root.Entities
}

func (root *EntityRoot) GetString() string {
	return root.String
}
