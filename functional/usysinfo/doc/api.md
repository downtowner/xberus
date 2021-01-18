# 系统信息组件

## 应用场景

- 监控软件系统运行状态是否正常
- 查看某一时间段软件系统的运行状态
- 用于排查问题时候，能清楚知道各个时间段的系统性能指标

## 注意点

buildInfo.go有3个字段需要在编译的时候动态注入

GitSHA   git SHA信息

GitURL   git URL信息

BuildTime 编译时间

```bash
go build -ldflags \
"-X git.vnnox.net/ncp/xframe/usysinfo.GitSHA=`git rev-parse HEAD` \
 -X git.vnnox.net/ncp/xframe/usysinfo.BuildTime=`date +%FT%T%z`" \
 -X git.vnnox.net/ncp/xframe/usysinfo.GitURL=yourGitURL "\
-o main *.go
```

或者使用Makefile

```Makefile
export PATH:=${PATH}:${GOPATH}/bin:

OUTPUT=main

GITCOMMITID=`git rev-parse HEAD`
BUILDTIME=`date +%FT%T%z`
GITURL=`http://xxxx`

LDFLAGS=-ldflags "-X git.vnnox.net/ncp/xframe/usysinfo.GitSHA=${GITCOMMITID} -X git.vnnox.net/ncp/xframe/usysinfo.BuildTime=${BUILDTIME} -X git.vnnox.net/ncp/xframe/usysinfo.GitURL=${GITURL}"

linux:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${OUTPUT} *.go

```


## 采集的数据

### cpu 

- cpu使用率

### 内存

- 内存使用率

### 硬盘

- 硬盘使用率
- 硬盘Inodes使用率（linux才有）

### goruntime

- cpu.count  cpu数
- cpu.goroutines  协程数
- cpu.cgoCall  cgo调用次数
- mem.alloc  Alloc是分配的堆对象的字节
- mem.total  TotalAlloc是分配给堆对象的累积字节
- mem.sys  OS获得的内存的总字节数
- mem.lookups  被runtime监视的指针数
- mem.malloc  Mallocs是分配的堆对象的累积计数
- mem.frees  Frees是已释放的堆对象的累积计数
- mem.heap.alloc  HeapAlloc是分配的堆对象的字节
- mem.heap.sys  从操作系统获得的堆内存字节
- mem.heap.idle  HeapIdle减去HeapReleased估计内存量
- mem.heap.inuse  HeapInuse是使用中的跨度中的字节
- mem.heap.released  HeapReleased是返回操作系统的物理内存字节
- mem.heap.objects  HeapObjects是分配的堆对象的数量
- mem.stack.inuse  StackInuse是堆栈跨度中的字节
- mem.stack.sys  StackSys是从操作系统获得的堆栈内存字节
- mem.stack.mspanInuse  MSpanInuse是分配的mspan结构的字节
- mem.stack.mspanSys  MSpanSys是从操作系统获取的mspan的内存字节
- mem.stack.mcacheInuse  MCacheInuse是分配的mcache结构的字节
- mem.stack.mcacheSys  MCacheSys是从OS获得的内存字节
- mem.othersys  OtherSys是其他堆中的内存字节
- gc.numGC  NumGC是已完成的GC周期数

### 网络

- 每秒接收字节数
- 每秒发送字节数
- 收包错误率
- 发包错误率
- 正确收包数
- 正确发包数

### IO

- 硬盘IO每秒读次数
- 硬盘IO每秒写次数
- 硬盘IO每秒读字节数
- 硬盘IO每秒写字节数

### 系统信息

- Go版本
- 系统平台
- 系统内核
- CPU
- CPU特性
- 内存
- 进程启动时间
- 系统启动时间
- 采集数据库文件大小
- 编译时间
- gitSHA
- gitURL

## 采集周期

目前1分钟采集1次

## 折线图显示粒度

| 时间范围  | 粒度 |
| ----- | ------ |
| 0-1分钟   | 1分钟 |
| 1到6分钟  | 5分钟 |
| 6分钟以上  | 15分钟 |

## 支持系统
- windows(64位)
- linux(64位)
- darwin(64位)

## 数据存储

目前采用bolt存储，db文件名sysinfo.db


## 接口

```golang
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
```


### HttpHandler

`HttpHandler(w http.ResponseWriter, r *http.Request)`

作用:http handler,这里实现的就是标准库net/http的Handler,所有的web框架都可以基于当前handler做二次封装


### IrisWrapRouter

`IrisWrapRouter(app *iris.Application, path string)`

作用:iris包装HttpHandler

#### 传入参数:

app:iris v12 的对象，也就是iris.New()

path:绑定的路由

#### 返回参数


### CrontabWrite

`CrontabWrite(storage StorageEngine, dayExpired int) (err error)`

作用:定时轮询写入 cpu使用率 内存使用量 系统平均负载 磁盘使用量 网络流入流出 tcp连接数

#### 传入参数:

storage:组件里自己实现了存储方法，当然，你也可以自己实现自己的存储的方法

```golang
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
```

dayExpired:日志保留天数

#### 返回参数

err:错误信息


### NewStorage

`NewStorage() StorageEngine`

作用:创建存储 只支持单机

#### 传入参数:

#### 返回参数

```golang
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
```

## 待完善

### 1.目前存储只支持单机文件存储，但留有接口可以扩展

### 2.路由包装支持iris,其他的web框架如有需要可以添加

欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [zousf@novastar.tech](mailto:zousf@novastar.tech)
