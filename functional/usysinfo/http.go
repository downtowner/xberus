package usysinfo

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/rakyll/statik/fs"
)

const (
	_ = iota
	storageError
	jsonDecodeError
	dateFormatError
	dateNumError
)

var ErrorHttpCode = map[int]string{
	storageError:    "storage未初始化",
	jsonDecodeError: "json decode 错误",
	dateFormatError: "日期格式不正确",
	dateNumError:    "结束日期不能小于开始日期",
}

// 静态资源处理
func (s *sysinfoObj) assets(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.SplitAfterN(r.URL.Path, s.urlPath, 2)
	if len(urlPath) != 2 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("urlPath != 2"))
		return
	}
	staticsData, err := s.assetsReadFile(urlPath[1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(r.URL.Path) >= 4 {
		suffix := strings.SplitAfterN(r.URL.Path[len(r.URL.Path)-4:], ".", 2)
		if len(suffix) == 2 {
			switch suffix[1] {
			case "css":
				w.Header().Add("content-type", "text/css")
			case "gif":
				w.Header().Add("content-type", "image/gif")
			case "ico":
				w.Header().Add("content-type", "image/x-icon")
			case "js":
				w.Header().Add("content-type", "application/javascript")
			case "json":
				w.Header().Add("content-type", "application/javascript")
			}
		}
	}
	w.Write(staticsData)
}

func (s *sysinfoObj) assetsReadFile(name string) ([]byte, error) {
	if runtime.GOOS == "windows" {
		return fs.ReadFile(s.fileSystem, strings.ReplaceAll(name, "/", "\\"))
	}
	return fs.ReadFile(s.fileSystem, name)
}

func (s *sysinfoObj) htmlHandle(w http.ResponseWriter, r *http.Request) {
	originContext, err := s.assetsReadFile("/index.html")
	if err != nil {
		originContext = []byte("fileSystem is error")
	}
	htmlContext := bytes.ReplaceAll(originContext, []byte("./"), []byte(r.URL.Path+"/"))
	w.Write(htmlContext)
}

func (s *sysinfoObj) apiHandle(w http.ResponseWriter, r *http.Request) {
	type apiResponse struct {
		Cpu        []ChartPoints `json:"cpu"`
		Disk       []ChartPoints `json:"disk"`
		Memory     []ChartPoints `json:"memory"`
		Network    []ChartPoints `json:"network"`
		GoRuntime  []ChartPoints `json:"goruntime"`
		IO         []ChartPoints `json:"io"`
		SystemInfo [][]string    `json:"systemInfo"`
	}
	originBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WFailure(w, jsonDecodeError)
		return
	}
	type ApiRequest struct {
		StartDate string `json:"startdate"`
		EndDate   string `json:"enddate"`
	}
	var (
		apiRequest      ApiRequest
		defaultDateFlag = false
	)

	err = json.Unmarshal(originBody, &apiRequest)
	nowTime := time.Now()
	startDateTime, err := StrToTime(apiRequest.StartDate, "2006/01/02 15:04:05")
	if err != nil {
		defaultDateFlag = true
	}
	endDateTime, err := StrToTime(apiRequest.EndDate, "2006/01/02 15:04:05")
	if err != nil {
		defaultDateFlag = true
	}
	if defaultDateFlag {
		//默认1小时
		startDateTime = nowTime.Add(-time.Hour)
		endDateTime = nowTime
	}
	if s.storage == nil {
		WFailure(w, storageError)
		return
	}
	if startDateTime.After(endDateTime) {
		WFailure(w, dateNumError)
		return
	}

	responseData := apiResponse{
		Cpu:        make([]ChartPoints, 0),
		Memory:     make([]ChartPoints, 0),
		GoRuntime:  make([]ChartPoints, 0),
		Network:    make([]ChartPoints, 0),
		Disk:       make([]ChartPoints, 0),
		IO:         make([]ChartPoints, 0),
		SystemInfo: make([][]string, 0),
	}
	allReqField, ok := r.URL.Query()["field"]
	if !ok || len(allReqField) == 0 {
		WFailure(w, jsonDecodeError)
		return
	}
	reqField := allReqField[0]

	timeInterval := 15 * time.Minute
	hourInterval := endDateTime.Sub(startDateTime).Hours()
	switch int(hourInterval) {
	case 0, 1:
		timeInterval = time.Minute
	case 2, 3, 4, 5, 6:
		timeInterval = 5 * time.Minute
	default:
		timeInterval = 15 * time.Minute
	}

	//cpu
	if strings.Contains(reqField, "cpu") {
		storageInfo := s.storage.GetCpu(startDateTime, endDateTime)
		for _, value := range cpuField {
			responseData.Cpu = value.CpuFilter(responseData.Cpu, storageInfo, timeInterval)
		}
	}
	//memory
	if strings.Contains(reqField, "memory") {
		storageInfo := s.storage.GetMemory(startDateTime, endDateTime)
		for _, value := range memoryField {
			responseData.Memory = value.MemoryFilter(responseData.Memory, storageInfo, timeInterval)
		}
	}
	//goruntime
	if strings.Contains(reqField, "goruntime") {
		storageInfo := s.storage.GetGoRuntime(startDateTime, endDateTime)
		for _, value := range goRuntimeField {
			responseData.GoRuntime = value.GoRuntimeFilter(responseData.GoRuntime, storageInfo, timeInterval)
		}
	}
	//disk
	if strings.Contains(reqField, "disk") {
		storageInfo := s.storage.GetDisk(startDateTime, endDateTime)
		for diskName, storeage := range storageInfo {
			for _, value := range diskField {
				responseData.Disk = value.DiskMapFilter(responseData.Disk, storeage, timeInterval, diskName)
			}
		}
	}
	//network
	if strings.Contains(reqField, "network") {
		storageInfo := s.storage.GetNetwork(startDateTime, endDateTime)
		for networkName, storeage := range storageInfo {
			for _, value := range networkField {
				responseData.Network = value.NetworkMapFilter(responseData.Network, storeage, timeInterval, networkName)
			}
		}
	}
	//io
	if strings.Contains(reqField, "io") {
		storageInfo := s.storage.GetIO(startDateTime, endDateTime)
		for _, value := range ioField {
			responseData.IO = value.IOFilter(responseData.IO, storageInfo, timeInterval)
		}
	}
	//systeminfo
	if strings.Contains(reqField, "systeminfo") {
		responseData.SystemInfo = GetSystemInfo()
	}
	WSuccess(w, responseData)
}

