package udpproxy

import (
	"fmt"
	"log"
	"net"
)

//NewTerminaler create a terminal obj
func NewTerminaler(netaddr string, msgqueue chan<- []byte) *terminaler {

	p := terminaler{}
	p.netaddr = netaddr
	p.msgqueue = msgqueue

	return &p
}

type terminaler struct {
	netaddr  string
	conn     net.Conn
	msgqueue chan<- []byte
}

func (t *terminaler) SendMsg(data []byte) {

	wCount := 0
	for {

		n, err := t.conn.Write(data[wCount:])
		if nil != err {

			log.Printf("send msg to des server[%s] err:%s", t.netaddr, err)
			break
		}
		wCount += n

		if wCount == len(data) {

			break
		}
	}
}

func (t *terminaler) Run() {

	var err error
	t.conn, err = net.Dial("udp", t.netaddr)
	if nil != err {

		log.Printf("dial server[%s] err: %s", t.netaddr, err.Error())
		return
	}

	log.Println("to terminal:", t.netaddr)

	go t.recvMsg()
}

func (t *terminaler) recvMsg() {
	for {

		msg := make([]byte, maxBufLen)
		n, err := t.conn.Read(msg)
		if nil != err {

			log.Printf("receive server[%sd] data err: %s", t.netaddr, err)
			continue
		}

		fmt.Printf("receive msg[%d] from des server[%s]: %s\n", n, t.netaddr, string(msg[:n]))

		t.msgqueue <- msg[:n]
	}
}
