package logs

import (
	"testing"
)

// Try each log level in decreasing order of priority.
func testConsoleCalls(bl *ULogger) {
	bl.Emergency("emergency")
	bl.Alert("alert")
	bl.Critical("critical")
	bl.Error("error")
	bl.Warning("warning")
	bl.Notice("notice")
	bl.Informational("informational")
	bl.Debug("debug")
}

// Test console logging by visually comparing the lines being output with and
// without a log level specification.
func TestConsole(t *testing.T) {
	log1 := NewLogger(10000)
	log1.EnableFuncCallDepth(true)
	log1.SetLogger("console", "")
	testConsoleCalls(log1)

	log1.Warning("==============================================")
	log2 := NewLogger(100)
	log2.EnableFuncCallDepth(true)
	log2.SetLogger("console", `{"level":0}`)
	testConsoleCalls(log2)

	log3 := NewLogger(100)
	log3.SetLogger(AdapterFile,`{"filename":"logs/error.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	log3.SetLevel(3)     // 设置日志写入缓冲区的等级
	log3.EnableFuncCallDepth(true)
	log3.SetLogger("console", `{"level":0}`)
	testConsoleCalls(log3)
}


