package bluetooth_tool

import (
	"context"
	"github.com/gookit/slog"
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
			time.Sleep(time.Microsecond * 10)
		}
	}
}
