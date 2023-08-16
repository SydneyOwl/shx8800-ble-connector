package test

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/sydneyowl/shx8800-ble-connector/internal/bluetooth_tool"
	"github.com/sydneyowl/shx8800-ble-connector/internal/stdout_fmt"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/logger"
	"testing"
)

func TestBtScanner(t *testing.T) {
	logger.InitLog(false, false)
	list, err := bluetooth_tool.GetAvailableBtDevList(bluetooth_tool.SHX8800Filter)
	if err != nil {
		slog.Fatalf("出错：%v", err)
		_, _ = fmt.Scanln()
	}
	if len(list) == 0 {
		slog.Warn("未找到SHX8800设备！")
		t.FailNow()
	}
	stdout_fmt.PrintAvailableShxDevices(list)
	var deviceNo int
	fmt.Println()
	fmt.Print("请输入要连接的设备编号")
	_, _ = fmt.Scanln(&deviceNo)
	slog.Info(deviceNo)
	deviceSHX := list[deviceNo]
	slog.Info(deviceSHX)
}
