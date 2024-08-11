package job

import (
	"context"
	"errors"
	"github.com/midnightsong/telegram-assistant/assistant"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/utils"
	"gorm.io/gorm"
	"os"
	"os/exec"
	"time"
)

// Init 直接调用init函数会导致fyne配置读取出现问题
func Init() {
	go cleanFwd()
	go checkAuth()
}
func checkAuth() {
	for {
		utils.LogInfo(context.Background(), "验证激活状态")
		time.Sleep(time.Hour)
		var err error
		var result *assistant.AuthResponse
		for i := 0; i < 3; i++ {
			result, err = assistant.Auth()
			if err != nil {
				time.Sleep(time.Minute * 10)
				continue
			} else {
				break
			}
		}
		if err != nil {
			restart()
		}
		if result.Code == 2000 {
			continue
		} else {
			restart()
		}
	}
}

// cleanFwd  删除超过一天的历史转发的消息id
func cleanFwd() {
	for {
		fwds, e := dao.FwdMsg{}.All()
		if e != nil && !errors.Is(e, gorm.ErrRecordNotFound) {
			utils.LogInfo(context.Background(), "查询转发的消息记录失败："+e.Error())
			continue
		}
		now := time.Now().Unix()
		for _, fwd := range fwds {
			if now-fwd.FwdTime >= 60*60*24 { //单位秒
				e = dao.FwdMsg{}.Delete(fwd)
				if e != nil {
					utils.LogError(context.Background(), "删除超过一天时间的转发记录失败："+e.Error())
					continue
				}
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

// restart  重启自身
func restart() {
	executable, err := os.Executable()
	if err != nil {
		os.Exit(0)
	}
	// 获取当前的命令行参数
	args := os.Args
	// 重新启动程序
	cmd := exec.Command(executable, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		os.Exit(0)
	}
	os.Exit(0)
}
