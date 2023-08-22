//go:build gui

package serial_tool

import (
	"context"
	"github.com/andlabs/ui"
	"github.com/gookit/slog"
	"time"
)

func GuiSerialDataProvider(ctx context.Context, serRecvChan chan<- []byte, label *ui.Label) {
	first := false
	for {
		var data []byte
		for {
			b := make([]byte, 1024)
			select {
			case <-ctx.Done():
				slog.Warn("I am going:GuiSerialDataProvider")
				return
			default:
				break
			}
			lens, _ := selConn.Read(b)
			// 即使超时也会返回nil
			//if err != nil {
			if lens == 0 {
				time.Sleep(time.Millisecond * 1)
				continue
			}
			data = append(data, b[:lens]...)
			break
			//}
		}
		if !first {
			first = !first
		}
		serRecvChan <- data
		if label.Visible() {
			label.Hide()
		} else {
			label.Show()
		}
	}
}

func GuiSerialDataWriter(ctx context.Context, btChan <-chan []byte, repErr chan<- error, label *ui.Label) {
	for {
		select {
		case <-ctx.Done():
			slog.Warn("ok-GuiSerialDataWriter")
			return
		case res := <-btChan:
			_, err := selConn.Write(res)
			if err != nil {
				repErr <- err
				return
			}
			if label.Visible() {
				label.Hide()
			} else {
				label.Show()
			}
		}
		time.Sleep(time.Millisecond)
	}
}
