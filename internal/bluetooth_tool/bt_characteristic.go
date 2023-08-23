package bluetooth_tool

import (
	"context"
	"golang.org/x/sync/semaphore"
	"tinygo.org/x/bluetooth"
)

var sem = semaphore.NewWeighted(1)
var ctx = context.Background()

type BTCharacteristic struct {
	CheckCharacteristic *bluetooth.DeviceCharacteristic
	RWCharacteristic    *bluetooth.DeviceCharacteristic

	CheckRecvHandler RecvHandler
	RWRecvHandler    RecvHandler
}

func (btc *BTCharacteristic) SetCheckReceiveHandler(handler RecvHandler, recvChan chan<- []byte) {
	//err := sem.Acquire(ctx, 1)
	//if err != nil {
	//	return
	//}
	//defer sem.Release(1)
	//btc.CheckRecvHandler = handler
	//_ = btc.CheckCharacteristic.EnableNotifications(handler(recvChan))
}

func (btc *BTCharacteristic) SetReadWriteReceiveHandler(handler RecvHandler, recvChan chan<- []byte, btIn chan<- struct{}) {
	err := sem.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer sem.Release(1)
	btc.RWRecvHandler = handler
	_ = btc.RWCharacteristic.EnableNotifications(handler(recvChan, btIn))
}
func (btc *BTCharacteristic) SendDataToCheck(data []byte) error {
	err := sem.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	defer sem.Release(1)
	_, err = btc.CheckCharacteristic.Write(data)
	return err
}
func (btc *BTCharacteristic) SendDataToRW(data []byte) error {
	err := sem.Acquire(ctx, 1)
	if err != nil {
		return err
	}
	defer sem.Release(1)
	_, err = btc.RWCharacteristic.Write(data)
	return err
}
