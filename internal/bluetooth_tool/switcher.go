package bluetooth_tool

import (
	"context"
	"github.com/gookit/slog"
	"time"
)

// 串口数据写蓝牙
func BTWriter(ctx context.Context, recv <-chan []byte) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Goroutine BTWR exited successfully!")
			return
		default:
			data := <-recv
			//slog.Warn("writ")
			_ = CurrentDevice.SendDataToRW(data)
			time.Sleep(time.Microsecond * 10)
		}
	}
}
