package usysinfo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	//db文件名
	DBPath = "sysinfo.db"
	//bucket前缀
	DBPrefix = "log_systeminfo_"

	cpuUsedPercent = "cpu.UsedPercent"

	memoryUsedPercent = "memory.UsedPercent"

	diskUsedPercent       = "disk.UsedPercent"
	diskInodesUsedPercent = "disk.InodesUsedPercent"

	networkRecvByteAvg = "network.RecvByteAvg"
	networkRecvPkgAvg  = "network.RecvPkgAvg"
	networkRecvErrRate = "network.RecvErrRate"
	networkSentByteAvg = "network.SentByteAvg"
	networkSentPkgAvg  = "network.SentPkgAvg"
	networkSentErrRate = "network.SentErrRate"

	goNumCpu       = "go.NumCpu"
	goNumGoroutine = "go.NumGoroutine"
	goNumCgoCall   = "go.NumCgoCall"
	goAlloc        = "go.Alloc"
	goTotalAlloc   = "go.TotalAlloc"
	goSys          = "go.Sys"
	goLookups      = "go.Lookups"
	goMallocs      = "go.Mallocs"
	goFrees        = "go.Frees"
	goHeapAlloc    = "go.HeapAlloc"
	goHeapSys      = "go.HeapSys"
	goHeapIdle     = "go.HeapIdle"
	goHeapInuse    = "go.HeapInuse"
	goHeapReleased = "go.HeapReleased"
	goHeapObjects  = "go.HeapObjects"
	goStackInuse   = "go.StackInuse"
	goStackSys     = "go.StackSys"
	goMSpanInuse   = "go.MSpanInuse"
	goMSpanSys     = "go.MSpanSys"
	goMCacheInuse  = "go.MCacheInuse"
	goMCacheSys    = "go.MCacheSys"
	goOtherSys     = "go.OtherSys"
	goNumGC        = "go.NumGC"

	ioReadCount  = "io.ReadCount"
	ioWriteCount = "io.WriteCount"
	ioReadByte   = "io.ReadByte"
	ioWriteByte  = "io.WriteByte"
)

// 存储的keyname
var storageKeyName = map[string]string{
	//cpu
	cpuUsedPercent: "cpu.UsedPercent",
	//memory
	memoryUsedPercent: "memory.UsedPercent",
	//disk
	diskUsedPercent:       "disk.UsedPercent.%s",
	diskInodesUsedPercent: "disk.InodesUsedPercent.%s", //windows无inodes
	//network
	networkRecvByteAvg: "network.RecvByteAvg.%s", //平均每秒接收字节数
	networkRecvPkgAvg:  "network.RecvPkgAvg.%s",  //平均每秒接收正确的包数
	networkRecvErrRate: "network.RecvErrRate.%s", //一个周期收包错误率
	networkSentByteAvg: "network.SentByteAvg.%s", //平均每秒发送字节数
	networkSentPkgAvg:  "network.SentPkgAvg.%s",  //平均每秒发送正确的包数
	networkSentErrRate: "network.SentErrRate.%s", //一个周期发包错误率
	//goruntime
	goNumCpu:       "go.NumCpu",
	goNumGoroutine: "go.NumGoroutine",
	goNumCgoCall:   "go.NumCgoCall",
	goAlloc:        "go.Alloc",
	goTotalAlloc:   "go.TotalAlloc",
	goSys:          "go.Sys",
	goLookups:      "go.Lookups",
	goMallocs:      "go.Mallocs",
	goFrees:        "go.Frees",
	goHeapAlloc:    "go.HeapAlloc",
	goHeapSys:      "go.HeapSys",
	goHeapIdle:     "go.HeapIdle",
	goHeapInuse:    "go.HeapInuse",
	goHeapReleased: "go.HeapReleased",
	goHeapObjects:  "go.HeapObjects",
	goStackInuse:   "go.StackInuse",
	goStackSys:     "go.StackSys",
	goMSpanInuse:   "go.MSpanInuse",
	goMSpanSys:     "go.MSpanSys",
	goMCacheInuse:  "go.MCacheInuse",
	goMCacheSys:    "go.MCacheSys",
	goOtherSys:     "go.OtherSys",
	goNumGC:        "go.NumGC",
	//io
	ioReadCount:  "io.ReadCount",
	ioWriteCount: "io.WriteCount",
	ioReadByte:   "io.ReadByte",
	ioWriteByte:  "io.WriteByte",
}

