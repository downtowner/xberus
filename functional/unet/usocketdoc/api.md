# 简要说明

和unet组件包配套使用的客户socket封装

# 详细说明

和unet组件包配套使用的客户socket封装，基于TCP协议，对数据包进行编解码，主要接口有：

```
func (u *UltralSocket) Connect(params ...interface{}) error
```

连接到目标服务器,该函数的一个参数会根据创建的模式而定，若是异步模式，则第一个参数必须为`func(int8, int32, string, []byte) error`函数

```
func (u *UltralSocket) SendIDMessage(id int, data []byte) error
```

发送ID模式的数据包，具体参见upackage组件

```
func (u *UltralSocket) SendCmdMessage(cmd string, data []byte) error 
```

发送CMD模式的数据包，具体参见upackage组件

```
func (u *UltralSocket) SendAndRecvIDMsg(id int, data []byte, handle func(int, []byte) error) error
```

同步发送和接收id模式数据包

```
func (u *UltralSocket) SendAndRecvCmdMsg(cmd string, data []byte, handle func(string, []byte) error) error
```

同步发送和接收CMD模式数据包

```
func (u *UltralSocket) Close()
```

关键客户连接



# 注意

```
func NewUSocket(addr string, syn bool) *UltralSocket
```

创建该套接字需要指定同步还是异步模式，如果是异步模式，内置方法`loop`会监听数据包，`SendAndRecvIDMsg`接口不能使用，它属于同步方法。



欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [xiaobing@novastar.tech](mailto:moubo@novastar.tech)