package main

import (
	"context"
	"errors"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gookit/slog"
	"github.com/sydneyowl/shx8800-ble-connector/internal/bluetooth_tool"
	"github.com/sydneyowl/shx8800-ble-connector/internal/serial_tool"
	"github.com/sydneyowl/shx8800-ble-connector/pkg/exceptions"
	"os"
	"strings"
	"time"
	"tinygo.org/x/bluetooth"
)

var btbox *ui.EditableCombobox
var mainwin *ui.Window
var entry *ui.MultilineEntry
var comList = make([]string, 0)
var btList = make([]string, 0)
var conn *bluetooth.Device
var scanBut *ui.Button
var comscan *ui.Button
var connButton *ui.Button
var bar *ui.ProgressBar
var scanStatus *ui.Label
var canceler context.CancelFunc
var globalDevList = make(map[string]bluetooth.ScanResult, 0)

var btrx *ui.Label
var bttx *ui.Label
var serx *ui.Label
var setx *ui.Label

var checkCharacteristic, rwCharacteristic *bluetooth.DeviceCharacteristic = nil, nil

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
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
	combox := ui.NewEditableCombobox()
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
	butAndChoices.Append(allChoices, false)
	butAndChoices.Append(ui.NewVerticalSeparator(), false)
	startBox := ui.NewHorizontalBox()
	connButton = ui.NewButton("开始连接")
	startBox.Append(connButton, false)
	butAndChoices.Append(startBox, false)
	lightBox := ui.NewVerticalBox()
	box_upper := ui.NewVerticalBox()
	btrx = ui.NewLabel("↓bt_rx")
	bttx = ui.NewLabel("↑bt_tx")
	box_upper.Append(btrx, false)
	box_upper.Append(bttx, false)
	box_upper.Append(ui.NewHorizontalSeparator(), false)
	box_sup := ui.NewVerticalBox()
	serx = ui.NewLabel("↓serial_rx")
	setx = ui.NewLabel("↑serial_tx")
	box_sup.Append(serx, false)
	box_sup.Append(setx, false)
	lightBox.Append(box_upper, false)
	lightBox.Append(box_sup, false)
	butAndChoices.Append(lightBox, false)
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
	group.SetMargined(true)
	vbox.Append(group, true)
	entry = ui.NewNonWrappingMultilineEntry()
	entry.SetReadOnly(true)
	group.SetChild(entry)

	comscan.OnClicked(func(button *ui.Button) {
		button.Disable()
		connButton.Disable()
		defer button.Enable()
		defer connButton.Enable()
		defer scanBut.Enable()
		ports, err := serial_tool.ScanPort()
		if err != nil {
			ui.MsgBoxError(mainwin,
				"错误",
				"扫描端口出现错误："+err.Error())
			return
		}
		for i := range ports {
			if !contains(comList, ports[i]) {
				combox.Append(ports[i])
				comList = append(comList, ports[i])
			}
		}
		if len(comList) == 0 {
			combox.SetText("未找到设备")
		} else {
			combox.SetText(ports[0])
		}
	})
	scanBut.OnClicked(func(button *ui.Button) {
		button.Disable()
		connButton.Disable()
		go updateComboBt()
	})
	connButton.OnClicked(func(button *ui.Button) {
		button.Disable()
		if button.Text() == "开始连接" {
			com := combox.Text()
			bt := btbox.Text()
			if com == "未找到设备" || bt == "未找到设备" || com == "-----请选择-----" || bt == "-----请选择-----" {
				ui.MsgBoxError(mainwin,
					"错误",
					"请选择正确的选项")
				addLog("未找到设备！")
				button.Enable()
				return
			}
			err := serial_tool.ConnPort(com)
			if err != nil {
				ui.MsgBoxError(mainwin,
					"错误",
					"连接端口失败:"+err.Error())
				addLog("端口连接失败！")
				button.Enable()
				return
			}
			addr := strings.Split(bt, "]")
			mac := addr[len(addr)-1]
			ctx, cancel := context.WithCancel(context.Background())
			canceler = cancel
			go updateBtConnStat(globalDevList[mac].Address, ctx)
			//go bluetooth_tool.ConnectByMacNoBlock(globalDevList[mac].Address, connChan)
		} else {
			doGuiShutup()
		}
	})
	return vbox
}
func doGuiShutup() {
	bluetooth_tool.DisconnectDevice(conn)
	serial_tool.ShutPort()
	addLog("连接断开！")
	addLog("清理...")
	canceler()
	time.Sleep(time.Second)
	os.Exit(0)
	//connButton.SetText("开始连接")
	//connButton.Enable()
}
func addLog(log string) {
	currData := time.Now().Format("2006-01-02 15:04:05")
	entry.Append(currData)
	entry.Append("\t")
	entry.Append(log)
	entry.Append("\n")
}
func updateBtConnStat(addr bluetooth.Address, ctx context.Context) {
	connChan := make(chan *bluetooth.Device, 0)
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
	defer doGuiShutup()
	addLog("连接成功：" + addr.String())
	addLog("发现服务中...")
	services, err := conn.DiscoverServices(nil)
	if err != nil {
		addLog("无法发现服务")
		ui.MsgBoxError(mainwin,
			"错误",
			"发现服务失败！")
		return
	}
	addLog("发现特征中...")
	var _, manufacturer, model, firmware = make([]byte, 10), make([]byte, 20), make([]byte, 20), make([]byte, 20)
	for _, service := range services {
		chs, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			addLog("无法发现特征")
			ui.MsgBoxError(mainwin,
				"错误",
				"发现特征失败！")
			return
		}
		slog.Trace(chs)
		for i, ch := range chs {
			if strings.Contains(ch.String(), bluetooth_tool.FIRMWARE_REVISION_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(firmware)
				addLog("固件版本：" + string(firmware))
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.MANUFACTURER_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(manufacturer)
				addLog("生产产商：" + string(manufacturer))
				continue
			}
			if strings.Contains(ch.String(), bluetooth_tool.MODEL_NUMBER_CHARACTERISTIC_UUID) {
				_, _ = ch.Read(model)
				addLog("设备型号：" + string(model))
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
		addLog("无法获取设备通道")
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
	bluetooth_tool.CurrentDevice.SetGuiReadWriteReceiveHandler(bluetooth_tool.GUIRWRecvHandler, btReplyChan, btrx)
	go bluetooth_tool.GuiBTWriter(ctx, serialChan, errChan, bttx)
	go serial_tool.GuiSerialDataProvider(ctx, serialChan, serx)
	go serial_tool.GuiSerialDataWriter(ctx, btReplyChan, errChan, setx)
	addLog("初始化完成！现在可以连接写频软件了！")
	addLog("如果遇到读频卡在4%，请点击取消后重新读频即可！\n手台写频完成重启后请重新连接！")
	addLog("如果一直写频失败，请使用写频线写入")
	go func() {
		err := <-errChan
		if errors.Is(err, exceptions.TransferDone) {
			addLog("传输完成！对讲机将重启，您可以退出了...")
			ui.MsgBox(mainwin,
				"提醒",
				"传输完成！对讲机将重启！")
		} else {
			addLog("出现异常，如果对讲机写频完成后重启了或者对讲机已经被关闭，您可以忽略提示" + err.Error())
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
	errChan := make(chan error, 0)
	bar.SetValue(-1)
	scanStatus.SetText("扫描中，请等待...")
	defer func() {
		scanBut.Enable()
		connButton.Enable()
		bar.SetValue(0)
		scanStatus.SetText("扫描结束")
	}()
	go bluetooth_tool.GetAvailableBtDevListViaChannel(bluetooth_tool.SHX8800Filter, errChan, &devList)
	for {
		err, ok := <-errChan
		if !ok {
			break
		}
		if err != nil {
			ui.MsgBoxError(mainwin,
				"错误",
				"扫描蓝牙出现错误："+err.Error())
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
	mainwin = ui.NewWindow("森海克斯写频工具", 0, 400, false)
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
