//go:build gui

package bluetooth_tool

import (
	"github.com/gookit/slog"
	"time"
	"tinygo.org/x/bluetooth"
)

func GetAvailableBtDevListViaChannel(filter DeviceFilter, devList *[]bluetooth.ScanResult, errChan chan error) {
	removeRepeat := make(map[string]bool)
	if err := adapter.Enable(); err != nil {
		errChan <- err
		close(errChan)
		return
	}
	go func(errChan chan error, devList *[]bluetooth.ScanResult) {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			if filter(device) {
				if !removeRepeat[device.Address.String()] {
					removeRepeat[device.Address.String()] = true
					*devList = append(*devList, device)
				}
			}
		})
		if err != nil {
			errChan <- err
		}
	}(errChan, devList)

	select {
	case <-time.After(BT_SCAN_TIMEOUT * time.Second):
		err := adapter.StopScan()
		if err != nil {
			errChan <- err
		}
		close(errChan)
		return
	}
}
func ConnectByMacNoBlock(mac bluetooth.Address, connChan chan *bluetooth.Device) {
	dev, err := adapter.Connect(mac, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(BT_CONNECT_TIMEOUT * time.Second),
	})
	if err != nil {
		slog.Error(err)
		connChan <- nil
		return
	}
	connChan <- dev
}
