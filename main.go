package main

import (
	"github.com/midnightsong/telegram-assistant/utils"
	"github.com/midnightsong/telegram-assistant/views"
	"os"
)

var m *bool //迁移数据库
func init() {
	name, _ := utils.FileName("/view/font/YaHei.ttf")
	fontPath := name
	os.Setenv("FYNE_FONT", fontPath)
	/*m = flag.Bool("m", false, "migrate tables")
	_ = flag.Bool("o", false, "debug")*/
}

func main() {
	//dao.Migrate()

	views.Run()
}
