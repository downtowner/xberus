# 使用实例

#### 消息中心

```go
corba.RunCenterNode(":8080")
```



#### 节点

```go
// 运行消息的子节点
corba.RunSubNodeDefault("127.0.0.1:8080")
```



#### 子节点对象注册

```go
// 注册对象
corba.RegisterObject("plugina")
```



#### 发消息

```go
// 发消息
msg := corba.Message{}
msg.Code = corba.SystemSendMessage
msg.ObjectName = "pluginb"

simpleMsg := xca.SimpleMessage{}
simpleMsg.Data = "sssss"
msg.Data = simpleMsg

_, err := corba.SendMessage(&msg)
if err != nil {
	fmt.Println(err)
}
```

