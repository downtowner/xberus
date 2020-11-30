package udpproxy

import (
	"errors"
	"flag"
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
func (u *UDPProxyServer) Run() error {

	flag.StringVar(&u.listenAddress, "lis", "", "local host listen address")

	var tmp string
	flag.StringVar(&tmp, "des", "", "transport to target host, please use ';' to split multiple host")

	flag.Parse()

	if 0 == len(u.listenAddress) || 0 == len(tmp) {

		log.Println("input params error, eg: -lis=ip:port -des=ip:port;ip:port")
		return errors.New("input params error, eg: -lis=ip:port -des=ip:port;ip:port")
	}

	u.targetAddress = strings.Split(tmp, ";")
	if 0 == len(u.targetAddress) {

		log.Println("input params error, eg: -lis=ip:port -des=ip:port;ip:port")
		return errors.New("input params error, eg: -lis=ip:port -des=ip:port;ip:port")
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
