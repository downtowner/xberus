package corba

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"

	"git.vnnox.net/ncp/xframe/xca"
)

// Message 消息
type Message struct {
	Code       uint8
	ObjectName string
	Data       xca.SimpleMessage
}

// GetCode 消息code
func (m *Message) GetCode() uint8 {
	return m.Code
}

// GetObectName 消息要发给哪个对象
func (m *Message) GetObectName() string {
	return m.ObjectName
}

// GetData 消息的内容
func (m *Message) GetData() xca.SimpleMessage {
	return m.Data
}

// MessageManager 消息处理
type MessageManager struct {
}

// FormatMsgBody 处理消息，格式化通讯包的内容
// 使用此方法封装包后，可以直接获取命令字、消息对名称，避免了不必要的反序列化操作
// 返回通讯包的内容，格式：总包长度（4字节） + 命令字（1字节）＋消息对象名称的长度（4字节）+ IMessage序列化后数据的长度（4字节）＋ 消息对象名称 ＋IMessage序列化后的数据
func (mm *MessageManager) FormatMsgBody(msg IMessage) []byte {

	var body []byte
	body = make([]byte, 4) // 用于存总包的长度

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(msg)

	var msgByte []byte
	if err == nil {
		msgByte = network.Bytes()
	}

	objName := []byte(msg.GetObectName())
	objLength := uint32(len(objName))
	msgLength := uint32(len(msgByte))

	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, msg.GetCode())
	body = append(body, bytebuf.Bytes()...)

	bytebuf.Reset()
	binary.Write(bytebuf, binary.BigEndian, objLength)
	body = append(body, bytebuf.Bytes()...)

	bytebuf.Reset()
	binary.Write(bytebuf, binary.BigEndian, msgLength)
	body = append(body, bytebuf.Bytes()...)

	body = append(body, objName...)
	body = append(body, msgByte...)

	var bodyLength uint32 //总包长度
	bodyLength = uint32(len(body))
	bytebuf.Reset()
	binary.Write(bytebuf, binary.BigEndian, bodyLength)
	headByte := bytebuf.Bytes()

	//替换包头4个字节位
	body[0] = headByte[0]
	body[1] = headByte[1]
	body[2] = headByte[2]
	body[3] = headByte[3]

	return body

}
