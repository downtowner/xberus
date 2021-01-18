# 简要说明

为xca提供跨进程消息服务，所需的组件

# 详细说明

通过socket通讯，为xca提供跨进程消息服务，消息的内容不限制可以是字符串或对象等，理论上消息体的大小也没有限制，但是大对象不建议使用，此组件包含 三个主要模块：

### CenterNode

消息中心模块，该模块负责消息调度，主要接口：

运行：

`func (cn *CenterNode) Run(addr string) error`

开始运行消息中心服务，接受1个参数，定义如下：

`addr`：监听的IP、端口，例如：`127.0.0.1:8080`



### SubNode

子节点模块，即需要通讯的部份，主要接口：

运行：

`func (s *SubNode) Run(addres string) error`

开始运行消息中心服务，接受1个参数，定义如下：

`addres `：消息中心监听的IP、端口，例如：`127.0.0.1:8080`



创建节点标识：

`func (s *SubNode) CreateTag() string `

返回一个唯一的字符串，非必需，也可以自行实现



设置节点标识：

`func (s *SubNode) SetNodeTag(tag string) error`

设置节点标识，tag必须是唯一的



获得运行状态：

`func (s *SubNode) GetState() bool`

只有节点在运行状态下才能被正确的发消息、推消息



注册对象：

`func (s *SubNode) RegisterObject(obj string) error`

将需要被通讯的对象名称注册到消息中心



发消息：

`func (s *SubNode) SendMessage(msg IMessage) (interface{}, error)`



推消息：

`func (s *SubNode) PostMessage(msg IMessage) error`



### MessageManager

消息模块，处理消息，格式化通讯包的内容，主要接口有：

格式化包：

`func (mm *MessageManager) FormatMsgBody(msg IMessage) []byte`







# 待完善

模块暂时只实现了核心功能，该模块还在扩展之中，目前的不足之处有：
##### 1.同步消息
目前的消息都是异步消息

##### 2.发消息的回传问题

组件有发消息，和推消息的功能，目前发消息未实现消息成功后的信息回传功能

##### 3、单方面断线，重联问题

当中心节点或子节点单方面重启、断线时的重联问题目前没有处理

##### 4. socket模块和解包模块的替换

目前这两个模块是自己独立编写，考虑替换为unet通讯，upackage解包



欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [luozp@novastar.tech](mailto:luozp@novastar.tech)





























