// Copyright (C) 2020 The Nova Company Ltd.All rights reserved.

package usysinfo

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/shirou/gopsutil/v3/host"

	_ "git.vnnox.net/ncp/xframe/functional/usysinfo/statik"

	"github.com/rakyll/statik/fs"
)

// ComponentNameSysinfoManager manager
const ComponentName = "sysinfoManager.sysinfo.novastar.tech"

var components = []map[string]interface{}{
	{
		"name":    ComponentName,
		"version": "1.0",
		"creator": func() interface{} { return newSysinfoObj() },
	},
}

// SysinfoManagerComponentInfo info
var SysinfoManagerComponentInfo = map[string]interface{}{
	"name":       ComponentName,
	"version":    "1.0",
	"creator":    func() interface{} { return newSysinfoObj() },
	"components": components,
}

type SysinfoManager interface {
	// httpHandler
	HttpHandler(w http.ResponseWriter, r *http.Request)
	// iris路由包装
	IrisWrapRouter(app *iris.Application, path string)
	// 定时轮询写入 cpu使用率 内存使用量 系统平均负载 磁盘使用量 网络流入流出 tcp连接数
	CrontabWrite(storage StorageEngine, dayExpired int) error
	// 创建存储 只支持单机，如果是多机，需要自己重新实现接口
	NewStorage() StorageEngine
}

type sysinfoObj struct {
	// 对外暴露的接口
	SysinfoManager
	// 存储引擎
	storage StorageEngine
	// 静态文件
	fileSystem http.FileSystem
	// url路由绑定路径
	urlPath string
}

func GetComponents() map[string]interface{} {
	return SysinfoManagerComponentInfo
}

func newSysinfoObj() interface{} {
	statikFS, err := fs.New()
	if err != nil {
		return nil
	}
	s := &sysinfoObj{
		fileSystem: statikFS,
	}
	return s
}

// httpHandler 目前iris的handler是基于当前方法封装
// 如果使用其他web框架也可以自己封装
func (s *sysinfoObj) HttpHandler(w http.ResponseWriter, r *http.Request) {
	// 静态资源处理
	if regexp.MustCompile(`\.(js|css|json|ico|gif)+$`).MatchString(r.URL.Path) {
		s.assets(w, r)
		return
	}

	api, ok := r.URL.Query()["field"]
	if !ok || len(api) < 1 {
		// 页面
		if len(s.urlPath) == 0 {
			s.urlPath = r.URL.Path
		}
		s.htmlHandle(w, r)
		return
	}
	// api
	s.apiHandle(w, r)
	return
}

// 初始化存储引擎并且定时入库，目前1分钟执行1次
// 采集的信息其实都是去调系统api去查询系统信息，短时间大量的查询会有系统损耗，所以这里不开放入库的时间控制
func (s *sysinfoObj) CrontabWrite(storage StorageEngine, dayExpired int) error {
	writeInterval := time.Minute
	s.storage = storage
	err := s.storage.InitEngine()
	if err != nil {
		return err
	}
	go func() {
		for {
			hostInfo, err := host.InfoWithContext(context.TODO())
			if err != nil {
				continue
			}
			cpuinfo := CpuInfo{}
			memoryInfo := MemoryInfo{}
			newWorkInfo := NewWorkInfo{}
			diskInfo := DiskInfo{}
			goRuntimeInfo := GoRuntimeInfo{}
			IOInfo := IOInfo{}
			storageInfo := StorageInfo{
				Time:      time.Now().Unix(),
				Hostname:  hostInfo.Hostname,
				Cpu:       cpuinfo.Collect(),
				Memory:    memoryInfo.Collect(),
				NetWork:   newWorkInfo.Collect(),
				Disk:      diskInfo.Collect(),
				GoRuntime: goRuntimeInfo.Collect(),
				IO:        IOInfo.Collect(),
			}
			s.storage.Set(storageInfo)
			<-time.After(writeInterval)
		}
	}()
	go func() {
		for {
			day := FormatTime(time.Now().AddDate(0, 0, -dayExpired), "20060102")
			s.storage.Delete(day)
			<-time.After(12 * time.Hour)
		}
	}()
	return nil
}
