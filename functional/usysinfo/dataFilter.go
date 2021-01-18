package usysinfo

import (
	"time"
)

type DataFilter struct {
	CpuFilter        func(respData []ChartPoints, data CpuData, timeInterval time.Duration) []ChartPoints
	MemoryFilter     func(respData []ChartPoints, data MemoryData, timeInterval time.Duration) []ChartPoints
	GoRuntimeFilter  func(respData []ChartPoints, data GoRuntimeData, timeInterval time.Duration) []ChartPoints
	IOFilter         func(respData []ChartPoints, data IOData, timeInterval time.Duration) []ChartPoints
	NetworkMapFilter func(respData []ChartPoints, data *NetworkData, timeInterval time.Duration, mapKey string) []ChartPoints
	DiskMapFilter    func(respData []ChartPoints, data *DiskData, timeInterval time.Duration, mapKey string) []ChartPoints
}

//图表坐标点
type ChartPoints struct {
	Name  string `json:"n"`
	Value uint64 `json:"v"`
	Time  string `json:"t"`
}

//根据时间间隔过滤
func FilterByTimeInterval(label string, collectData []OneCollectData, timeInterval time.Duration) []ChartPoints {
	chartPoints := make([]ChartPoints, 0)
	lastTime := time.Unix(0, 0)
	for _, v := range collectData {
		vTime := time.Unix(int64(v.Time), 0)
		if lastTime.Add(timeInterval).After(vTime) {
			continue
		}
		lastTime = vTime
		chartPoints = append(chartPoints, ChartPoints{
			Name:  label,
			Value: v.Metrics,
			Time:  FormatTime(vTime),
		})
	}
	return chartPoints
}

