package main

import (
	"context"
	"errors"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gookit/slog"
	"github.com/sydneyowl/shx8800-ble-connector/config"
	"github.com/sydneyowl/shx8800-ble-connector/internal/bluetooth_tool"
	"github.com/sydneyowl/shx8800-ble-connector/internal/gui_tool"
	"github.com/sydneyowl/shx8800-ble-connector/internal/serial_tool"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/exceptions"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"tinygo.org/x/bluetooth"
)

var comList = make([]string, 0)
var btList = make([]string, 0)
var canceler context.CancelFunc
var globalDevList = make(map[string]bluetooth.ScanResult)
var btIn = make(chan struct{}, 100)
var btOut = make(chan struct{}, 100)
var checkCharacteristic, rwCharacteristic *bluetooth.DeviceCharacteristic = nil, nil

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
func chgColorOnRecv(ctx context.Context, btIn chan struct{}, btOut chan struct{}) {
	go func(ctx context.Context, btIn chan struct{}) {
		for {
			select {
			case <-btIn:
				//green
				btrx.SetColor(0.167882, 0.704918, 0.109562, 1)
				time.Sleep(time.Millisecond * 3)
				//white
				btrx.SetColor(0.99, 0.99, 0.98, 1)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, btIn)
	go func(ctx context.Context, btOut chan struct{}) {
		for {
			select {
			case <-btOut:
				bttx.SetColor(0.9672130, 0.0618894, 0.0312005, 1)
				time.Sleep(time.Millisecond * 3)
				bttx.SetColor(0.99, 0.99, 0.98, 1)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, btOut)

}
func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	butAndChoices := ui.NewHorizontalBox()
	butAndChoices.SetPadded(true)
	allChoices := ui.NewVerticalBox()
	allChoices.SetPadded(true)
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	combox = ui.NewEditableCombobox()
	combox.SetText("-----请选择-----")
	comscan = ui.NewButton("扫描端口")
	hbox.Append(combox, false)
	hbox.Append(comscan, false)
	allChoices.Append(hbox, false)
	hbox2 := ui.NewHorizontalBox()
	hbox2.SetPadded(true)
	btbox = ui.NewEditableCombobox()
	btbox.SetText("-----请选择-----")
	hbox2.Append(btbox, false)
	scanBut = ui.NewButton("扫描蓝牙")
	scanBut.Disable()
	hbox2.Append(scanBut, false)
	allChoices.Append(hbox2, false)
	hbox3 := ui.NewHorizontalBox()
	hbox3.SetPadded(true)
	bandRateLabel = ui.NewLabel("波特率：")
	hbox3.Append(bandRateLabel, false)
	chgBandrate = ui.NewEditableCombobox()
	chgBandrate.SetText("9600（如无需要请勿更改）")
	for _, v := range serial_tool.AVAILABLE_BAUDRATE {
		chgBandrate.Append(strconv.Itoa(v))
	}
	chgBandrate.OnChanged(baudrateCallback)
	hbox3.Append(chgBandrate, false)
	allChoices.Append(hbox3, false)
	butAndChoices.Append(allChoices, false)
	butAndChoices.Append(ui.NewVerticalSeparator(), false)
	startBox := ui.NewHorizontalBox()
	connButton = ui.NewButton("开始连接")
	logButton = ui.NewCheckbox("原始数据输出(可能减慢传输速率)")
	logButton.OnToggled(logButtonToggled)
	startBox.Append(connButton, false)
	startBox.Append(logButton, false)
	startBox.SetPadded(true)
	butAndChoices.Append(startBox, false)
	pgbar := ui.NewVerticalBox()
	scanStatus = ui.NewLabel("等待扫描...")
	pgbar.Append(scanStatus, false)
	bar = ui.NewProgressBar()
	pgbar.Append(bar, false)
	vbox.Append(butAndChoices, false)
	vbox.Append(pgbar, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	//status := ui.New()
	group := ui.NewGroup("信息")
	entry = ui.NewNonWrappingMultilineEntry()
	entry.SetReadOnly(true)
	gui_tool.SetEntry(entry)
	group.SetChild(entry)
	group.SetMargined(true)
	vbox.Append(group, true)
	btrx = ui.NewColorButton()
	bttx = ui.NewColorButton()
	bttx.Disable()
	btrx.Disable()
	bttx.SetColor(0.99, 0.99, 0.98, 1)
	btrx.SetColor(0.99, 0.99, 0.98, 1)
	butBox := ui.NewHorizontalBox()
	butBox.Enable()
	butBox.SetPadded(true)
	BTL1 := ui.NewLabel("BT-RX")
	butBox.Append(BTL1, false)
	butBox.Append(bttx, true)
	BTL2 := ui.NewLabel("BT-TX")
	butBox.Append(BTL2, false)
	butBox.Append(btrx, true)
	ee := ui.NewEntry()
	ee.SetText("--使用前请做好备份工作--")
	ee.Disable()
	ef := ui.NewEntry()
	ef.SetText("shx8800-ble-connector " + config.VER)
	ef.Disable()
	butBox.Append(ee, false)
	butBox.Append(ef, false)
	vbox.Append(butBox, false)
	comscan.OnClicked(pressComscan)
	scanBut.OnClicked(pressBtScan)
	connButton.OnClicked(clickConnectButton)
	return vbox
}
func updateBtConnStat(addr bluetooth.Address, ctx context.Context) {
	connChan := make(chan *bluetooth.Device)
	bar.SetValue(-1)
	defer connButton.Enable()
	defer bar.SetValue(0)
	go bluetooth_tool.ConnectByMacNoBlock(addr, connChan)
	conn = <-connChan
	if conn == nil {
		ui.MsgBoxError(mainwin,
			"错误",
			"连接蓝牙失败")
		return
	}
	//defer doGuiShutup()
	gui_tool.AddLog("连接成功：" + addr.String())
	bluetooth_tool.SetConnectStatus(true)
	gui_tool.AddLog("发现服务中...")
	services, err := conn.DiscoverServices(nil)
	if err != nil {
		gui_tool.AddLog("无法发现服务")
		ui.MsgBoxError(mainwin,
			"错误",
			"发现服务失败！"+err.Error())
		return
	}
	gui_tool.AddLog("发现特征中...")
	var _, manufacturer, model, firmware = make([]byte, 10), make([]byte, 20), make([]byte, 20), make([]byte, 20)
	for _, service := range services {
		chs, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			gui_tool.AddLog("无法发现特征")
			ui.MsgBoxError(mainwin,
				"错误",
				"发现特征失败！"+err.Error())
			return
		}
		slog.Trace(chs)
		for i, ch := range chs {
			if strings.Contains(ch.String(), bluetooth_tool.FIRMWARE_REVISION_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(firmware)
				gui_tool.AddLog("固件版本：" + string(firmware))
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.MANUFACTURER_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(manufacturer)
				gui_tool.AddLog("生产产商：" + string(manufacturer))
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.MODEL_NUMBER_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(model)
				gui_tool.AddLog("设备型号：" + string(model))
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
		gui_tool.AddLog("无法获取设备通道")
		ui.MsgBoxError(mainwin,
			"错误",
			"无法获取设备信息，请检查设备！")
		return
	}
	bluetooth_tool.CurrentDevice = &bluetooth_tool.BTCharacteristic{
		CheckCharacteristic: checkCharacteristic, //暂时No
		RWCharacteristic:    rwCharacteristic,
	}
	btReplyChan := make(chan []byte, 5)
	serialChan := make(chan []byte, 10)
	errChan := make(chan error, 3)
	bluetooth_tool.CurrentDevice.SetReadWriteReceiveHandler(bluetooth_tool.RWRecvHandler, btReplyChan, btOut)
	go bluetooth_tool.BTWriter(ctx, serialChan, errChan, btIn)
	go serial_tool.SerialDataProvider(ctx, serialChan)
	go serial_tool.SerialDataWriter(ctx, btReplyChan, errChan)
	chgColorOnRecv(ctx, btIn, btOut)
	gui_tool.AddLog("初始化完成！现在可以连接写频软件了！")
	gui_tool.AddLog("如果遇到读频卡在4%，请点击取消后重新读频即可！\n手台写频完成重启后请重新连接！")
	gui_tool.AddLog("如果一直写频失败，请使用写频线写入")
	go func() {
		err := <-errChan
		if errors.Is(err, exceptions.TransferDone) {
			gui_tool.AddLog("传输完成！对讲机将重启，您可以退出了...")
			ui.MsgBox(mainwin,
				"提醒",
				"传输完成！对讲机将重启！")
		} else {
			gui_tool.AddLog("出现异常，如果对讲机写频完成后重启了或者对讲机已经被关闭，您可以忽略提示" + err.Error())
			ui.MsgBoxError(mainwin,
				"注意",
				"出现异常，如果对讲机写频完成后重启了或者对讲机已经被关闭，您可以忽略提示")
		}
		doGuiShutup()
	}()
	bar.SetValue(0)
	connButton.SetText("断开连接并退出")
	connButton.Enable()
	<-ctx.Done()
}
func updateComboBt() {
	devList := make([]bluetooth.ScanResult, 0)
	errChan := make(chan error)
	bar.SetValue(-1)
	scanStatus.SetText("扫描中，请等待...")
	defer func() {
		scanBut.Enable()
		connButton.Enable()
		bar.SetValue(0)
		scanStatus.SetText("扫描结束")
	}()
	go bluetooth_tool.GetAvailableBtDevListViaChannel(bluetooth_tool.SHX8800Filter, &devList, errChan)
	for {
		err, ok := <-errChan
		if !ok {
			break
		}
		if err != nil {
			ui.MsgBoxError(mainwin,
				"错误",
				"扫描蓝牙出现错误："+err.Error()+", 原因我也不知道...重新启动程序即可解决！按下确定重重启程序")
			doRestart()
			return
		}
	}
	for i := range devList {
		if !contains(btList, devList[i].Address.String()) {
			globalDevList[devList[i].Address.String()] = devList[i]
			btbox.Append("[" + devList[i].LocalName() + "]" + devList[i].Address.String())
			comList = append(comList, devList[i].Address.String())
		}
	}
	if len(devList) == 0 {
		btbox.SetText("未找到设备")
	} else {
		btbox.SetText("[" + devList[0].LocalName() + "]" + devList[0].Address.String())
	}
}
func setupUI() {
	mainwin = ui.NewWindow("森海克斯写频工具", 550, 400, false)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	mainwin.SetMargined(true)
	mainwin.SetChild(makeBasicControlsPage())
	mainwin.Show()
}
func GUI() {
	//defer doGuiShutup()
	_ = ui.Main(setupUI)
}
func doGuiShutup() {
	if bluetooth_tool.GetConnectStatus() {
		bluetooth_tool.DisconnectDevice(conn)
	}
	if serial_tool.GetConnectedStatus() {
		serial_tool.ShutPort()
	}
	gui_tool.AddLog("连接断开！")
	gui_tool.AddLog("清理...")
	if canceler != nil {
		canceler()
	}
	time.Sleep(time.Second)
	os.Exit(0)
	//connButton.SetText("开始连接")
	//connButton.Enable()
}
func doRestart() {
	defer doGuiShutup()
	exePath, _ := os.Executable()
	go exec.Command(exePath, os.Args...).Run()
	time.Sleep(time.Second * 2)
}
