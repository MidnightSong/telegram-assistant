package main

import (
	"fmt"
	"go/format"
	"html/template"
	"os"
	"strings"

	"github.com/midnightsong/telegram-assistant/gotgproto/generator/parser"
)

var helperFuncsCUTempl = template.Must(template.New("cuHelpers").Parse(helperFuncsCU))

var hardCodedReplacements = map[string]string{
	"EditAdminOpts": "ext.EditAdminOpts",
}

func readContextFile() []byte {
	b, err := os.ReadFile("ext/context.go")
	if err != nil {
		panic("failed to read context file: " + err.Error())
	}
	return b
}

func generateCUHelpers() {
	fmt.Println("Reading context.go")
	ctxFile := readContextFile()
	builder := strings.Builder{}
	builder.WriteString(predefinedCU)
	fmt.Println("Parsing all context methods...")
	for _, method := range parser.ParseMethods(string(ctxFile)) {
		if strings.ToLower(string(method.Name[0])) == string(method.Name[0]) {
			continue
		}
		params := method.Params
		if !(strings.Contains(params, "chatId ") ||
			strings.Contains(params, "chatId, ") ||
			strings.Contains(params, "userId ") ||
			strings.Contains(params, "userId, ")) {
			continue
		}
		params = strings.ReplaceAll(params, "chatId, ", "")
		params = strings.ReplaceAll(params, "chatId int64, ", "")
		params = strings.ReplaceAll(params, "chatId int64", "")
		params = strings.ReplaceAll(params, "userId, ", "")
		params = strings.ReplaceAll(params, "userId int64, ", "")
		params = strings.ReplaceAll(params, "userId int64", "")
		for repl, valrepl := range hardCodedReplacements {
			params = strings.ReplaceAll(params, repl, valrepl)
		}
		chatFrame, userFrame := getFrames(method)
		inputIdParams, fetchedIdParams := getIdParams(chatFrame, userFrame)
		fmt.Printf("Executing generic helper function template for Context.%s\n", method.Name)
		err := helperFuncsCUTempl.Execute(&builder, contextHelpers{
			FuncName:        method.Name,
			FuncParams:      params,
			FuncReturn:      method.Return,
			FilledParams:    filledParams(params),
			ChatFrame:       chatFrame,
			UserFrame:       userFrame,
			InputIdParams:   inputIdParams,
			FetchedIdParams: fetchedIdParams,
			// DefaultValues: goodErrReturns(method.Return),
		})
		if err != nil {
			fmt.Printf("failed to generate helper for ext.Context.%s because %s\n", method.Name, err.Error())
			continue
		}
	}
	fmt.Println("Writing gen_cu.go for context generic helpers...")
	_ = writeFile(builder, "generic/gen_cu.go")
}

func getFrames(method *parser.Method) (chatFrame, userFrame string) {
	returnVals := goodErrReturns(method.Return)
	for _, param := range getParamNamesArray(method.Params) {
		switch param {
		case "chatId":
			chatFrame = fmt.Sprintf(FrameProperty, "chat", "chat", returnVals)
		case "userId":
			userFrame = fmt.Sprintf(FrameProperty, "user", "user", returnVals)
		}
	}
	return
}

func getIdParams(chatFrame, userFrame string) (inputIdParam, fetchedIdParam string) {
	switch {
	case chatFrame != "" && userFrame != "":
		inputIdParam = "chat, user chatUnion"
		fetchedIdParam = "chatId, userId"
	case chatFrame != "":
		inputIdParam = "chat chatUnion"
		fetchedIdParam = "chatId"
	case userFrame != "":
		inputIdParam = "user chatUnion"
		fetchedIdParam = "userId"
	}
	return
}

const predefinedCU = `
// GoTGProto Generic Helpers 
// WARNING: This file is autogenerated, please DO NOT EDIT

package generic

import (
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/types"
	"github.com/gotd/td/tg"
)

type ChatUnion interface {
	int | int64 | string
}

func getIdByUnion[chatUnion ChatUnion](ctx *ext.Context, chat chatUnion) (int64, error) {
	switch val := any(chat).(type) {
	case string:
		username := val
		peer := ctx.PeerStorage.GetPeerByUsername(username)
		if peer.ID != 0 {
			return peer.ID, nil
		}
		chat, err := ctx.ResolveUsername(username)
		return chat.GetID(), err
	case int64:
		return val, nil
	case int:
		return int64(val), nil
	}
	// Unreachable
	return 0, nil
}
`

type contextHelpers struct {
	FuncName     string
	FuncParams   string
	FuncReturn   string
	FilledParams string
	// DefaultValues   string
	ChatFrame       string
	UserFrame       string
	InputIdParams   string
	FetchedIdParams string
}

// chatId, err := getIdByUnion(ctx, chat)
//
//	if err != nil {
//		return {{.DefaultValues}}
//	}
const FrameProperty = `
%sId, err := getIdByUnion(ctx, %s)
	if err != nil {return %s}
`

const helperFuncsCU = `
// {{.FuncName}} is a generic helper for ext.Context.{{.FuncName}} method.
func {{.FuncName}}[chatUnion ChatUnion] (ctx *ext.Context, {{.InputIdParams}}, {{.FuncParams}}) ({{.FuncReturn}}) {
	{{.ChatFrame}}
	{{.UserFrame}}
	return ctx.{{.FuncName}}({{.FetchedIdParams}}, {{.FilledParams}})
}`

func writeFile(builder strings.Builder, filename string) error {
	write, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}

	bs := []byte(builder.String())

	_, err = write.WriteAt(bs, 0)
	if err != nil {
		return fmt.Errorf("failed to write unformatted file %s: %w", filename, err)
	}

	fmted, err := format.Source(bs)
	if err != nil {
		return fmt.Errorf("failed to format file %s: %w", filename, err)
	}

	err = write.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to truncate file %s: %w", filename, err)
	}

	_, err = write.WriteAt(fmted, 0)
	if err != nil {
		return fmt.Errorf("failed to write final file %s: %w", filename, err)
	}

	return nil
}

func getParamNamesArray(paramStr string) []string {
	paramArr := make([]string, 0)
	params := strings.Split(paramStr, ", ")
	for _, param := range params {
		paramFields := strings.Fields(param)
		if len(paramFields) == 0 {
			continue
		}
		paramArr = append(paramArr, paramFields[0])
	}
	return paramArr
}

func filledParams(paramStr string) string {
	return strings.Join(getParamNamesArray(paramStr), ", ")
}

func goodErrReturns(s string) string {
	returns := strings.Split(s, ", ")
	if len(returns) == 0 {
		return `
	_ = err
	return`
	}
	goodReturns := make([]string, 0)
	for _, returnStr := range returns {
		switch {
		case strings.HasPrefix(returnStr, "*"), strings.HasSuffix(returnStr, "Class"), strings.HasPrefix(returnStr, "[]"):
			goodReturns = append(goodReturns, "nil")
		case returnStr == "error":
			goodReturns = append(goodReturns, "err")
		case returnStr == "bool":
			goodReturns = append(goodReturns, "false")
		case returnStr == "string":
			goodReturns = append(goodReturns, "\"\"")
		case returnStr == "int64", returnStr == "int":
			goodReturns = append(goodReturns, "0")
		default:
			goodReturns = append(goodReturns, returnStr+"{}")
		}
	}
	return strings.Join(goodReturns, ", ")
}
