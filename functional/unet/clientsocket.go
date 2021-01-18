package unet

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"git.vnnox.net/ncp/xframe/functional/upackage"
)

/**
author: icarus
data: 2020-10-27

note:

**/

//NewClientSocket create a clientsocket obj, c tcp connection. params Optional configuration parameters.
/**

The meaning of the parameters is as follows:
readchecktime check socket cycle of read data, closed when it time out, 0,uncheck, unit:second

**/

//Client must be implemented method
type HandleMessage interface {

	//cmd
	OnCmdMessage(cmd string, data []byte) error

	//id
	OnIDMessage(id int, data []byte) error

	
}

func NewClientSocket() Client {

	p := ClientSocket{}
	p.Init()
	return &p
}

//ClientSocket structure of recv connection
type ClientSocket struct {
	//net connection
	conn *net.TCPConn

	//dead time for read
	readchecktime int

	//remote host net address
	remoteAddr string

	//content for listening exit signal
	ctx context.Context

	//go goroutine
	wg *sync.WaitGroup

	//subclass redirect to subclass
	hm HandleMessage

	//for close signal
	status chan struct{}
}

//Loop ...
func (c *ClientSocket) loop() {

	for {

		if c.readchecktime > 0 {

			c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.readchecktime) * time.Second))
		}

		err := c.receivePackage()
		if nil != err {

			break
		}
	}

	c.Close()

	c.status <- struct{}{}

	c.wg.Done()
}

func (c *ClientSocket) Init() {

	c.readchecktime = 0
	c.wg = &sync.WaitGroup{}
	
	c.status = make(chan struct{})
}

func (c *ClientSocket) receivePackage() error {

	recvFixedData := func(length int) ([]byte, error) {

		var buf []byte
		recvSize := 0
		for {

			tmpBuf := make([]byte, length-recvSize)
			if nil == tmpBuf {

				log.Panic("out of memery")
			}

			cn, err := c.conn.Read(tmpBuf)

			if nil != err {

				log.Println("read err:", err)
				return nil, err
			}

			recvSize += cn
			buf = append(buf, tmpBuf[:cn]...)
			if recvSize >= length {

				break
			}
		}

		return buf, nil
	}

	mark := make([]byte, 1)
	if nil == mark {

		log.Panic("out of memery")
	}

	cn, err := c.conn.Read(mark)
	if err != nil {

		log.Println("read mark error:", err)
		return err
	}

	if 0 == cn || upackage.HeaderMarkUninit == upackage.HeaderType(mark[0]) {

		return fmt.Errorf("mark uninit")
	}

	headerlen := 0
	switch upackage.HeaderType(mark[0]) {

	case upackage.HeaderMarkID:

		headerlen = 6
	case upackage.HeaderMarkCmd:

		headerlen = 12
	}

	var headerBuf []byte
	if headerBuf, err = recvFixedData(headerlen); nil != err {

		return err
	}

	pkgHeader := upackage.NewBPackage(headerBuf)
	id := 0
	cmd := ""
	switch upackage.HeaderType(mark[0]) {

	case upackage.HeaderMarkID:

		id = int(pkgHeader.ReadPackageID())
	case upackage.HeaderMarkCmd:

		cmd = pkgHeader.ReadPackageCmd()
	}

	var bodyBuf []byte
	if bodyBuf, err = recvFixedData(int(pkgHeader.ReadDataLength())); nil != err {

		return err
	}

	return c.HandleMessage(int8(mark[0]), id, cmd, bodyBuf)
}

//HandleMessage handle message
func (c *ClientSocket) HandleMessage(mark int8, id int, cmd string, data []byte) error {

	var err error
	if upackage.HeaderMarkID == upackage.HeaderType(mark) {

		if nil != c.hm {

			err = c.hm.OnIDMessage(id, data)
		}

	} else {

		if nil != c.hm {

			err = c.hm.OnCmdMessage(cmd, data)
		}
	}

	return err
}

//Close close socket
func (c *ClientSocket) Close() {

	c.conn.Close()
}

//BroadcastIDMsg broadcast id by command to all client
func (c *ClientSocket) BroadcastIDMsg(id int, data []byte) {

}

//BroadcastCmdMsg broadcast msg by command to all client
func (c *ClientSocket) BroadcastCmdMsg(cmd string, data []byte) {

}

//Initialize init data
/*
param1: ctx
param2: conn
param3: readdeadtime
*/
func (c *ClientSocket) Initialize(ctx context.Context, signal chan int32, conn *net.TCPConn, id int32, checkTime int) {

	c.ctx = ctx
	c.conn = conn
	c.readchecktime = checkTime

	c.wg.Add(2)

	//listen ctx status
	go func(cs *ClientSocket, signal chan int32, id int32) {

		var err error
		select {

		case <-c.ctx.Done():

			//just for server exit,close read
			conn.CloseRead()

			err = ctx.Err()

			<-c.status
		case <-c.status:

			err = fmt.Errorf("an existing connection was forcibly closed by the remote host")
		}

		c.wg.Done()

		c.wg.Wait()

		signal <- id

		log.Printf("client[%s] exiting... reason: %s\n", c.conn.RemoteAddr().String(), err)

	}(c, signal, id)

	//for read data
	go c.loop()
}

//Redirect2Sub Runtime polymorphism
func (c *ClientSocket) Redirect2Sub(hm HandleMessage) {

	c.hm = hm
}

//RemoteAddress remote host address info
func (c *ClientSocket) RemoteAddress() string {

	return c.conn.RemoteAddr().String()
}

//SendIDMessage ...
func (c *ClientSocket) SendIDMessage(id int, data []byte) error {

	pkg := upackage.BPackage{}
	pkg.AddPackageID(int16(id))
	pkg.AddBytes(data)
	pkg.Done()

	wData := pkg.GetData()
	leftCount := len(wData)
	for leftCount > 0 {

		slen, err := c.conn.Write(wData)
		if nil != err {

			return err
		}

		wData = wData[slen:]
		leftCount -= slen
	}

	return nil
}

//SendCmdMessage ...
func (c *ClientSocket) SendCmdMessage(cmd string, data []byte) error {

	pkg := upackage.BPackage{}
	pkg.AddPackageCmd(cmd)
	pkg.AddBytes(data)
	pkg.Done()

	wData := pkg.GetData()
	leftCount := len(wData)
	for leftCount > 0 {

		slen, err := c.conn.Write(wData)
		if nil != err {

			return err
		}

		wData = wData[slen:]
		leftCount -= slen
	}

	return nil
}