var (
	cpuField = map[string]DataFilter{
		cpuUsedPercent: DataFilter{
			CpuFilter: func(respData []ChartPoints, cpuData CpuData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("cpu使用率", cpuData.UsedPercent, timeInterval)...)
			},
		},
	}
	memoryField = map[string]DataFilter{
		memoryUsedPercent: DataFilter{
			MemoryFilter: func(respData []ChartPoints, memoryData MemoryData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("内存使用率", memoryData.UsedPercent, timeInterval)...)
			},
		},
	}
	goRuntimeField = map[string]DataFilter{
		goNumCpu: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("cpu.numCpu", goRuntimeData.NumCpu, timeInterval)...)
			},
		},
		goNumGoroutine: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("cpu.goroutines", goRuntimeData.NumGoroutine, timeInterval)...)
			},
		},
		goNumCgoCall: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("cpu.cgoCall", goRuntimeData.NumCgoCall, timeInterval)...)
			},
		},
		goAlloc: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.alloc", goRuntimeData.Alloc, timeInterval)...)
			},
		},
		goTotalAlloc: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.totalAlloc", goRuntimeData.TotalAlloc, timeInterval)...)
			},
		},
		goSys: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.sys", goRuntimeData.Sys, timeInterval)...)
			},
		},
		goLookups: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.lookups", goRuntimeData.Lookups, timeInterval)...)
			},
		},
		goMallocs: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.mallocs", goRuntimeData.Mallocs, timeInterval)...)
			},
		},
		goFrees: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.frees", goRuntimeData.Frees, timeInterval)...)
			},
		},
		goHeapAlloc: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.heap.alloc", goRuntimeData.HeapAlloc, timeInterval)...)
			},
		},
		goHeapSys: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.heap.sys", goRuntimeData.HeapSys, timeInterval)...)
			},
		},
		goHeapIdle: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.heap.idle", goRuntimeData.HeapIdle, timeInterval)...)
			},
		},
		goHeapInuse: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.heap.inuse", goRuntimeData.HeapInuse, timeInterval)...)
			},
		},
		goHeapReleased: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.heap.released", goRuntimeData.HeapReleased, timeInterval)...)
			},
		},
		goHeapObjects: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.heap.objects", goRuntimeData.HeapObjects, timeInterval)...)
			},
		},
		goStackInuse: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.stack.inuse", goRuntimeData.StackInuse, timeInterval)...)
			},
		},
		goStackSys: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.stack.sys", goRuntimeData.StackSys, timeInterval)...)
			},
		},
		goMSpanInuse: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.stack.mspanInuse", goRuntimeData.MSpanInuse, timeInterval)...)
			},
		},
		goMSpanSys: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.stack.mspanSys", goRuntimeData.MSpanSys, timeInterval)...)
			},
		},
		goMCacheInuse: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.stack.mcacheInuse", goRuntimeData.MCacheInuse, timeInterval)...)
			},
		},
		goMCacheSys: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.stack.mcacheSys", goRuntimeData.MCacheSys, timeInterval)...)
			},
		},
		goOtherSys: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.otherSys", goRuntimeData.OtherSys, timeInterval)...)
			},
		},
		goNumGC: DataFilter{
			GoRuntimeFilter: func(respData []ChartPoints, goRuntimeData GoRuntimeData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("mem.gc.numGC", goRuntimeData.NumGC, timeInterval)...)
			},
		},
	}
	diskField = map[string]DataFilter{
		diskUsedPercent: DataFilter{
			DiskMapFilter: func(respData []ChartPoints, diskData *DiskData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"硬盘使用率", diskData.UsedPercent, timeInterval)...)
			},
		},
		diskInodesUsedPercent: DataFilter{
			DiskMapFilter: func(respData []ChartPoints, diskData *DiskData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"硬盘Inodes使用率", diskData.InodesUsedPercent, timeInterval)...)
			},
		},
	}
	networkField = map[string]DataFilter{
		networkRecvByteAvg: DataFilter{
			NetworkMapFilter: func(respData []ChartPoints, networkData *NetworkData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"每秒接收字节数", networkData.RecvByteAvg, timeInterval)...)
			},
		},
		networkSentByteAvg: DataFilter{
			NetworkMapFilter: func(respData []ChartPoints, networkData *NetworkData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"每秒发送字节数", networkData.SentByteAvg, timeInterval)...)
			},
		},
		networkRecvErrRate: DataFilter{
			NetworkMapFilter: func(respData []ChartPoints, networkData *NetworkData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"收包错误率", networkData.RecvErrRate, timeInterval)...)
			},
		},
		networkSentErrRate: DataFilter{
			NetworkMapFilter: func(respData []ChartPoints, networkData *NetworkData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"发包错误率", networkData.SentErrRate, timeInterval)...)
			},
		},
		networkRecvPkgAvg: DataFilter{
			NetworkMapFilter: func(respData []ChartPoints, networkData *NetworkData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"正确收包数", networkData.RecvPkgAvg, timeInterval)...)
			},
		},
		networkSentPkgAvg: DataFilter{
			NetworkMapFilter: func(respData []ChartPoints, networkData *NetworkData, timeInterval time.Duration, mapKey string) []ChartPoints {
				return append(respData, FilterByTimeInterval(mapKey+"正确发包数", networkData.SentPkgAvg, timeInterval)...)
			},
		},
	}
	ioField = map[string]DataFilter{
		ioReadCount: DataFilter{
			IOFilter: func(respData []ChartPoints, ioData IOData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("硬盘IO每秒读次数", ioData.ReadCount, timeInterval)...)
			},
		},
		ioWriteCount: DataFilter{
			IOFilter: func(respData []ChartPoints, ioData IOData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("硬盘IO每秒写次数", ioData.WriteCount, timeInterval)...)
			},
		},
		ioReadByte: DataFilter{
			IOFilter: func(respData []ChartPoints, ioData IOData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("硬盘IO每秒读字节数", ioData.ReadByte, timeInterval)...)
			},
		},
		ioWriteByte: DataFilter{
			IOFilter: func(respData []ChartPoints, ioData IOData, timeInterval time.Duration) []ChartPoints {
				return append(respData, FilterByTimeInterval("硬盘IO每秒写字节数", ioData.WriteByte, timeInterval)...)
			},
		},
	}
)
