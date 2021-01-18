package usysinfo

import (
	"context"
	"net"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	utilNet "github.com/shirou/gopsutil/v3/net"
)

var IfiMap = map[string]*Ifi{}
var IOStat = &IOStatLog{}

type Ifi struct {
	Name string
	// Ip                string
	// Speed             float64
	// OutRecvPkgErrRate float64 //外网收包错误率
	// OutSendPkgErrRate float64 //外网发包错误率
	RecvByte uint64 //接收的字节数
	RecvPkg  uint64 //接收正确的包数
	RecvErr  uint64 //接收错误的包数
	SentByte uint64 //发送的字节数
	SentPkg  uint64 //发送正确的包数
	SentErr  uint64 //发送错误的包数

	RecvByteAvg float64 //一个周期平均每秒接收字节数
	SentByteAvg float64 //一个周期平均每秒发送字节数
	RecvErrRate float64 //一个周期收包错误率
	SentErrRate float64 //一个周期发包错误率
	RecvPkgAvg  float64 //一个周期平均每秒收包数
	SentPkgAvg  float64 //一个周期平均每秒发包数

	Last int64 //上次采集时间
}

type MemoryInfo struct {
	Total uint64 `json:"total"`
	// Used        uint64 `json:"used"`
	UsedPercent uint64 `json:"usedPercent"`
}

func (m MemoryInfo) Collect() MemoryInfo {
	meminfo, err := mem.VirtualMemoryWithContext(context.TODO())
	if err != nil {
		return m
	}
	m.Total = meminfo.Total
	// m.Used = meminfo.Used
	m.UsedPercent = uint64(meminfo.UsedPercent)
	return m
}

type DiskInfo struct {
	Path string `json:"path"`
	// Fstype            string `json:"fstype"`
	// Total uint64 `json:"total"`
	// Free              uint64 `json:"free"`
	// Used        uint64 `json:"used"`
	UsedPercent uint64 `json:"usedPercent"`
	// InodesTotal       uint64 `json:"inodesTotal"`
	// InodesUsed        uint64 `json:"inodesUsed"`
	// InodesFree        uint64 `json:"inodesFree"`
	InodesUsedPercent uint64 `json:"inodesUsedPercent"`
}

func (d DiskInfo) Collect() (diskMapInfo map[string]DiskInfo) {
	diskMapInfo = map[string]DiskInfo{}
	parts, err := disk.PartitionsWithContext(context.TODO(), true)
	if err != nil {
		return
	}
	var diskInfo *disk.UsageStat
	for _, part := range parts {
		diskInfo, err = disk.UsageWithContext(context.TODO(), part.Mountpoint)
		if err != nil {
			return
		}
		if diskInfo.Total == 0 {
			continue
		}
		diskMapInfo[diskInfo.Path] = DiskInfo{
			Path: diskInfo.Path,
			// Fstype:            diskInfo.Fstype,
			// Total:	 		  diskInfo.Total,
			// Free:              diskInfo.Free,
			// Used:        	  diskInfo.Used,
			UsedPercent: uint64(diskInfo.UsedPercent),
			// InodesTotal:       diskInfo.InodesTotal,
			// InodesUsed:        diskInfo.InodesUsed,
			// InodesFree:        diskInfo.InodesFree,
			InodesUsedPercent: uint64(diskInfo.InodesUsedPercent),
		}
	}
	return
}

type GoRuntimeInfo struct {
	// CPU
	NumCpu       uint64 `json:"cpu.count"`      //cpu数
	NumGoroutine uint64 `json:"cpu.goroutines"` //协程数
	NumCgoCall   uint64 `json:"cpu.cgo_calls"`  //cgo调用次数

	// General
	Alloc      uint64 `json:"mem.alloc"`   //Alloc是分配的堆对象的字节
	TotalAlloc uint64 `json:"mem.total"`   //TotalAlloc是分配给堆对象的累积字节
	Sys        uint64 `json:"mem.sys"`     //OS获得的内存的总字节数
	Lookups    uint64 `json:"mem.lookups"` //被runtime监视的指针数
	Mallocs    uint64 `json:"mem.malloc"`  //Mallocs是分配的堆对象的累积计数
	Frees      uint64 `json:"mem.frees"`   //Frees是已释放的堆对象的累积计数

	// Heap
	HeapAlloc    uint64 `json:"mem.heap.alloc"`    //HeapAlloc是分配的堆对象的字节
	HeapSys      uint64 `json:"mem.heap.sys"`      //从操作系统获得的堆内存字节
	HeapIdle     uint64 `json:"mem.heap.idle"`     //HeapIdle减去HeapReleased估计内存量
	HeapInuse    uint64 `json:"mem.heap.inuse"`    //HeapInuse是使用中的跨度中的字节
	HeapReleased uint64 `json:"mem.heap.released"` //HeapReleased是返回操作系统的物理内存字节
	HeapObjects  uint64 `json:"mem.heap.objects"`  //HeapObjects是分配的堆对象的数量

	// Stack
	StackInuse  uint64 `json:"mem.stack.inuse"`        //StackInuse是堆栈跨度中的字节
	StackSys    uint64 `json:"mem.stack.sys"`          //StackSys是从操作系统获得的堆栈内存字节
	MSpanInuse  uint64 `json:"mem.stack.mspan_inuse"`  //MSpanInuse是分配的mspan结构的字节
	MSpanSys    uint64 `json:"mem.stack.mspan_sys"`    //MSpanSys是从操作系统获取的mspan的内存字节
	MCacheInuse uint64 `json:"mem.stack.mcache_inuse"` //MCacheInuse是分配的mcache结构的字节
	MCacheSys   uint64 `json:"mem.stack.mcache_sys"`   //MCacheSys是从OS获得的内存字节

	OtherSys uint64 `json:"mem.othersys"` //OtherSys是其他堆中的内存字节

	// GC
	NumGC uint64 `json:"mem.gc.count"` //NumGC是已完成的GC周期数
}

