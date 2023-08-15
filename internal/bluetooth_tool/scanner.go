package bluetooth_tool

import (
	"github.com/gookit/slog"
	"time"
	"tinygo.org/x/bluetooth"
)

var adapter = bluetooth.DefaultAdapter

func GetAvailableBtDevList(filter DeviceFilter) ([]bluetooth.ScanResult, error) {
	devList := make([]bluetooth.ScanResult, 0)
	errChan := make(chan error, 0)
	removeRepeat := make(map[string]bool)
	if err := adapter.Enable(); err != nil {
		return nil, err
	}

	go func(errChan chan error, devList *[]bluetooth.ScanResult) {
		err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			slog.Tracef("找到设备: %s, %d, %s", device.Address.String(), device.RSSI, device.LocalName())
			if filter(device) {
				if !removeRepeat[device.Address.String()] {
					slog.Debugf("设备: %s, %d, %s", device.Address.String(), device.RSSI, device.LocalName())
					removeRepeat[device.Address.String()] = true
					*devList = append(*devList, device)
				}
			}
		})
		if err != nil {
			errChan <- err
		}
	}(errChan, &devList)

	select {
	case err := <-errChan:
		return []bluetooth.ScanResult{}, err
	case <-time.After(BT_SCAN_TIMEOUT * time.Second):
		err := adapter.StopScan()
		return devList, err
	}
}
func ConnectByMac(mac bluetooth.Address) (*bluetooth.Device, error) {
	return adapter.Connect(mac, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(BT_CONNECT_TIMEOUT * time.Second),
	})
}
func SetHandler(connHandler func(address bluetooth.Address, connected bool)) {
	adapter.SetConnectHandler(connHandler)
}
func DisconnectDevice(shx *bluetooth.Device) {
	call_disconnect = true
	_ = shx.Disconnect()
}
