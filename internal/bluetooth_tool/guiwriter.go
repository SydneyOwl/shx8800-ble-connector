//go:build gui

package bluetooth_tool

import (
	"bytes"
	"context"
	"github.com/andlabs/ui"
	"github.com/gookit/slog"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/exceptions"
	"time"
)

// 串口数据写蓝牙
func GuiBTWriter(ctx context.Context, recv <-chan []byte, repErr chan<- error, label *ui.Label) {
	for {
		select {
		case <-ctx.Done():
			slog.Warn("ok-GuiBTWriter")
			return
		case data := <-recv:
			if label.Visible() {
				label.Hide()
			} else {
				label.Show()
			}
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
		}
		time.Sleep(time.Millisecond)
	}
}
