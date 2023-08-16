package serial_tool

import (
	"context"
	"github.com/gookit/slog"
	"go.bug.st/serial"
	"sort"
	"time"
)

var selConn serial.Port

func ScanPort() ([]string, error) {
	ports, err := serial.GetPortsList()
	sort.Strings(ports)
	return ports, err
}
func ConnPort(portName string) error {
	mode := &serial.Mode{
		BaudRate: BAUDRATE,
	}
	sel, err := serial.Open(portName, mode)
	if err != nil {
		return err
	}
	selConn = sel

	_ = selConn.SetReadTimeout(time.Millisecond * SERIAL_READ_TIMEOUT_MS)
	return nil
}

func SerialDataProvider(ctx context.Context, serRecvChan chan<- []byte) {
	first := false
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Goroutine SDPR exited successfully!")
			return
		default:
			var data []byte
			for {
				b := make([]byte, 1024)
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
				slog.Info("软件已连接！")
			}
			serRecvChan <- data
		}
	}
}

func SerialDataWriter(ctx context.Context, btChan <-chan []byte, repErr chan<- error) {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Goroutine SDWE exited successfully!")
			return
		default:
			res := <-btChan
			_, err := selConn.Write(res)
			if err != nil {
				repErr <- err
				return
			}
			time.Sleep(time.Millisecond * 1)
		}
	}
}

func ShutPort() {
	_ = selConn.Close()
}
