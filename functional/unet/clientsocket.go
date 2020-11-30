package unet

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"xframe/functional/upackage"
)

/**
author: irarus
data: 2020-10-27

note:

**/

//NewClientSocket create a clientsocket obj, c tcp connection. params Optional configuration parameters.
/**

The meaning of the parameters is as follows:
readchecktime check socket cycle of read data, closed when it time out, 0,uncheck, unit:second

**/
func NewClientSocket() Client {

	p := ClientSocket{}
	p.readchecktime = 0
	p.wg = &sync.WaitGroup{}

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
}

//Loop ...
func (c *ClientSocket) loop() {

	for {

		if c.readchecktime > 0 {

			c.conn.SetReadDeadline(time.Now().Add(time.Duration(c.readchecktime) * time.Second))
		}

		if err := c.receivePackage(); nil != err {

			break
		}
	}

	c.Close()

	c.wg.Done()
}

func (c *ClientSocket) receivePackage() error {

	recvFixedData := func(buf []byte, length int) error {

		recvSize := 0
		for {

			tmpBuf := make([]byte, length-recvSize)
			if nil == tmpBuf {

				log.Panic("out of memery")
			}

			cn, err := c.conn.Read(tmpBuf)

			if nil != err {

				log.Println("read err:", err)
				return err
			}

			recvSize += cn
			buf = append(buf, buf...)
			if recvSize >= length {

				break
			}
		}

		return nil
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

	headerBuf := []byte{}
	if err := recvFixedData(headerBuf, headerlen); nil != err {

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

	bodyBuf := []byte{}
	if err := recvFixedData(bodyBuf, int(pkgHeader.ReadDataLength())); nil != err {

		return err
	}

	return c.HandleMessage(int8(mark[0]), id, cmd, bodyBuf)
}

//HandleMessage handle message
func (c *ClientSocket) HandleMessage(mark int8, id int, cmd string, data []byte) error {

	var err error
	if upackage.HeaderMarkID == upackage.HeaderType(mark) {

		err = c.OnIDMessage(id, data)
	} else {

		err = c.OnCmdMessage(cmd, data)
	}

	return err
}

//OnCmdMessage handle cmd message
func (c *ClientSocket) OnCmdMessage(cmd string, data []byte) error {

	return nil
}

//OnIDMessage handle id message
func (c *ClientSocket) OnIDMessage(id int, data []byte) error {

	return nil
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
func (c *ClientSocket) Initialize(ctx context.Context, conn *net.TCPConn, checkTime int) {

	c.ctx = ctx
	c.conn = conn
	c.readchecktime = checkTime

	c.wg.Add(2)

	//listen ctx status
	go func(ctx context.Context, c *net.TCPConn, wg *sync.WaitGroup) {

		select {

		case <-ctx.Done():

			//just for server exit,close read
			c.CloseRead()

			wg.Done()
		}
	}(c.ctx, c.conn, c.wg)

	//for read data
	go c.loop()

	c.Close()
}

//Done safe exit service
func (c *ClientSocket) Done() {

	c.wg.Wait()

	log.Printf("client[%s] exiting... reason: %s\n", c.remoteAddr, c.ctx.Err())
}