func GetSystemInfo() [][]string {
	systemInfo := make([][]string, 0)

	cpuName := "unknown"
	cpuFlags := "unknown"
	systemPlatform := "unknown"
	memoryTotal := "unknown"
	systemBootTime := "unknown"
	processBootTime := "unknown"
	kernelVersion := "unknown"
	dbSize := "unknown"
	cpuInfo, err := cpu.Info()
	if err == nil {
		cpuName = fmt.Sprintf("%s %s", runtime.GOARCH, cpuInfo[0].ModelName)
		cpuFlags = strings.Join(cpuInfo[0].Flags, " ")
	}
	bootTime, err := host.BootTimeWithContext(context.Background())
	if err == nil {
		systemBootTime = FormatTime(time.Unix(int64(bootTime), 0), "2006/01/02 15:04:05 UTC+8")
	}

	memoryInfo := MemoryInfo{}
	meminfo := memoryInfo.Collect()
	memoryTotal = fmt.Sprintf("%d MB", meminfo.Total/1024/1024)

	hostInfo, err := host.Info()
	if err == nil {
		kernelVersion = fmt.Sprintf("%s %s", hostInfo.KernelVersion, hostInfo.KernelArch)
		systemPlatform = fmt.Sprintf("%s %s%s %s", hostInfo.OS, hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.PlatformFamily)
	}
	processInfo, err := process.NewProcess(int32(os.Getpid()))
	if err == nil {
		createTime, err := processInfo.CreateTime()
		if err == nil {
			processBootTime = FormatTime(time.Unix(createTime/1000, 0), "2006/01/02 15:04:05 UTC+8")
		}

	}
	dbSize = fmt.Sprintf("%.02f MB", float64(GetDBSize())/1024/1024)

	systemInfo = append(systemInfo, []string{"Go版本", runtime.Version()})
	systemInfo = append(systemInfo, []string{"系统平台", systemPlatform})
	systemInfo = append(systemInfo, []string{"系统内核", kernelVersion})
	systemInfo = append(systemInfo, []string{"CPU", cpuName})
	systemInfo = append(systemInfo, []string{"CPU特性", cpuFlags})
	systemInfo = append(systemInfo, []string{"内存", memoryTotal})
	systemInfo = append(systemInfo, []string{"进程启动时间", processBootTime})
	systemInfo = append(systemInfo, []string{"系统启动时间", systemBootTime})
	systemInfo = append(systemInfo, []string{"采集数据库文件大小", dbSize})
	systemInfo = append(systemInfo, []string{"编译时间", BuildTime})
	systemInfo = append(systemInfo, []string{"gitSHA", GitSHA})
	systemInfo = append(systemInfo, []string{"gitURL", GitURL})
	return systemInfo
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func wResponse(w http.ResponseWriter, resp Response) {
	jsonResp, err := json.Marshal(resp)
	w.Header().Add("content-type", "application/json")
	if err != nil {
		w.Write([]byte("{\"code\":0, \"msg\":\"json转码错误\"}"))
		return
	}
	w.Write(jsonResp)
}

func WSuccess(w http.ResponseWriter, data interface{}) {
	response := Response{
		Msg:  "操作成功",
		Data: data,
	}
	wResponse(w, response)
}

func WFailure(w http.ResponseWriter, code int) {
	response := Response{
		Code: -1,
		Msg:  "未定义的错误",
		Data: map[string]interface{}{},
	}
	if _, ok := ErrorHttpCode[code]; ok {
		response.Code = code
		response.Msg = ErrorHttpCode[code]
	}
	wResponse(w, response)
}
