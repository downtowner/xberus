package udpproxy

import (
	"errors"
	"log"
	"strings"
)

//NewUProxyServer ...
func NewUProxyServer() *UDPProxyServer {

	p := UDPProxyServer{}

	return &p
}

//UDPProxyServer ...
type UDPProxyServer struct {

	// listen on which net address
	listenAddress string

	// forward data to target server
	targetAddress []string
}

// run run service
func (u *UDPProxyServer) Run(listen, targets string) error {

	if 0 == len(listen) || 0 == len(targets) {

		log.Println("run start params error")
		return errors.New("run start params error")
	}

	u.listenAddress = listen

	u.targetAddress = strings.Split(targets, ";")
	if 0 == len(u.targetAddress) {

		log.Println("parse targets param err,eg: xxx.xxx.xx.x:a;xxx.xxx.x.x:b")
		return errors.New("parse targets param err,eg: xxx.xxx.xx.x:a;xxx.xxx.x.x:b")
	}

	c := NewConnector(u.listenAddress)
	err := c.Listen()
	if nil != err {

		return err
	}

	mgr := NewTerminalerMgr(u.targetAddress, c.TerminalChannel(), c.OriginChannel())
	mgr.Run()

	return nil
}
