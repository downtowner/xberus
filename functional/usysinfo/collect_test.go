package usysinfo

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestMemoryInfoCollect(t *testing.T) {
	memory := MemoryInfo{}
	info := memory.Collect()

	if info.Total == 0 {
		t.Error("Memory collect error")
	}
}
func TestCpuInfoCollect(t *testing.T) {
	cpu := CpuInfo{}
	cpu.Collect()

}

func TestDiskInfoCollect(t *testing.T) {
	disk := DiskInfo{}
	info := disk.Collect()

	if len(info) == 0 {
		t.Error("disk collect error")
	}
}

func TestNetworkInfoCollect(t *testing.T) {
	network := NewWorkInfo{}
	info := network.Collect()

	if len(info) == 0 {
		t.Error("network collect error")
	}
}

func TestGoruntimeInfoCollect(t *testing.T) {
	goRuntime := GoRuntimeInfo{}
	info := goRuntime.Collect()

	if info.NumGoroutine == 0 {
		t.Error("GoruntimeInfo collect error")
	}
}

func TestIOInfo(t *testing.T) {
	ioInfo := IOInfo{}
	iodata := ioInfo.Collect()
	t.Log(iodata)

	filePath := "./test.txt"
	for i := 0; i < 10; i++ {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			t.Error("open file error", err)
		}
		write := bufio.NewWriter(file)
		write.WriteString("test \n")
		write.Flush()
		file.Close()
		time.Sleep(time.Second)
	}

	iodata = ioInfo.Collect()
	if iodata.WriteByte == 0 {
		t.Error("iostat.Collect is error")
	}
	os.Remove(filePath)
	t.Log(iodata.ReadByte)
	t.Log(iodata.ReadCount)
	t.Log(iodata.WriteByte)
	t.Log(iodata.WriteCount)
}
