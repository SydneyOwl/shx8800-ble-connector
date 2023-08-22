//go:build gui

package bluetooth_tool

import (
	"github.com/andlabs/ui"
)

func GuiRWRecvHandler(recvChan chan<- []byte, label *ui.Label) func(c []byte) {
	return func(c []byte) {
		if label.Visible() {
			label.Hide()
		} else {
			label.Show()
		}
		recvChan <- c
	}
}
