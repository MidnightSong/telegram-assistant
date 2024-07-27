package msg

import (
	"fmt"
	"time"
)

var Log = make(chan string, 10)

func AddLog(msg string) {
	go func() {
		now := time.Now().Format(time.DateTime)
		Log <- fmt.Sprintf("%s\n%s", now, msg)
	}()
}