//采集数据存储格式
type OneCollectData struct {
	Metrics uint64 `json:"m"` //指标
	Time    int64  `json:"t"` //采集时间
}

type StorageEngine interface {
	//初始化
	InitEngine() error
	//定时轮询写入 cpu使用率 内存使用量 系统平均负载 磁盘使用量 网络流入流出 tcp连接数
	Set(StorageInfo) error
	//获取某一天的数据
	GetCpu(startdate, enddate time.Time) CpuData
	GetMemory(startdate, enddate time.Time) MemoryData
	GetGoRuntime(startdate, enddate time.Time) GoRuntimeData
	GetDisk(startdate, enddate time.Time) map[string]*DiskData
	GetNetwork(startdate, enddate time.Time) map[string]*NetworkData
	GetIO(startdate, enddate time.Time) IOData
	//删除某一天的数据
	Delete(day string)
}

type CpuData struct {
	UsedPercent []OneCollectData
}

func (c CpuData) GetCpu() CpuData {
	return c
}

type MemoryData struct {
	UsedPercent []OneCollectData
}

type DiskData struct {
	UsedPercent       []OneCollectData
	InodesUsedPercent []OneCollectData
}

type NetworkData struct {
	RecvByteAvg []OneCollectData
	RecvPkgAvg  []OneCollectData
	RecvErrRate []OneCollectData
	SentByteAvg []OneCollectData
	SentPkgAvg  []OneCollectData
	SentErrRate []OneCollectData
}

type GoRuntimeData struct {
	NumCpu       []OneCollectData `json:"goNumCpu"`
	NumGoroutine []OneCollectData `json:"goNumGoroutine"`
	NumCgoCall   []OneCollectData `json:"goNumCgoCall"`
	Alloc        []OneCollectData `json:"goAlloc"`
	TotalAlloc   []OneCollectData `json:"goTotalAlloc"`
	Sys          []OneCollectData `json:"goSys"`
	Lookups      []OneCollectData `json:"goLookups"`
	Mallocs      []OneCollectData `json:"goMallocs"`
	Frees        []OneCollectData `json:"goFrees"`
	HeapAlloc    []OneCollectData `json:"goHeapAlloc"`
	HeapSys      []OneCollectData `json:"goHeapSys"`
	HeapIdle     []OneCollectData `json:"goHeapIdle"`
	HeapInuse    []OneCollectData `json:"goHeapInuse"`
	HeapReleased []OneCollectData `json:"goHeapReleased"`
	HeapObjects  []OneCollectData `json:"goHeapObjects"`
	StackInuse   []OneCollectData `json:"goStackInuse"`
	StackSys     []OneCollectData `json:"goStackSys"`
	MSpanInuse   []OneCollectData `json:"goMSpanInuse"`
	MSpanSys     []OneCollectData `json:"goMSpanSys"`
	MCacheInuse  []OneCollectData `json:"goMCacheInuse"`
	MCacheSys    []OneCollectData `json:"goMCacheSys"`
	OtherSys     []OneCollectData `json:"goOtherSys"`
	NumGC        []OneCollectData `json:"goNumGC"`
}

type IOData struct {
	ReadByte   []OneCollectData `json:"readByte"`   //每秒读字节数
	ReadCount  []OneCollectData `json:"readCount"`  //每秒读次数
	WriteByte  []OneCollectData `json:"writeByte"`  //每秒写字节数
	WriteCount []OneCollectData `json:"writeCount"` //每秒写次数
}

type Storage struct {
	once sync.Once
	db   *bolt.DB
	StorageEngine
}

type StorageInfo struct {
	Time      int64                  `json:"time"`
	Hostname  string                 `json:"hostname"`
	Memory    MemoryInfo             `json:"memory"`
	Disk      map[string]DiskInfo    `json:"disk"`
	Cpu       CpuInfo                `json:"cpu"`
	NetWork   map[string]NewWorkInfo `json:"network"`
	GoRuntime GoRuntimeInfo          `json:"goRuntime"`
	IO        IOInfo                 `json:"io"`
}

