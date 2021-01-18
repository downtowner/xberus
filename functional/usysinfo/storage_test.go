package usysinfo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

func TestStorage(t *testing.T) {
	obj := createObj()
	storage := obj.NewStorage()
	err := storage.InitEngine()
	if err != nil {
		t.Fatal(err)
	}
	hostInfo, err := host.InfoWithContext(context.TODO())
	if err != nil {
		t.Fatal(err)
	}
	cpuinfo := CpuInfo{}
	memoryInfo := MemoryInfo{}
	newWorkInfo := NewWorkInfo{}
	diskInfo := DiskInfo{}
	goRuntimeInfo := GoRuntimeInfo{}
	testNowTime := time.Now()
	storageInfo := StorageInfo{
		Time:      testNowTime.AddDate(0, 0, -7).Unix(),
		Hostname:  hostInfo.Hostname,
		Cpu:       cpuinfo.Collect(),
		Memory:    memoryInfo.Collect(),
		NetWork:   newWorkInfo.Collect(),
		Disk:      diskInfo.Collect(),
		GoRuntime: goRuntimeInfo.Collect(),
	}
	storage.Set(storageInfo)
	storageInfo.Time = testNowTime.Unix() + 1
	storage.Set(storageInfo)
	storageInfo.Time = testNowTime.Unix() + 1
	storage.Set(storageInfo)

	startdate := testNowTime
	enddate := startdate.AddDate(0, 0, 1)
	cpuData := storage.GetCpu(startdate, enddate)
	if len(cpuData.UsedPercent) != 2 {
		t.Error("storage.GetCpu error")
	}
	memoryData := storage.GetMemory(startdate, enddate)
	if len(memoryData.UsedPercent) != 2 {
		t.Error("storage.GetMemory error")
	}
	goRuntime := storage.GetGoRuntime(startdate, enddate)
	if len(goRuntime.NumGoroutine) != 2 {
		t.Error("storage.GetGoRuntime error")
	}
	disk := storage.GetDisk(startdate, enddate)
	for k, _ := range disk {
		if len(disk[k].UsedPercent) != 2 {
			t.Error("storage.GetDisk error")
		}
		break
	}
	network := storage.GetNetwork(startdate, enddate)
	for k, _ := range network {
		if len(network[k].RecvByteAvg) != 2 {
			t.Error("storage.GetNetwork error")
		}
		break
	}
	os.Remove(DBPath)

}
