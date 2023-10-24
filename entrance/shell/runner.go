package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gookit/slog"
	"github.com/sydneyowl/shx8800-ble-connector/internal/bluetooth_tool"
	"github.com/sydneyowl/shx8800-ble-connector/internal/serial_tool"
	"github.com/sydneyowl/shx8800-ble-connector/internal/stdout_fmt"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/exceptions"
	"os"
	"strings"
	"sync"
	"tinygo.org/x/bluetooth"
)

var cancelFunc context.CancelFunc
var once sync.Once
var btConn *bluetooth.Device

func doShutup() {
	// Clean up...
	cancelFunc()
	bluetooth_tool.DisconnectDevice(btConn)
	serial_tool.ShutPort()
}

func StartAndRun() {
	defer func() {
		once.Do(doShutup)
		if err := recover(); err != nil {
			slog.Fatalf("出现异常：%v. 按下回车键退出", err)
			_, _ = fmt.Scanln()
		}
	}()
	// port
	slog.Info("搜索端口...")
	ports, err := serial_tool.ScanPort()
	if err != nil {
		slog.Fatalf("无法检测端口：%v", err)
		_, _ = fmt.Scanln()
		os.Exit(-1)
	}
	stdout_fmt.PrintAllPorts(ports)
	fmt.Println()
	fmt.Print("请输入目标端口编号: ")
	var portNo int
	_, _ = fmt.Scanln(&portNo)
	for !(portNo > 0 && portNo <= len(ports)) {
		slog.Warnf("输入错误，请重新输入")
		fmt.Print("请输入目标端口编号:")
		_, _ = fmt.Scanln(&portNo)
	}
	targetPort := ports[portNo-1]
	err = serial_tool.ConnPort(targetPort)
	if err != nil {
		slog.Fatalf("无法连接端口：%v", err)
		_, _ = fmt.Scanln()
		os.Exit(-1)
	}
	slog.Info("端口连接成功！")

	// Device-BT
	slog.Infof("正在扫描设备,请等待%d秒...", bluetooth_tool.BT_SCAN_TIMEOUT)
	list, err := bluetooth_tool.GetAvailableBtDevList(bluetooth_tool.SHX8800Filter)
	if err != nil {
		slog.Fatalf("扫描失败: %v", err)
		_, _ = fmt.Scanln()
		os.Exit(-1)
	}
	if len(list) == 0 {
		slog.Fatal("未找到SHX8800设备！")
		_, _ = fmt.Scanln()
		os.Exit(-1)
	}
	stdout_fmt.PrintAvailableShxDevices(list)
	var deviceNo int
	fmt.Println()
	fmt.Print("请输入要连接的设备编号: ")
	_, _ = fmt.Scanln(&deviceNo)
	for !(deviceNo > 0 && deviceNo <= len(list)) {
		slog.Warnf("输入错误，请重新输入")
		fmt.Print("请输入要连接的设备编号: ")
		_, _ = fmt.Scanln(&deviceNo)
	}
	deviceSHX := list[deviceNo-1]
	slog.Info("连接设备...")
	//windows不可用
	//bluetooth_tool.SetHandler(bluetooth_tool.DisconnectHandler)
	conn, err := bluetooth_tool.ConnectByMac(deviceSHX.Address)
	btConn = conn
	if err != nil {
		slog.Fatalf("无法连接设备:%v", err)
		_, _ = fmt.Scanln()
		os.Exit(-1)
	}
	slog.Info("连接成功！")
	slog.Debug("正在发现服务...")
	services, err := conn.DiscoverServices(nil)
	if err != nil {
		slog.Fatalf("无法发现服务:%v", err)
		_, _ = fmt.Scanln()
		os.Exit(-1)
	}
	slog.Trace(services)
	slog.Debug("正在发现特征...")
	var _, manufacturer, model, firmware = make([]byte, 10), make([]byte, 20), make([]byte, 20), make([]byte, 20)
	var checkCharacteristic, rwCharacteristic *bluetooth.DeviceCharacteristic = nil, nil
	for _, service := range services {
		chs, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			slog.Fatalf("无法发现特征:%v", err)
			_, _ = fmt.Scanln()
			os.Exit(-1)
		}
		slog.Trace(chs)
		for i, ch := range chs {
			if strings.Contains(ch.String(), bluetooth_tool.FIRMWARE_REVISION_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(firmware)
				slog.Noticef("固件版本：%s", string(firmware))
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.MANUFACTURER_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(manufacturer)
				slog.Noticef("生产产商：%s", string(manufacturer))
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.MODEL_NUMBER_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(model)
				slog.Noticef("设备型号：%s", string(model))
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.CHECK_CHARACTERISTIC_UUID) {
				checkCharacteristic = &chs[i]
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.RW_CHARACTERISTIC_UUID) {
				rwCharacteristic = &chs[i]
				continue
			}
		}
	}
	if checkCharacteristic == nil || rwCharacteristic == nil {
		slog.Fatalf("无法获取设备通道")
		_, _ = fmt.Scanln()
		os.Exit(-1)
	}
	bluetooth_tool.CurrentDevice = &bluetooth_tool.BTCharacteristic{
		CheckCharacteristic: checkCharacteristic, //暂时No
		RWCharacteristic:    rwCharacteristic,
	}
	btReplyChan := make(chan []byte, 5)
	serialChan := make(chan []byte, 10)
	errChan := make(chan error, 3)
	bluetooth_tool.CurrentDevice.SetReadWriteReceiveHandler(bluetooth_tool.RWRecvHandler, btReplyChan, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancelFunc = cancel
	go bluetooth_tool.BTWriter(ctx, serialChan, errChan, nil)
	go serial_tool.SerialDataProvider(ctx, serialChan)
	go serial_tool.SerialDataWriter(ctx, btReplyChan, errChan)
	slog.Notice("初始化完成！现在可以连接写频软件了！输入任意字符退出软件！")
	slog.Notice("如果遇到读频卡在4%，请点击取消后重新读频即可！手台写频完成重启后请关闭软件重新打开！")
	slog.Notice("如果一直写频失败，请使用写频线写入")
	slog.Noticef("------------------------------")
	go func() {
		err := <-errChan
		if errors.Is(err, exceptions.TransferDone) {
			slog.Noticef("传输完成！对讲机将重启，按下回车键退出...")
		} else {
			slog.Warnf("出现异常：%v，如果对讲机写频完成后重启了或者对讲机已经被关闭，您可以忽略提示，按下回车键以退出", err)
			once.Do(doShutup)
		}
	}()
	_, _ = fmt.Scanln()
	once.Do(doShutup)
}
