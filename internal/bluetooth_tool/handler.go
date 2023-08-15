package bluetooth_tool

import (
	"github.com/gookit/slog"
	"os"
	"tinygo.org/x/bluetooth"
)

type RecvHandler func(recvChan chan<- []byte) func(c []byte)

func CheckRecvHandler(c []byte) {
	slog.Debug(c)
}

func RWRecvHandler(recvChan chan<- []byte) func(c []byte) {
	return func(c []byte) {
		//slog.Warn("repl")
		recvChan <- c
	}
}

func DisconnectHandler(_ bluetooth.Address, connected bool) {
	if !connected && !call_disconnect {
		// Call cleanup!
		slog.Warnf("蓝牙断开！程序退出")
		os.Exit(-1)
	}
}
