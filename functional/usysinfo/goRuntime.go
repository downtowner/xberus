package usysinfo

import (
	"runtime"
	"time"
)

// FieldsFunc represents a callback after successfully gathering statistics
type GoRuntimeFieldsFunc func(GoRuntimeFields)

// Collector implements the periodic grabbing of informational data from the
// runtime package and outputting the values to a GaugeFunc.
type GoRuntimeCollector struct {
	// PauseDur represents the interval in-between each set of stats output.
	// Defaults to 10 seconds.
	PauseDur time.Duration

	// EnableCPU determines whether CPU statistics will be output. Defaults to true.
	EnableCPU bool

	// EnableMem determines whether memory statistics will be output. Defaults to true.
	EnableMem bool

	// EnableGC determines whether garbage collection statistics will be output. EnableMem
	// must also be set to true for this to take affect. Defaults to true.
	EnableGC bool

	// Done, when closed, is used to signal Collector that is should stop collecting
	// statistics and the Run function should return.
	Done <-chan struct{}

	fieldsFunc GoRuntimeFieldsFunc
}

// New creates a new Collector that will periodically output statistics to fieldsFunc. It
// will also set the values of the exported fields to the described defaults. The values
// of the exported defaults can be changed at any point before Run is called.
func NewGoRuntimeCollector(fieldsFunc GoRuntimeFieldsFunc) *GoRuntimeCollector {
	if fieldsFunc == nil {
		fieldsFunc = func(GoRuntimeFields) {}
	}

	return &GoRuntimeCollector{
		PauseDur:   10 * time.Second,
		EnableCPU:  true,
		EnableMem:  true,
		EnableGC:   true,
		fieldsFunc: fieldsFunc,
	}
}

// Run gathers statistics then outputs them to the configured PointFunc every
// PauseDur. Unlike OneOff, this function will return until Done has been closed
// (or never if Done is nil), therefore it should be called in its own go routine.
func (g *GoRuntimeCollector) Run() {
	g.fieldsFunc(g.collectStats())

	tick := time.NewTicker(g.PauseDur)
	defer tick.Stop()
	for {
		select {
		case <-g.Done:
			return
		case <-tick.C:
			g.fieldsFunc(g.collectStats())
		}
	}
}

// OneOff gathers returns a map containing all statistics. It is safe for use from
// multiple go routines
func (g *GoRuntimeCollector) OneOff() GoRuntimeFields {
	return g.collectStats()
}

func (g *GoRuntimeCollector) collectStats() GoRuntimeFields {
	fields := GoRuntimeFields{}

	if g.EnableCPU {
		cStats := cpuStats{
			NumGoroutine: uint64(runtime.NumGoroutine()),
			NumCgoCall:   uint64(runtime.NumCgoCall()),
			NumCpu:       uint64(runtime.NumCPU()),
		}
		g.collectCPUStats(&fields, &cStats)
	}
	if g.EnableMem {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)
		g.collectMemStats(&fields, m)
		if g.EnableGC {
			g.collectGCStats(&fields, m)
		}
	}

	fields.Goos = runtime.GOOS
	fields.Goarch = runtime.GOARCH
	fields.Version = runtime.Version()

	return fields
}

func (g *GoRuntimeCollector) collectCPUStats(fields *GoRuntimeFields, s *cpuStats) {
	fields.NumCpu = s.NumCpu
	fields.NumGoroutine = s.NumGoroutine
	fields.NumCgoCall = s.NumCgoCall
}

func (g *GoRuntimeCollector) collectMemStats(fields *GoRuntimeFields, m *runtime.MemStats) {
	// General
	fields.Alloc = m.Alloc
	fields.TotalAlloc = m.TotalAlloc
	fields.Sys = m.Sys
	fields.Lookups = m.Lookups
	fields.Mallocs = m.Mallocs
	fields.Frees = m.Frees

	// Heap
	fields.HeapAlloc = m.HeapAlloc
	fields.HeapSys = m.HeapSys
	fields.HeapIdle = m.HeapIdle
	fields.HeapInuse = m.HeapInuse
	fields.HeapReleased = m.HeapReleased
	fields.HeapObjects = m.HeapObjects

	// Stack
	fields.StackInuse = m.StackInuse
	fields.StackSys = m.StackSys
	fields.MSpanInuse = m.MSpanInuse
	fields.MSpanSys = m.MSpanSys
	fields.MCacheInuse = m.MCacheInuse
	fields.MCacheSys = m.MCacheSys

	fields.OtherSys = m.OtherSys
}

