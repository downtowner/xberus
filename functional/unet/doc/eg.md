# 使用实例

```go
实例1,均采用默认的管理
func TestTcpServer(t *testing.T) {

	app := NewTCPServer()

	app.Shutdown(time.Second*5, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	app.Run("tcp", "127.0.0.1:10086", NewConnMgr(NewClientSocket))
}


```

##### 说明

例子中均采用默认的组件包

`NewConnMgr`: 实现了Connections接口

`NewClientSocket`: 对象工厂

```go
实例2,扩展mgr和clientsocket
package main

import (
	"git.vnnox.net/ncp/xframe/functional/unet"
	"log"
	"os"
	"syscall"
	"time"
)

type LocalManager struct {
	unet.ConnectionMgr
}

type LocalClient struct {
	unet.ClientSocket
}

func NewLocalClient() unet.Client {
	p := &LocalClient{}
	p.Init()
	return p
}

//OnCmdMessage handle cmd message
func (l *LocalClient) OnCmdMessage(cmd string, data []byte) error {

	log.Println("Local OnCmdMessage")



	return nil
}

//OnIDMessage handle id message
func (l *LocalClient) OnIDMessage(id int, data []byte) error {

	log.Println("Local OnIDMessage")

	log.Println("id: ", id, "data: ", string(data))

	l.SendIDMessage(id, data)
	return nil
}


func main() {

	log.Println("start")

	app := unet.NewTCPServer()

	app.Shutdown(time.Second*5, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	mgr := &LocalManager{}
	mgr.Init()
	mgr.SetConnecter(NewLocalClient)


	app.Run("tcp", "127.0.0.1:10086", mgr)

	log.Println("endl")
}
```