func (s *sysinfoObj) NewStorage() StorageEngine {
	return &Storage{}
}

func (s *Storage) InitEngine() error {
	var err error
	s.once.Do(func() {
		s.db, err = bolt.Open(DBPath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	})
	return err
}

func (s *Storage) Set(info StorageInfo) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucketName := s.bucketName(time.Unix(info.Time, 0).Format("20060102"))
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		err = s.setCpu(bucket, info.Cpu, info.Time)
		if err != nil {
			return err
		}
		err = s.setMemory(bucket, info.Memory, info.Time)
		if err != nil {
			return err
		}
		err = s.setGoRuntime(bucket, info.GoRuntime, info.Time)
		if err != nil {
			return err
		}
		err = s.setNetwork(bucket, info.NetWork, info.Time)
		if err != nil {
			return err
		}
		err = s.setDisk(bucket, info.Disk, info.Time)
		if err != nil {
			return err
		}
		err = s.setIO(bucket, info.IO, info.Time)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (s *Storage) GetCpu(startDatetime, endDatetime time.Time) CpuData {
	cpuData := CpuData{}
	tempDatetime := startDatetime
	s.db.View(func(tx *bolt.Tx) error {
		for {
			day := FormatTime(tempDatetime, "20060102")
			bucketName := s.bucketName(day)
			bucket := tx.Bucket(bucketName)
			if bucket != nil {
				cursor := bucket.Cursor()
				prefix := []byte("cpu.")
				for key, jsonValue := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, jsonValue = cursor.Next() {
					var collectData []OneCollectData
					err := json.Unmarshal(jsonValue, &collectData)
					if err != nil {
						continue
					}
					for _, one := range collectData {
						if startDatetime.Unix() <= one.Time && one.Time <= endDatetime.Unix() {
							switch string(key) {
							case cpuUsedPercent:
								cpuData.UsedPercent = append(cpuData.UsedPercent, one)
							}
						}
					}
				}
			}
			mustBreak := false
			mustBreak, tempDatetime = datetimeOffset(tempDatetime, endDatetime)
			if mustBreak {
				break
			}
		}
		return nil
	})
	return cpuData
}

func (s *Storage) GetMemory(startDatetime, endDatetime time.Time) MemoryData {
	memoryData := MemoryData{}
	tempDatetime := startDatetime
	s.db.View(func(tx *bolt.Tx) error {
		for {
			day := FormatTime(tempDatetime, "20060102")
			bucketName := s.bucketName(day)
			bucket := tx.Bucket(bucketName)
			if bucket != nil {
				cursor := bucket.Cursor()
				prefix := []byte("memory.")
				for key, jsonValue := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, jsonValue = cursor.Next() {
					var collectData []OneCollectData
					err := json.Unmarshal(jsonValue, &collectData)
					if err != nil {
						continue
					}
					for _, one := range collectData {
						if startDatetime.Unix() <= one.Time && one.Time <= endDatetime.Unix() {
							switch string(key) {
							case memoryUsedPercent:
								memoryData.UsedPercent = append(memoryData.UsedPercent, one)
							}
						}
					}
				}
			}
			mustBreak := false
			mustBreak, tempDatetime = datetimeOffset(tempDatetime, endDatetime)
			if mustBreak {
				break
			}
		}
		return nil
	})
	return memoryData
}

func (s *Storage) GetGoRuntime(startDatetime, endDatetime time.Time) GoRuntimeData {
	goRuntimeData := GoRuntimeData{}
	tempDatetime := startDatetime
	s.db.View(func(tx *bolt.Tx) error {
		for {
			day := FormatTime(tempDatetime, "20060102")
			bucketName := s.bucketName(day)
			bucket := tx.Bucket(bucketName)
			if bucket != nil {
				cursor := bucket.Cursor()
				prefix := []byte("go.")
				for key, jsonValue := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, jsonValue = cursor.Next() {
					var collectData []OneCollectData
					err := json.Unmarshal(jsonValue, &collectData)
					if err != nil {
						continue
					}
					for _, one := range collectData {
						if startDatetime.Unix() <= one.Time && one.Time <= endDatetime.Unix() {
							switch string(key) {
							case goNumCpu:
								goRuntimeData.NumCpu = append(goRuntimeData.NumCpu, one)
							case goNumGoroutine:
								goRuntimeData.NumGoroutine = append(goRuntimeData.NumGoroutine, one)
							case goNumCgoCall:
								goRuntimeData.NumCgoCall = append(goRuntimeData.NumCgoCall, one)
							case goAlloc:
								goRuntimeData.Alloc = append(goRuntimeData.Alloc, one)
							case goTotalAlloc:
								goRuntimeData.TotalAlloc = append(goRuntimeData.TotalAlloc, one)
							case goSys:
								goRuntimeData.Sys = append(goRuntimeData.Sys, one)
							case goLookups:
								goRuntimeData.Lookups = append(goRuntimeData.Lookups, one)
							case goMallocs:
								goRuntimeData.Mallocs = append(goRuntimeData.Mallocs, one)
							case goFrees:
								goRuntimeData.Frees = append(goRuntimeData.Frees, one)
							case goHeapAlloc:
								goRuntimeData.HeapAlloc = append(goRuntimeData.HeapAlloc, one)
							case goHeapSys:
								goRuntimeData.HeapSys = append(goRuntimeData.HeapSys, one)
							case goHeapIdle:
								goRuntimeData.HeapIdle = append(goRuntimeData.HeapIdle, one)
							case goHeapInuse:
								goRuntimeData.HeapInuse = append(goRuntimeData.HeapInuse, one)
							case goHeapReleased:
								goRuntimeData.HeapReleased = append(goRuntimeData.HeapReleased, one)
							case goHeapObjects:
								goRuntimeData.HeapObjects = append(goRuntimeData.HeapObjects, one)
							case goStackInuse:
								goRuntimeData.StackInuse = append(goRuntimeData.StackInuse, one)
							case goStackSys:
								goRuntimeData.StackSys = append(goRuntimeData.StackSys, one)
							case goMSpanInuse:
								goRuntimeData.MSpanInuse = append(goRuntimeData.MSpanInuse, one)
							case goMSpanSys:
								goRuntimeData.MSpanSys = append(goRuntimeData.MSpanSys, one)
							case goMCacheInuse:
								goRuntimeData.MCacheInuse = append(goRuntimeData.MCacheInuse, one)
							case goMCacheSys:
								goRuntimeData.MCacheSys = append(goRuntimeData.MCacheSys, one)
							case goOtherSys:
								goRuntimeData.OtherSys = append(goRuntimeData.OtherSys, one)
							case goNumGC:
								goRuntimeData.NumGC = append(goRuntimeData.NumGC, one)
							}
						}
					}
				}
			}
			mustBreak := false
			mustBreak, tempDatetime = datetimeOffset(tempDatetime, endDatetime)
			if mustBreak {
				break
			}
		}
		return nil
	})
	return goRuntimeData
}

func (s *Storage) GetIO(startDatetime, endDatetime time.Time) IOData {
	ioData := IOData{}
	tempDatetime := startDatetime
	s.db.View(func(tx *bolt.Tx) error {
		for {
			day := FormatTime(tempDatetime, "20060102")
			bucketName := s.bucketName(day)
			bucket := tx.Bucket(bucketName)
			if bucket != nil {
				cursor := bucket.Cursor()
				prefix := []byte("io.")
				for key, jsonValue := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, jsonValue = cursor.Next() {
					var collectData []OneCollectData
					err := json.Unmarshal(jsonValue, &collectData)
					if err != nil {
						continue
					}
					for _, one := range collectData {
						if startDatetime.Unix() <= one.Time && one.Time <= endDatetime.Unix() {
							switch string(key) {
							case ioReadByte:
								ioData.ReadByte = append(ioData.ReadByte, one)
							case ioReadCount:
								ioData.ReadCount = append(ioData.ReadCount, one)
							case ioWriteByte:
								ioData.WriteByte = append(ioData.WriteByte, one)
							case ioWriteCount:
								ioData.WriteCount = append(ioData.WriteCount, one)
							}
						}
					}
				}
			}
			mustBreak := false
			mustBreak, tempDatetime = datetimeOffset(tempDatetime, endDatetime)
			if mustBreak {
				break
			}
		}
		return nil
	})
	return ioData
}

func (s *Storage) GetDisk(startDatetime, endDatetime time.Time) map[string]*DiskData {
	diskData := map[string]*DiskData{}
	tempDatetime := startDatetime
	s.db.View(func(tx *bolt.Tx) error {
		for {
			day := FormatTime(tempDatetime, "20060102")
			bucketName := s.bucketName(day)
			bucket := tx.Bucket(bucketName)
			if bucket != nil {
				cursor := bucket.Cursor()
				prefix := []byte("disk.")
				for key, jsonValue := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, jsonValue = cursor.Next() {
					var collectData []OneCollectData
					err := json.Unmarshal(jsonValue, &collectData)
					if err != nil {
						continue
					}
					for _, one := range collectData {
						if startDatetime.Unix() <= one.Time && one.Time <= endDatetime.Unix() {
							strKey := string(key)
							strSplit := strings.Split(strKey, ".")
							if len(strSplit) != 3 {
								continue
							}
							metrics := strings.ReplaceAll(strKey, ".%s", "")
							if _, ok := diskData[strSplit[2]]; !ok {
								diskData[strSplit[2]] = &DiskData{}
							}
							if strings.Contains(metrics, diskInodesUsedPercent) {
								diskData[strSplit[2]].InodesUsedPercent = append(diskData[strSplit[2]].InodesUsedPercent, one)
								continue
							}
							if strings.Contains(metrics, diskUsedPercent) {
								diskData[strSplit[2]].UsedPercent = append(diskData[strSplit[2]].UsedPercent, one)
								continue
							}
						}
					}
				}
			}
			mustBreak := false
			mustBreak, tempDatetime = datetimeOffset(tempDatetime, endDatetime)
			if mustBreak {
				break
			}
		}
		return nil
	})
	return diskData
}

func (s *Storage) GetNetwork(startDatetime, endDatetime time.Time) map[string]*NetworkData {
	networkData := map[string]*NetworkData{}
	tempDatetime := startDatetime
	s.db.View(func(tx *bolt.Tx) error {
		for {
			day := FormatTime(tempDatetime, "20060102")
			bucketName := s.bucketName(day)
			bucket := tx.Bucket(bucketName)
			if bucket != nil {
				cursor := bucket.Cursor()
				prefix := []byte("network.")
				for key, jsonValue := cursor.Seek(prefix); key != nil && bytes.HasPrefix(key, prefix); key, jsonValue = cursor.Next() {
					var collectData []OneCollectData
					err := json.Unmarshal(jsonValue, &collectData)
					if err != nil {
						continue
					}
					for _, one := range collectData {
						if startDatetime.Unix() <= one.Time && one.Time <= endDatetime.Unix() {
							strKey := string(key)
							strSplit := strings.Split(strKey, ".")
							if len(strSplit) != 3 {
								continue
							}
							metrics := strings.ReplaceAll(strKey, ".%s", "")
							if _, ok := networkData[strSplit[2]]; !ok {
								networkData[strSplit[2]] = &NetworkData{}
							}
							if strings.Contains(metrics, networkRecvByteAvg) {
								networkData[strSplit[2]].RecvByteAvg = append(networkData[strSplit[2]].RecvByteAvg, one)
								continue
							}
							if strings.Contains(metrics, networkRecvPkgAvg) {
								networkData[strSplit[2]].RecvPkgAvg = append(networkData[strSplit[2]].RecvPkgAvg, one)
								continue
							}
							if strings.Contains(metrics, networkRecvErrRate) {
								networkData[strSplit[2]].RecvErrRate = append(networkData[strSplit[2]].RecvErrRate, one)
								continue
							}
							if strings.Contains(metrics, networkSentByteAvg) {
								networkData[strSplit[2]].SentByteAvg = append(networkData[strSplit[2]].SentByteAvg, one)
								continue
							}
							if strings.Contains(metrics, networkSentPkgAvg) {
								networkData[strSplit[2]].SentPkgAvg = append(networkData[strSplit[2]].SentPkgAvg, one)
								continue
							}
							if strings.Contains(metrics, networkSentErrRate) {
								networkData[strSplit[2]].SentErrRate = append(networkData[strSplit[2]].SentErrRate, one)
								continue
							}
						}
					}
				}
			}
			mustBreak := false
			mustBreak, tempDatetime = datetimeOffset(tempDatetime, endDatetime)
			if mustBreak {
				break
			}
		}
		return nil
	})
	return networkData
}

func (s *Storage) Delete(day string) {
	s.db.View(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(s.bucketName(day))
	})
}

