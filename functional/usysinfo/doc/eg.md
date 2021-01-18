# 系统信息组件

```golang

package main

import (
	"fmt"
	"log"
	"net/http"

	"git.vnnox.net/ncp/xframe/functional/usysinfo"
	"git.vnnox.net/ncp/xframe/xca"
	"github.com/kataras/iris/v12"
)

type SysinfoManager interface {
	// httpHandler
	HttpHandler(w http.ResponseWriter, r *http.Request)
	// iris路由包装
	IrisWrapRouter(app *iris.Application, path string)
	// 定时轮询写入 cpu使用率 内存使用量 系统平均负载 磁盘使用量 网络流入流出 tcp连接数
	CrontabWrite(storage usysinfo.StorageEngine, dayExpired int) error
	// 创建存储 只支持单机，如果是多机，需要自己重新实现接口
	NewStorage() usysinfo.StorageEngine
}

func main() {
	log.SetFlags(log.Llongfile)
	obj, err := createObj()
	if err != nil {
		log.Fatal(err)
	}
	app := iris.New()
	//注册路由
	url := "/sysinfo"
	obj.IrisWrapRouter(app, url)
	//自定义数据存储过期时间和存储引擎
	obj.CrontabWrite(obj.NewStorage(), 7)
	app.Run(iris.Addr(":8080"))
}

func createObj() (SysinfoManager, error) {
	//	本地内置信息
	var localComponents = []map[string]interface{}{
		usysinfo.SysinfoManagerComponentInfo,
	}
	// 注册内置组件
	err := xca.RegisterComponents(localComponents)
	if err != nil {
		return nil, err
	}
	obj, ok := xca.CreateNamedObject(usysinfo.ComponentName, usysinfo.ComponentName).(SysinfoManager)
	if !ok {
		return nil, fmt.Errorf("defined SysinfoManager error")
	}
	return obj, nil
}


```

## 待完善

### 1.目前存储只支持单机文件存储，但留有接口可以扩展

### 2.路由包装支持iris,其他的web框架如有需要可以添加

欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [zousf@novastar.tech](mailto:zousf@novastar.tech)
