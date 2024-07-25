package test

import (
	"fmt"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/utils"
	"testing"
)

func Test_config(t *testing.T) {
	err := dao.Config{}.Set("pp", "哈哈哈")
	if err != nil {
		t.Logf("未通过：%v", err)
	}
}

func Test_String(t *testing.T) {
	str := "撒福a点击阿里可视对讲阿里卡受打击"
	fmt.Println(utils.InsertNewLine(str, 20))
}