func (s *Storage) bucketName(day string) []byte {
	return []byte(DBPrefix + day)
}

func (s *Storage) setCpu(bucket *bolt.Bucket, cpu CpuInfo, timeUnix int64) error {
	saveToDB(bucket, storageKeyName[cpuUsedPercent], cpu.UsedPercent, timeUnix)
	return nil
}

func (s *Storage) setMemory(bucket *bolt.Bucket, mem MemoryInfo, timeUnix int64) error {
	saveToDB(bucket, storageKeyName[memoryUsedPercent], mem.UsedPercent, timeUnix)
	return nil
}

func (s *Storage) setGoRuntime(bucket *bolt.Bucket, goruntime GoRuntimeInfo, timeUnix int64) error {
	saveToDB(bucket, storageKeyName[goNumCpu], goruntime.NumCpu, timeUnix)
	saveToDB(bucket, storageKeyName[goNumGoroutine], goruntime.NumGoroutine, timeUnix)
	saveToDB(bucket, storageKeyName[goNumCgoCall], goruntime.NumCgoCall, timeUnix)
	saveToDB(bucket, storageKeyName[goAlloc], goruntime.Alloc, timeUnix)
	saveToDB(bucket, storageKeyName[goTotalAlloc], goruntime.TotalAlloc, timeUnix)
	saveToDB(bucket, storageKeyName[goSys], goruntime.Sys, timeUnix)
	saveToDB(bucket, storageKeyName[goLookups], goruntime.Lookups, timeUnix)
	saveToDB(bucket, storageKeyName[goMallocs], goruntime.Mallocs, timeUnix)
	saveToDB(bucket, storageKeyName[goFrees], goruntime.Frees, timeUnix)
	saveToDB(bucket, storageKeyName[goHeapAlloc], goruntime.HeapAlloc, timeUnix)
	saveToDB(bucket, storageKeyName[goHeapSys], goruntime.HeapSys, timeUnix)
	saveToDB(bucket, storageKeyName[goHeapIdle], goruntime.HeapIdle, timeUnix)
	saveToDB(bucket, storageKeyName[goHeapInuse], goruntime.HeapInuse, timeUnix)
	saveToDB(bucket, storageKeyName[goHeapReleased], goruntime.HeapReleased, timeUnix)
	saveToDB(bucket, storageKeyName[goHeapObjects], goruntime.HeapObjects, timeUnix)
	saveToDB(bucket, storageKeyName[goStackInuse], goruntime.StackInuse, timeUnix)
	saveToDB(bucket, storageKeyName[goStackSys], goruntime.StackSys, timeUnix)
	saveToDB(bucket, storageKeyName[goMSpanInuse], goruntime.MSpanInuse, timeUnix)
	saveToDB(bucket, storageKeyName[goMSpanSys], goruntime.MSpanSys, timeUnix)
	saveToDB(bucket, storageKeyName[goMCacheInuse], goruntime.MCacheInuse, timeUnix)
	saveToDB(bucket, storageKeyName[goMCacheSys], goruntime.MCacheSys, timeUnix)
	saveToDB(bucket, storageKeyName[goOtherSys], goruntime.OtherSys, timeUnix)
	saveToDB(bucket, storageKeyName[goNumGC], goruntime.NumGC, timeUnix)

	return nil
}

