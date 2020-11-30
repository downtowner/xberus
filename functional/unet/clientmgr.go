package unet

import (
	"context"
	"log"
	"net"
	"reflect"
	"sync"
	"xframe/functional/upackage"
)

//Client base interface
type Client interface {
	Initialize(ctx context.Context, conn *net.TCPConn, checkTime int)
	Done()
}

//WClient io write interface
type WClient interface {
	SendMessage([]byte) error
}

//ObjFactory object factory
type ObjFactory func() Client

//NewConnMgr ...
func NewConnMgr(of ObjFactory) *ConnectionMgr {

	if nil == of {

		log.Panicf("ObjFactory must be non-empty")
	}

	p := &ConnectionMgr{}
	p.broadcast = make(chan *upackage.BPackage, 100)
	p.lock = &sync.Mutex{}

	p.enableBroadcast()

	return p
}

//ConnectionMgr ...
type ConnectionMgr struct {

	//Factory for creating objects
	of ObjFactory

	//Establish a channel between client and mgr
	broadcast chan *upackage.BPackage

	//manage client list
	cList []Client

	//safe lock for ctx/clist
	lock *sync.Mutex

	//use ctx for listening server status
	ctx context.Context
}

//HandleConnections ...
func (c *ConnectionMgr) HandleConnections(ctx context.Context, conn *net.TCPConn) {

	if nil == c.ctx {

		c.lock.Lock()
		defer c.lock.Unlock()
		if nil == c.ctx {

			c.ctx = ctx
		}
	}

	client := c.of()

	if reflect.Ptr != reflect.TypeOf(client).Kind() {

		log.Panicf("ObjFactory func return value must be ptr")
	}

	if nil == client {

		log.Panicf("Value of clientFactory.NewSocketClinet must be non-nil")
	}

	c.lock.Lock()
	c.cList = append(c.cList, client)
	c.lock.Unlock()

	client.Initialize(ctx, conn, 0)
}

//Done safe exit service
func (c *ConnectionMgr) Done() {

	for _, v := range c.cList {

		v.Done()
	}
}

func (c *ConnectionMgr) enableBroadcast() {

	go func(mgr *ConnectionMgr) {

		for {

			select {

			case msg := <-c.broadcast:

				for _, v := range mgr.cList {

					if sed, ok := v.(WClient); ok {

						sed.SendMessage(msg.GetData())
					}
				}
			}
		}
	}(c)
}
