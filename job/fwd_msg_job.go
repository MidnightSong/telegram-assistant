package job

import (
	"context"
	"github.com/midnightsong/telegram-assistant/dao"
	"github.com/midnightsong/telegram-assistant/utils"
	"time"

	"gorm.io/gorm"
)

func init() {
	go func() {
		for {
			fwds, e := dao.FwdMsg{}.All()
			if e != nil && e != gorm.ErrRecordNotFound {
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
	}()
}