func (s *Storage) setIO(bucket *bolt.Bucket, ioInfo IOInfo, timeUnix int64) error {
	saveToDB(bucket, storageKeyName[ioReadByte], ioInfo.ReadByte, timeUnix)
	saveToDB(bucket, storageKeyName[ioReadCount], ioInfo.ReadCount, timeUnix)
	saveToDB(bucket, storageKeyName[ioWriteByte], ioInfo.WriteByte, timeUnix)
	saveToDB(bucket, storageKeyName[ioWriteCount], ioInfo.WriteCount, timeUnix)
	return nil
}

func (s *Storage) setNetwork(bucket *bolt.Bucket, neworkInfo map[string]NewWorkInfo, timeUnix int64) error {
	for ifiName, _ := range neworkInfo {
		keyName := fmt.Sprintf(storageKeyName[networkRecvByteAvg], ifiName)
		saveToDB(bucket, keyName, neworkInfo[ifiName].RecvByteAvg, timeUnix)
		keyName = fmt.Sprintf(storageKeyName[networkRecvPkgAvg], ifiName)
		saveToDB(bucket, keyName, neworkInfo[ifiName].RecvPkgAvg, timeUnix)
		keyName = fmt.Sprintf(storageKeyName[networkRecvErrRate], ifiName)
		saveToDB(bucket, keyName, neworkInfo[ifiName].RecvErrRate, timeUnix)
		keyName = fmt.Sprintf(storageKeyName[networkSentByteAvg], ifiName)
		saveToDB(bucket, keyName, neworkInfo[ifiName].SentByteAvg, timeUnix)
		keyName = fmt.Sprintf(storageKeyName[networkSentPkgAvg], ifiName)
		saveToDB(bucket, keyName, neworkInfo[ifiName].SentPkgAvg, timeUnix)
		keyName = fmt.Sprintf(storageKeyName[networkSentErrRate], ifiName)
		saveToDB(bucket, keyName, neworkInfo[ifiName].SentErrRate, timeUnix)
	}
	return nil
}

