package serial_tool

const (
	BAUDRATE               = 9600
	SERIAL_READ_TIMEOUT_MS = 20
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
