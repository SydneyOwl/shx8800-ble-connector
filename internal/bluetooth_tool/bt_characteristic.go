package bluetooth_tool

import "tinygo.org/x/bluetooth"

type BTCharacteristic struct {
	CheckCharacteristic *bluetooth.DeviceCharacteristic
	RWCharacteristic    *bluetooth.DeviceCharacteristic

	CheckRecvHandler RecvHandler
	RWRecvHandler    RecvHandler
}

func (btc *BTCharacteristic) SetCheckReceiveHandler(handler RecvHandler, recvChan chan<- []byte) {
	btc.CheckRecvHandler = handler
	_ = btc.CheckCharacteristic.EnableNotifications(handler(recvChan))
}

func (btc *BTCharacteristic) SetReadWriteReceiveHandler(handler RecvHandler, recvChan chan<- []byte) {
	btc.RWRecvHandler = handler
	_ = btc.RWCharacteristic.EnableNotifications(handler(recvChan))
}
func (btc *BTCharacteristic) SendDataToCheck(data []byte) error {
	_, err := btc.CheckCharacteristic.Write(data)
	return err
}
func (btc *BTCharacteristic) SendDataToRW(data []byte) error {
	BtMutex.Lock()
	_, err := btc.RWCharacteristic.Write(data)
	BtMutex.Unlock()
	return err
}
