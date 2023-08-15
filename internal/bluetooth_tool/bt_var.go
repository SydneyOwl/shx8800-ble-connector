package bluetooth_tool

import "sync"

var (
	CurrentDevice *BTCharacteristic
	BtMutex       = sync.Mutex{}
)
