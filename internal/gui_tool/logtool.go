package gui_tool

import (
	"sync"
	"time"
)
import "github.com/andlabs/ui"

var checked bool
var entry *ui.MultilineEntry
var mutex = &sync.Mutex{}

func AddLog(log string) {
	mutex.Lock()
	currData := time.Now().Format("2006-01-02 15:04:05")
	entry.Append(currData)
	entry.Append("\t")
	entry.Append(log)
	entry.Append("\n")
	mutex.Unlock()
}
func LogStatus(checkOutput bool) {
	checked = checkOutput
}
func CheckedLog() bool { return checked }
func SetEntry(entrylg *ui.MultilineEntry) {
	entry = entrylg
}