func (g GoRuntimeInfo) Collect() GoRuntimeInfo {
	goRuntimeInfo := NewGoRuntimeCollector(nil).OneOff()
	g.NumCpu = goRuntimeInfo.NumCpu
	g.NumGoroutine = goRuntimeInfo.NumGoroutine
	g.NumCgoCall = goRuntimeInfo.NumCgoCall
	g.Alloc = goRuntimeInfo.Alloc
	g.TotalAlloc = goRuntimeInfo.TotalAlloc
	g.Sys = goRuntimeInfo.Sys
	g.Lookups = goRuntimeInfo.Lookups
	g.Mallocs = goRuntimeInfo.Mallocs
	g.Frees = goRuntimeInfo.Frees
	g.HeapAlloc = goRuntimeInfo.HeapAlloc
	g.HeapSys = goRuntimeInfo.HeapSys
	g.HeapIdle = goRuntimeInfo.HeapIdle
	g.HeapInuse = goRuntimeInfo.HeapInuse
	g.HeapReleased = goRuntimeInfo.HeapReleased
	g.HeapObjects = goRuntimeInfo.HeapObjects
	g.StackInuse = goRuntimeInfo.StackInuse
	g.StackSys = goRuntimeInfo.StackSys
	g.MSpanInuse = goRuntimeInfo.MSpanInuse
	g.MSpanSys = goRuntimeInfo.MSpanSys
	g.MCacheInuse = goRuntimeInfo.MCacheInuse
	g.MCacheSys = goRuntimeInfo.MCacheSys
	g.OtherSys = goRuntimeInfo.OtherSys
	g.NumGC = goRuntimeInfo.NumGC
	return g
}

type CpuInfo struct {
	UsedPercent uint64 `json:"usedPercent"`
}

func (c CpuInfo) Collect() CpuInfo {
	cpuinfo, err := cpu.PercentWithContext(context.TODO(), time.Second, false)
	if err != nil {
		return c
	}
	if len(cpuinfo) < 1 {
		return c
	}
	c.UsedPercent = uint64(cpuinfo[0])
	return c
}

type NewWorkInfo struct {
	RecvByteAvg uint64 `json:"recvByteAvg"` //平均每秒接收字节数
	RecvPkgAvg  uint64 `json:"recvPkgAvg"`  //平均每秒接收正确的包数
	RecvErrRate uint64 `json:"recvErrRate"` //一个周期收包错误率
	SentByteAvg uint64 `json:"sentByteAvg"` //平均每秒发送字节数
	SentPkgAvg  uint64 `json:"sentPkgAvg"`  //平均每秒发送正确的包数
	SentErrRate uint64 `json:"sentErrRate"` //一个周期发包错误率
}

