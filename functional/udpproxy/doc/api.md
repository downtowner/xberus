# 简要信息

一个简单的1对多UDP代理服务组件

# 详细信息

一个简单的1对多UDP代理服务组件，用来解决远程调试不方便的问题，消息流程如下：

![avatar][png]

该组件主要模块有:

### connector

和目标服务对接的模块，主要的接口：

```
func (c *connector) Listen() error
```

开始监听数据源服务器

```
func (c *connector) TerminalChannel() chan []byte
```

获得需要发送到目标服务器的消息管道

```
func (c *connector) OriginChannel() chan []byte
```

获得需要发送到源服务器的消息管道

### proxy

对各个子模块进行管理和调度

```
func (u *UDPProxyServer) Run() error
```

运行代理服务器实例



### terminaler

每个目标服务器的抽象管理，主要接口有:

```
func (t *terminaler) SendMsg(data []byte)
```

向目标终端发送数据

```
func (t *terminaler) Run()
```

开始监听目标终端回应的数据

### terminalermgr

对N个终端进行管理，主要接口有:

```
func (t *terminalMgr) Run()
```

初始化并开始运行管理模块





欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [xiaobing@novastar.tech](mailto:moubo@novastar.tech)



[png]:data:image/png;base64,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