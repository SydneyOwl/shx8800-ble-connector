package bluetooth_tool

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/sydneyowl/shx8800-ble-connector/internal/gui_tool"
	"os"
	"tinygo.org/x/bluetooth"
)

type RecvHandler func(recvChan chan<- []byte, btOut chan<- struct{}) func(c []byte)

func CheckRecvHandler(c []byte) {
	slog.Debug(c)
}

func RWRecvHandler(recvChan chan<- []byte, btOut chan<- struct{}) func(c []byte) {
	return func(c []byte) {
		//slog.Warn("repl")
		if btOut != nil {
			if gui_tool.CheckedLog() {
				go gui_tool.AddLog(fmt.Sprintf("Bluetooth Recv: {%x}", c))
			}
			btOut <- struct{}{}
		}
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
