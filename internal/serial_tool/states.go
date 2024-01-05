package serial_tool

const (
	SERIAL_READ_TIMEOUT_MS = 20
)

var (
	AVAILABLE_BAUDRATE = []int{2400, 4800, 9600, 19200, 38400, 57600, 115200}
	BAUDRATE           = 9600
)

const (
	HandShakeStep1 = iota
	HandShakeStep2
	HandShakeStep3
	HandShakeStep4
	ReadStep1
	ReadStep2
	ReadStep3
	WriteStep1
	WriteStep
)

var connected = false

func SetConnectedStatus(stat bool) {
	connected = stat
}
func GetConnectedStatus() bool {
	return connected
}