func (s *Storage) setDisk(bucket *bolt.Bucket, diskInfo map[string]DiskInfo, timeUnix int64) error {
	for diskName, _ := range diskInfo {
		keyName := fmt.Sprintf(storageKeyName[diskUsedPercent], diskName)
		saveToDB(bucket, keyName, diskInfo[diskName].UsedPercent, timeUnix)
		keyName = fmt.Sprintf(storageKeyName[diskInodesUsedPercent], diskName)
		saveToDB(bucket, keyName, diskInfo[diskName].InodesUsedPercent, timeUnix)
	}
	return nil
}

func GetDBSize() int64 {
	var result int64
	filepath.Walk(DBPath, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func saveToDB(bucket *bolt.Bucket, key string, value uint64, detailTime int64) error {
	dbKey := []byte(key)
	oldJson := bucket.Get(dbKey)
	var oldData []OneCollectData
	err := json.Unmarshal(oldJson, &oldData)
	if err != nil {
		oldData = []OneCollectData{}
	}
	oldData = append(oldData, OneCollectData{
		Metrics: value,
		Time:    detailTime,
	})
	newJson, err := json.Marshal(oldData)
	if err != nil {
		return err
	}
	bucket.Put(dbKey, newJson)
	return err
}

func datetimeOffset(tempDatetime, endDatetime time.Time) (bool, time.Time) {
	stopOffset := false
	subTime := endDatetime.Unix() - tempDatetime.Unix()
	hour := subTime / 3600
	if tempDatetime.After(endDatetime) || tempDatetime.Equal(endDatetime) {
		stopOffset = true
		return stopOffset, tempDatetime
	}
	if hour >= 24 {
		tempDatetime = tempDatetime.AddDate(0, 0, 1)
	} else {
		minute := (subTime - 3600*hour) / 60
		second := subTime - 3600*hour - minute*60
		tempDatetime = tempDatetime.Add(time.Duration(hour) * time.Hour)
		tempDatetime = tempDatetime.Add(time.Duration(minute) * time.Minute)
		tempDatetime = tempDatetime.Add(time.Duration(second) * time.Second)
	}

	return stopOffset, tempDatetime
}
