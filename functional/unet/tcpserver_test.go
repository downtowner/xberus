package unet

import (
	"log"
	"os"
	"syscall"
	"testing"
	"time"
)

type LocalManager struct {
	ConnectionMgr
}

type LocalClient struct {
	ClientSocket
}

func NewLocalClient() Client {
	p := &LocalClient{}
	p.Init()
	return p
}

//OnCmdMessage handle cmd message
func (l *LocalClient) OnCmdMessage(cmd string, data []byte) error {

	log.Println("Local OnCmdMessage")

	return nil
}

//OnIDMessage handle id message
func (l *LocalClient) OnIDMessage(id int, data []byte) error {

	//log.Println("Local OnIDMessage", l.RemoteAddress(), "id: ", id, "data: ", string(data))

	l.SendIDMessage(id, data)
	return nil
}

func TestTcpServer(t *testing.T) {

	log.Println("start")

	app := NewTCPServer()

	app.Shutdown(time.Second*5, os.Interrupt, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	// mgr := &LocalManager{}
	// mgr.Init()
	// mgr.SetConnecter(NewLocalClient)

	app.Run("tcp", "127.0.0.1:10086", NewConnMgr(NewClientSocket))

	log.Println("endl")
}
