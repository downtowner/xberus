package corba

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"git.vnnox.net/ncp/xframe/xca"
)

// SubNode 子节点
type SubNode struct {
	readChan  chan []byte
	writeChan chan []byte
	state     bool
}

var subnode *SubNode

// NewSubNode return a *SubNode
func NewSubNode() *SubNode {
	subnode = &SubNode{}
	messageManager = &MessageManager{}

	subnode.readChan = make(chan []byte, 256)
	subnode.writeChan = make(chan []byte, 256)

	return subnode
}

// CreateTag 创建节点标识
func (s *SubNode) CreateTag() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}

	str := base64.URLEncoding.EncodeToString(b)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GetState 状态，为true时才可以正常发送消息
func (s *SubNode) GetState() bool {
	return s.state
}

// SetNodeTag 设置节点标识
func (s *SubNode) SetNodeTag(tag string) error {

	msg := Message{}
	msg.Code = SystemSetNodeTag
	msg.ObjectName = tag

	// xcaMsg := XCAMessage{}
	// xcaMsg.Data = tag
	// msg.Data = xcaMsg

	s.writeChan <- messageManager.FormatMsgBody(&msg)
	return nil
}

// RegisterObject 注册对象
func (s *SubNode) RegisterObject(obj string) error {
	msg := Message{}
	msg.Code = SystemRegisterObject
	msg.ObjectName = obj

	// xcaMsg := XCAMessage{}
	// xcaMsg.Data = tag
	// msg.Data = xcaMsg

	s.writeChan <- messageManager.FormatMsgBody(&msg)
	return nil
}

// SendMessage 发消息
func (s *SubNode) SendMessage(msg IMessage) (interface{}, error) {
	if !s.state {
		return nil, errors.New("The child node is not running")
	}

	body := messageManager.FormatMsgBody(msg)
	s.writeChan <- body
	return nil, nil
}

// PostMessage 推消息
func (s *SubNode) PostMessage(msg IMessage) error {
	body := messageManager.FormatMsgBody(msg)
	s.writeChan <- body
	return nil
}

// Run 运行
func (s *SubNode) Run(addres string) error {

	addr, err := net.ResolveTCPAddr("tcp4", addres)
	checkError(err)

	//建立tcp连接
	conn, err := net.DialTCP("tcp4", nil, addr)
	checkError(err)

	closeChan := make(chan error)

	s.state = true

	// 读数据
	go func() {
		for {

			// 读包头，确定包后续包的长度
			bufHead := make([]byte, 4)
			count, err := conn.Read(bufHead)
			if err != nil {
				closeChan <- err

				break
			}

			bytebuff := bytes.NewBuffer(bufHead[0:count])
			var bodyLength int32 //总包长度
			binary.Read(bytebuff, binary.BigEndian, &bodyLength)

			// 读包体
			bufBody := make([]byte, bodyLength-4)
			countBody, err := conn.Read(bufBody)
			if err != nil {
				closeChan <- err

				break
			}

			var data []byte
			data = append(data, bufHead[0:count]...)
			data = append(data, bufBody[0:countBody]...)

			s.readChan <- data

		}
	}()

	for {
		select {
		case buf := <-s.readChan:
			s.handleRead(buf)

		case writeBuf := <-s.writeChan:
			conn.Write(writeBuf)

		case err := <-closeChan:
			fmt.Println("read failed, err:", err)
		}
	}
}

// 处理读消息
func (s *SubNode) handleRead(msgData []byte) {

	// 解包
	var cmd uint8      // 命令字
	var objLen uint32  // 对象长度
	var msgLen uint32  // IMessage序列化后数据的长度
	var objName string // 消息对象名称
	var msgByte []byte // IMessage序列化后的数据

	bytebuff := bytes.NewBuffer(msgData[4:5])
	binary.Read(bytebuff, binary.BigEndian, &cmd)

	bytebuff.Reset()
	bytebuff = bytes.NewBuffer(msgData[5:9])
	binary.Read(bytebuff, binary.BigEndian, &objLen)

	bytebuff.Reset()
	bytebuff = bytes.NewBuffer(msgData[9:13])
	binary.Read(bytebuff, binary.BigEndian, &msgLen)

	objName = string(msgData[13 : 13+objLen])
	msgByte = msgData[13+objLen : 13+objLen+msgLen]

	// 反序列化得到msg
	var network bytes.Buffer // 替代网络连接
	network.Write(msgByte)
	dec := gob.NewDecoder(&network) // 将从网络上读取。
	var msg Message
	err := dec.Decode(&msg)
	if err != nil {
		log.Fatal("decode error 1:", err)
	}

	switch cmd {
	case SystemSendMessage:
		// 调用xca发消息
		xca.SendMessage(objName, &msg.Data)

	case SystemPostMessage:
		// 调用xca发消息
		xca.PostMessage(objName, &msg.Data)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
