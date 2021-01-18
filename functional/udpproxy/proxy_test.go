package udpproxy

import (
	"flag"
	"log"
	"testing"
)

//start params:-lis=127.0.0.1:10086 -des=192.168.0.12
func TestUDPProxy(t *testing.T) {

	var listen string
	flag.StringVar(&listen, "lis", "", "local host listen address")

	var targets string
	flag.StringVar(&targets, "des", "", "transport to target host, please use ';' to split multiple host")

	flag.Parse()

	log.Println("start up params:", listen, targets)
	app := NewUProxyServer()
	app.Run(listen, targets)
}