func (n NewWorkInfo) Collect() (networkMap map[string]NewWorkInfo) {
	networkMap = map[string]NewWorkInfo{}
	netIO, err := utilNet.IOCounters(true)
	if err != nil {
		return
	}
	for _, perIO := range netIO {
		recvByte := perIO.BytesRecv
		recvPkg := perIO.PacketsRecv
		recvErr := perIO.Errin
		sentByte := perIO.BytesSent
		sentPkg := perIO.PacketsSent
		sentErr := perIO.Errout
		//根据网卡名得到对应的网络接口
		ethname := perIO.Name
		netifi, err := net.InterfaceByName(ethname)
		if err != nil {
			continue
		}
		var addrs []net.Addr
		addrs, err = netifi.Addrs()
		if err != nil {
			continue
		}
		if len(addrs) == 0 {
			continue
		}
		moniTag := true
		for _, addr := range addrs {
			cidr := addr.String()
			if strings.Contains(cidr, "0.0.0.0") || strings.Contains(cidr, "127.0.0.1") {
				//0.0.0.0 127.0.0.1 不监控
				moniTag = false
				break
			}
		}
		if moniTag == false {
			continue
		}
		_, exists := IfiMap[ethname]
		if !exists {
			IfiMap[ethname] = &Ifi{}
		}
		ifi, _ := IfiMap[ethname]
		var (
			recvByteAvg float64
			recvPkgAvg  float64
			recvErrRate float64
			sentByteAvg float64
			sentPkgAvg  float64
			sentErrRate float64
		)
		now := time.Now().Unix()
		difftime := float64(now - ifi.Last)
		if ifi.Last == 0 {
			//第一次采集，没有时间差，不计算
		} else {
			if difftime > 0 {
				recvByteAvg = float64(recvByte-ifi.RecvByte) / difftime //平均每秒接收字节数
				recvPkgAvg = float64(recvPkg-ifi.RecvPkg) / difftime    //平均每秒接收正确的包数
				if recvPkg-ifi.RecvPkg > 0 {
					recvErrRate = float64(recvErr-ifi.RecvErr) / float64(recvPkg-ifi.RecvPkg) //一个周期收包错误率
				}
				sentByteAvg = float64(sentByte-ifi.SentByte) / difftime //平均每秒发送字节数
				sentPkgAvg = float64(sentPkg-ifi.SentPkg) / difftime    //平均每秒发送正确的包数
				if sentPkg-ifi.SentPkg > 0 {
					sentErrRate = float64(sentErr-ifi.SentErr) / float64(sentPkg-ifi.SentPkg) //一个周期发包错误率
				}
			}
		}
		ifi.Name = ethname
		// ifi.Ip = addrs[0].String()
		ifi.RecvByte = recvByte
		ifi.RecvPkg = recvPkg
		ifi.RecvErr = recvErr
		ifi.SentByte = sentByte
		ifi.SentPkg = sentPkg
		ifi.SentErr = sentErr
		ifi.RecvByteAvg = recvByteAvg
		ifi.SentByteAvg = sentByteAvg
		ifi.RecvErrRate = recvErrRate
		ifi.SentErrRate = sentErrRate
		ifi.RecvPkgAvg = recvPkgAvg
		ifi.SentPkgAvg = sentPkgAvg
		ifi.Last = now
		networkMap[ethname] = NewWorkInfo{
			RecvByteAvg: uint64(recvByteAvg),
			RecvPkgAvg:  uint64(recvPkgAvg),
			RecvErrRate: uint64(recvErrRate * 100),
			SentByteAvg: uint64(sentByteAvg),
			SentPkgAvg:  uint64(sentPkgAvg),
			SentErrRate: uint64(sentErrRate * 100),
		}
	}
	return
}

type IOStatLog struct {
	ReadByte   uint64
	ReadCount  uint64
	WriteByte  uint64
	WriteCount uint64

	Last int64 //上次采集时间
}

type IOInfo struct {
	ReadByte   uint64
	ReadCount  uint64
	WriteByte  uint64
	WriteCount uint64
}

func (i IOInfo) Collect() IOInfo {
	iostat, err := disk.IOCountersWithContext(context.Background())
	if err != nil {
		return IOInfo{}
	}
	var (
		allReadByte   uint64 = 0
		allReadCount  uint64 = 0
		allWriteByte  uint64 = 0
		allWriteCount uint64 = 0
	)
	for diskName, v := range iostat {
		if runtime.GOOS == "linux" {
			reg := regexp.MustCompile(`\d+`)
			if reg.MatchString(diskName) {
				continue
			}
		}
		allReadByte += v.ReadBytes
		allReadCount += v.ReadCount
		allWriteByte += v.WriteBytes
		allWriteCount += v.WriteCount
	}
	var (
		readByteAvg   float64
		readCountAvg  float64
		writeByteAvg  float64
		writeCountAvg float64
	)
	now := time.Now().Unix()
	difftime := float64(now - IOStat.Last)
	if IOStat.Last == 0 {
		//第一次采集，没有时间差，不计算
	} else {
		if difftime > 0 {
			readByteAvg = float64(allReadByte-IOStat.ReadByte) / difftime
			readCountAvg = float64(allReadCount-IOStat.ReadCount) / difftime
			writeByteAvg = float64(allWriteByte-IOStat.WriteByte) / difftime
			writeCountAvg = float64(allWriteCount-IOStat.WriteCount) / difftime
		}
	}
	IOStat.Last = now
	IOStat.ReadByte = allReadByte
	IOStat.ReadCount = allReadCount
	IOStat.WriteByte = allWriteByte
	IOStat.WriteCount = allWriteCount
	return IOInfo{
		ReadByte:   uint64(readByteAvg),
		ReadCount:  uint64(readCountAvg),
		WriteByte:  uint64(writeByteAvg),
		WriteCount: uint64(writeCountAvg),
	}
}
