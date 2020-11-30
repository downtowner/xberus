package udpproxy

import (
	"testing"
)

func TestUDPProxy(t *testing.T) {

	// flag.Set("lis", "127.0.0.1:10086")
	// flag.Set("des", "192.168.1.153:10086")
	app := NewUProxyServer()
	app.Run()
}
