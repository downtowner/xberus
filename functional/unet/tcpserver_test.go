package unet

import (
	"os"
	"syscall"
	"testing"
	"time"
)

func TestTcpServer(t *testing.T) {

	app := NewTCPServer()

	app.Shutdown(time.Second*5, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	app.Run("tcp", "127.0.0.1:10086", NewConnMgr(NewClientSocket))
}
