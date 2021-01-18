package corba

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

// CenterNode 中心节点
type CenterNode struct {
	mutex           sync.Mutex
	messageManager  IMessageManager
	ConnectCacheMap map[string]chan []byte // ConnectCacheMap 连接缓存的map，key是tag ,value是writeChan
	ObjcetCacheMap  map[string]string      // 注册对象的缓存 key是对象名称 value是节点的tag
}

var centerNode *CenterNode

// NewCenterNode return a *CenterNode
func NewCenterNode() *CenterNode {
	centerNode = &CenterNode{}
	centerNode.ConnectCacheMap = make(map[string]chan []byte)
	centerNode.ObjcetCacheMap = make(map[string]string)
	centerNode.messageManager = &MessageManager{}

	return centerNode
}

// Run 启动
func (cn *CenterNode) Run(addr string) error {
	tcpServer, _ := net.ResolveTCPAddr("tcp4", addr)
	listener, _ := net.ListenTCP("tcp", tcpServer)
	for {
		//当有新的客户端请求来的时候，拿到与客户端的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		//处理逻辑
		go cn.handle(conn)
	}
}

func (cn *CenterNode) handle(conn net.Conn) {
	var readChan chan []byte
	var writeChan chan []byte
	var tag string

	readChan = make(chan []byte, 256)
	writeChan = make(chan []byte, 256)

	closeChan := make(chan error)
	defer func() {
		fmt.Println("conn.Close")
		conn.Close()
	}()

	//读取客户端传送的消息
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

			readChan <- data
		}
	}()

	//监听信号
	for {
		select {
		case buf := <-readChan:
			var cmd uint8     // 命令字
			var objLen uint32 // 对象长度

			bytebuff := bytes.NewBuffer(buf[4:5])
			binary.Read(bytebuff, binary.BigEndian, &cmd)

			bytebuff.Reset()
			bytebuff = bytes.NewBuffer(buf[5:9])
			binary.Read(bytebuff, binary.BigEndian, &objLen)

			switch cmd {
			case SystemSetNodeTag: // 设置tag
				tag = string(buf[13 : 13+objLen])
				cn.mutex.Lock()
				cn.ConnectCacheMap[tag] = make(chan []byte, 256)
				cn.mutex.Unlock()
				//fmt.Println("设置tag :", tag)

				// 启动数据转发的信号
				go func() {
					for {
						select {
						case b := <-cn.ConnectCacheMap[tag]:
							writeChan <- b
						}
					}

				}()

			case SystemRegisterObject: //注册对象
				obj := string(buf[13 : 13+objLen])
				cn.mutex.Lock()
				cn.ObjcetCacheMap[obj] = tag
				cn.mutex.Unlock()

			case SystemSendMessage:
				obj := string(buf[13 : 13+objLen])
				connectKey, ok := cn.ObjcetCacheMap[obj]
				if !ok {

					break
				}
				cn.ConnectCacheMap[connectKey] <- buf

			case SystemPostMessage:
				obj := string(buf[13 : 13+objLen])
				connectKey, ok := cn.ObjcetCacheMap[obj]
				if !ok {

					break
				}
				cn.ConnectCacheMap[connectKey] <- buf

			}

		case writeBuf := <-writeChan:
			conn.Write(writeBuf)

		case err := <-closeChan:
			fmt.Println("read failed, err:", err)
			return
		}
	}
}
