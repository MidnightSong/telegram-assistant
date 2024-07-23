package msg

var Log = make(chan string, 10)

func AddLog(msg string) {
	go func() {
		Log <- msg
	}()
}
