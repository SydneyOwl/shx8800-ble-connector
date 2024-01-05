package bluetooth_tool

var (
	call_disconnect = false
	connected       = false
)

func SetConnectStatus(connectStat bool) {
	connected = connectStat
}

func GetConnectStatus() bool {
	return connected
}
