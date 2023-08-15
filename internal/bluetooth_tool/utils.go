package bluetooth_tool

import "math/rand"

// 这个函数应该是用来和对讲机鉴权用的？
func GenerateSerialNum() []byte {
	bArr := make([]byte, 20)
	for i := 0; i < 20; i++ {
		if i < 5 {
			bArr[i] = 63
		} else {
			bArr[i] = byte(int(rand.Float64() * 99.0))
		}
	}
	return bArr
}
