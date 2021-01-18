unet
# 简要说明

提供高并发，高性能服务所需的核心功能组件

# 详细说明

提供一个高并发，高性能，轻量级的服务所需的核心功能组件，此组件包含三个主要模块：

### TCPServer

服务的核心框架模块，主要接口：

`func (t *TCPServer) Run(network, address string, hc Connections)`

功能：开始运行服务，接受3个参数，意义分别如下：

`network`: 协议类型, 后期可扩展

`address `: 监听的网络地址

`hc`: 处理链接的接口实现，该接口具体如下:

```go
type Connections interface {
HandleConnections(context.Context, *net.TCPConn)
Done()
}
```

`HandleConnections(context.Context, *net.TCPConn)`: 处理请求链接，该包已经提供默认的处理模块`ConnectionMgr `,可以方便的进行扩展

Done()：优雅的退出接口   



`func (t *TCPServer) Shutdown(timeout time.Duration, sig ...os.Signal) `

功能：关闭服务器，接受2个参数，意义分别如下：

`timeout `: 退出的超时时间

`sig`: 退出监听的信号

### ConnectionMgr

管理连接请求的核心模块，可以重写其行为，实现了Connections接口，主要接口有：

`func (c *ConnectionMgr) enableBroadcast() `:  启用广播功能

```go
type Client interface {
	Initialize(ctx context.Context, conn *net.TCPConn, checkTime int)
	Done()
}   
```

功能：定义了请求连接的必要行为

`Initialize(ctx context.Context, conn *net.TCPConn, checkTime int)` 

初始化函数,参数意义如下:

`ctx`: 上下文

`checkTime `: 读数据的检查时间

`conn`: 连接指针



`Done()`

退出连接



`func (c *ConnectionMgr) Init()`

初始化数据，通过`NewConnMgr`创建的mgr对象已经被初始化

`func (c *ConnectionMgr) SetConnecter(of ObjFactory)`

设置连接对象工厂,通过`NewConnMgr`创建的mgr对象可以不需要调用此过程

### ClientSocket 

默认的客户连接请求处理模块, 模块中采用自定义协议 `upackage` 进行管理,具体参见该组件详细信息，主要接口有：

`func (c *ClientSocket) loop()`

异步接收消息

`func (c *ClientSocket) receivePackage() error`

消息解码

`func (c *ClientSocket) HandleMessage(mark int8, id int, cmd string, data []byte) error`

消息处理，`mark` 消息格式，`id` 消息id，`cmd` 消息标识, `data` 消息内容,具体参见`upackage`

`func (c *ClientSocket) BroadcastIDMsg(id int, data []byte)`

`func (c *ClientSocket) BroadcastCmdMsg(cmd string, data []byte)`

广播消息接口



# 待完善

模块暂时只实现了核心功能，该模块还在扩展之中，当然你可以根据实际情况对模块进行扩展，目前的不足之处有：

##### 1.消息的调试打印问题

该服务是基于传输层协议而进行直接解码，所以不像webserver那样可以方便的对接口进行调试。

##### 2. 封装粒度问题

该模块依然处理开发状态,后期的版本可能会针对通用的业务进行更精细的封装。

##### 3.计划

绿色精简的服务核心框架，可以定制功能模块，代码简单易读，下一步会优先融入对缓存、数据库、CGO封装的实现。



欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [xiaobing@novastar.tech](mailto:moubo@novastar.tech)


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

