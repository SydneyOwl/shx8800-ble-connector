package main

import (
	"context"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/sydneyowl/shx8800-ble-connector/internal/gui_tool"
	"github.com/sydneyowl/shx8800-ble-connector/internal/serial_tool"
	"strings"
	"tinygo.org/x/bluetooth"
)

var bttx *ui.ColorButton
var btrx *ui.ColorButton
var btbox *ui.EditableCombobox
var mainwin *ui.Window
var entry *ui.MultilineEntry
var conn *bluetooth.Device
var scanBut *ui.Button
var comscan *ui.Button
var combox *ui.EditableCombobox
var logButton *ui.Checkbox
var connButton *ui.Button
var bar *ui.ProgressBar
var scanStatus *ui.Label

func clickConnectButton(button *ui.Button) {
	button.Disable()
	if button.Text() == "开始连接" {
		com := combox.Text()
		bt := btbox.Text()
		if com == "未找到设备" || bt == "未找到设备" || com == "-----请选择-----" || bt == "-----请选择-----" {
			ui.MsgBoxError(mainwin,
				"错误",
				"请选择正确的选项")
			gui_tool.AddLog("未找到设备！")
			button.Enable()
			return
		}
		err := serial_tool.ConnPort(com)
		if err != nil {
			ui.MsgBoxError(mainwin,
				"错误",
				"连接端口失败:"+err.Error())
			gui_tool.AddLog("端口连接失败！")
			button.Enable()
			return
		}
		serial_tool.SetConnectedStatus(true)
		addr := strings.Split(bt, "]")
		mac := addr[len(addr)-1]
		ctx, cancel := context.WithCancel(context.Background())
		canceler = cancel
		go updateBtConnStat(globalDevList[mac].Address, ctx)
		//go bluetooth_tool.ConnectByMacNoBlock(globalDevList[mac].Address, connChan)
	} else {
		doGuiShutup()
	}
}

func pressBtScan(button *ui.Button) {
	button.Disable()
	connButton.Disable()
	go updateComboBt()
}
func pressComscan(button *ui.Button) {
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
}

func logButtonToggled(checkbox *ui.Checkbox) {
	gui_tool.LogStatus(checkbox.Checked())
}
