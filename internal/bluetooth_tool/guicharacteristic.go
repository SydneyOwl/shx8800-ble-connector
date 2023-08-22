//go:build gui

package bluetooth_tool

import "github.com/andlabs/ui"

type GUIRecvHandler func(recvChan chan<- []byte, label *ui.Label) func(c []byte)

func (btc *BTCharacteristic) SetGuiReadWriteReceiveHandler(handler GUIRecvHandler, recvChan chan<- []byte, label *ui.Label) {
	err := sem.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer sem.Release(1)
	_ = btc.RWCharacteristic.EnableNotifications(handler(recvChan, label))
}
