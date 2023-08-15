package stdout_fmt

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"tinygo.org/x/bluetooth"
)

func PrintAvailableShxDevices(shxs []bluetooth.ScanResult) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"No.", "Mac", "信号", "名称", "其他信息"})
	for i, v := range shxs {
		t.AppendRow(table.Row{
			i + 1, v.Address, v.RSSI, v.LocalName(), v.Address.IsRandom(),
		})
	}
	t.Render()
}
func PrintAllPorts(ports []string) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"No.", "端口号", "其他信息"})
	for i, v := range ports {
		t.AppendRow(table.Row{
			i + 1, v, "",
		})
	}
	t.Render()
}
