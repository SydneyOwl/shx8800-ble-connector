package bluetooth_tool

import "tinygo.org/x/bluetooth"

type DeviceFilter func(dev bluetooth.ScanResult) bool

func SHX8800Filter(dev bluetooth.ScanResult) bool {
	return dev.LocalName() == BTNAME_SHX8800
}

func EmptyFilter(dev bluetooth.ScanResult) bool {
	return true
}
