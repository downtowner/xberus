package corba

import "git.vnnox.net/ncp/xframe/xca"

// IMessage 消息体
type IMessage interface {
	GetCode() uint8
	GetObectName() string
	GetData() xca.SimpleMessage
}

// IMessageManager 消息管理
type IMessageManager interface {
	FormatMsgBody(IMessage) []byte //处理包
}

// ICenterNode 中心节点
type ICenterNode interface {
	Run(string) error // 运行
	//SendMessage(string, IMessage) error
}

// ISubNode 子节点
type ISubNode interface {
	CreateTag() string           // 创建节点标识
	SetNodeTag(string) error     // 设置节点标识
	RegisterObject(string) error // 注册对象
	Run(string) error            // 运行
	GetState() bool              // 状态，为true时才可以正常发送消息

	SendMessage(IMessage) (interface{}, error)
	PostMessage(IMessage) error
}

/* 以下内容为默认的基础实现*/
const (
	SystemSetNodeTag     uint8 = 0x00
	SystemRegisterObject uint8 = 0x01
	SystemCheckMessage   uint8 = 0x02
	SystemSendMessage    uint8 = 0x10
	SystemPostMessage    uint8 = 0x11
)