func (g *GoRuntimeCollector) collectGCStats(fields *GoRuntimeFields, m *runtime.MemStats) {
	fields.GCSys = m.GCSys
	fields.NextGC = m.NextGC
	fields.LastGC = m.LastGC
	fields.PauseTotalNs = m.PauseTotalNs
	fields.PauseNs = m.PauseNs[(m.NumGC+255)%256]
	fields.NumGC = uint64(m.NumGC)
	fields.GCCPUFraction = m.GCCPUFraction
}

type cpuStats struct {
	NumCpu       uint64
	NumGoroutine uint64
	NumCgoCall   uint64
}

type GoRuntimeFields struct {
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
	GCSys         uint64  `json:"mem.gc.sys"`          //GCSys是垃圾回收元数据中的内存字节
	NextGC        uint64  `json:"mem.gc.next"`         //NextGC是下一个GC周期的目标堆大小
	LastGC        uint64  `json:"mem.gc.last"`         //LastGC是最后一个垃圾回收完成的时间
	PauseTotalNs  uint64  `json:"mem.gc.pause_total"`  //PauseTotalNs是GC中的累积纳秒
	PauseNs       uint64  `json:"mem.gc.pause"`        //PauseNs是最近的GC停止发行的循环缓冲区
	NumGC         uint64  `json:"mem.gc.count"`        //NumGC是已完成的GC周期数
	GCCPUFraction float64 `json:"mem.gc.cpu_fraction"` //自程序启动以来GC使用的CPU时间

	Goarch  string `json:"cpu_arch"`
	Goos    string `json:"os"`
	Version string `json:"go_version"`
}

func (g *GoRuntimeFields) Tags() map[string]string {
	return map[string]string{
		"go.os":      g.Goos,
		"go.arch":    g.Goarch,
		"go.version": g.Version,
	}
}

func (g *GoRuntimeFields) Values() map[string]interface{} {
	return map[string]interface{}{
		"cpu.count":      g.NumCpu,
		"cpu.goroutines": g.NumGoroutine,
		"cpu.cgo_calls":  g.NumCgoCall,

		"mem.alloc":   g.Alloc,
		"mem.total":   g.TotalAlloc,
		"mem.sys":     g.Sys,
		"mem.lookups": g.Lookups,
		"mem.malloc":  g.Mallocs,
		"mem.frees":   g.Frees,

		"mem.heap.alloc":    g.HeapAlloc,
		"mem.heap.sys":      g.HeapSys,
		"mem.heap.idle":     g.HeapIdle,
		"mem.heap.inuse":    g.HeapInuse,
		"mem.heap.released": g.HeapReleased,
		"mem.heap.objects":  g.HeapObjects,

		"mem.stack.inuse":        g.StackInuse,
		"mem.stack.sys":          g.StackSys,
		"mem.stack.mspan_inuse":  g.MSpanInuse,
		"mem.stack.mspan_sys":    g.MSpanSys,
		"mem.stack.mcache_inuse": g.MCacheInuse,
		"mem.stack.mcache_sys":   g.MCacheSys,
		"mem.othersys":           g.OtherSys,

		"mem.gc.sys":          g.GCSys,
		"mem.gc.next":         g.NextGC,
		"mem.gc.last":         g.LastGC,
		"mem.gc.pause_total":  g.PauseTotalNs,
		"mem.gc.pause":        g.PauseNs,
		"mem.gc.count":        g.NumGC,
		"mem.gc.cpu_fraction": float64(g.GCCPUFraction),
	}
}
