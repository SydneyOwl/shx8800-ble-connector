package bluetooth_tool

import (
	"bytes"
	"context"
	"github.com/gookit/slog"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/exceptions"
	"time"
)

// 串口数据写蓝牙
func BTWriter(ctx context.Context, recv <-chan []byte, repErr chan<- error) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Goroutine BTWR exited successfully!")
			return
		default:
			data := <-recv
			//slog.Warn("writ")
			err := CurrentDevice.SendDataToRW(data)
			if err != nil {
				repErr <- err
				return
			}
			//遇到最后一个帧对讲机将重启
			if bytes.Equal(FINAL_DATA_STARTER, data[0:3]) {
				repErr <- exceptions.TransferDone
				return
			}
			time.Sleep(time.Microsecond * 5)
		}
	}
}
