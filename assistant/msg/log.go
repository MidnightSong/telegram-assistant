package msg

import (
	"fmt"
	"time"
)

var Log = make(chan string, 10)
var LogSwitch bool

func AddLog(msg string) {
	go func() {
		if LogSwitch {
			now := time.Now().Format(time.DateTime)
			Log <- fmt.Sprintf("%s \n\n %s \n\n", now, msg)
		}
	}()
}
