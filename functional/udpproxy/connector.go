package udpproxy

import (
	"fmt"
	"log"
	"net"
)

const (
	maxBufLen   int = 1024
	maxMsgCount int = 256
)

//NewConnector docking class between orgin server and terminal server
func NewConnector(netaddr string) *connector {

	p := connector{}
	p.netAddress = netaddr
	p.toTerminals = make(chan []byte, maxMsgCount)
	p.toOrigin = make(chan []byte, maxMsgCount*10)

	return &p
}

type connector struct {
	netAddress string

	conn *net.UDPConn

	connAddr *net.UDPAddr

	toTerminals chan []byte
	toOrigin    chan []byte
}

func (c *connector) Listen() error {

	addr, err := net.ResolveUDPAddr("udp", c.netAddress)
	if err != nil {

		fmt.Println(err)
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {

		fmt.Println(err)
		return err
	}

	log.Println("listen net address: ", c.netAddress)

	c.conn = conn

	go c.recvOriginMsg()
	go c.recvTerminalMsg()

	return nil
}

func (c *connector) recvTerminalMsg() {

	for {
		select {

		case data := <-c.toOrigin:

			wCount := 0
			for {

				n, err := c.conn.WriteToUDP(data[wCount:], c.connAddr)
				if nil != err {

					log.Printf("send msg to client server[%s] err:%s", c.connAddr.String(), err)
					break
				}
				wCount += n

				if wCount == len(data) {

					break
				}
			}
		}
	}
}

func (c *connector) recvOriginMsg() {

	for {

		data := make([]byte, maxBufLen)
		dlen, addr, err := c.conn.ReadFromUDP(data)
		if err != nil {

			fmt.Println(err)
			continue
		}

		log.Printf("recv client[%s] server data[%d]: %s\n", addr.String(), dlen, string(data[:dlen]))

		c.connAddr = addr

		c.toTerminals <- data[:dlen]
	}
}

func (c *connector) TerminalChannel() chan []byte {

	return c.toTerminals
}

func (c *connector) OriginChannel() chan []byte {

	return c.toOrigin
}
