package unet

import (
	"context"
	"log"
	"net"
	"reflect"
	"sync"

	"git.vnnox.net/ncp/xframe/functional/upackage"
)

//Client base interface
type Client interface {
	Initialize(ctx context.Context, signal chan int32, conn *net.TCPConn, id int32, checkTime int)
	Redirect2Sub(hm HandleMessage)
	Close()
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
	p.Init()
	p.of = of

	return p
}

//ConnectionMgr ...
type ConnectionMgr struct {

	//Factory for creating objects
	of ObjFactory

	//Establish a channel between client and mgr
	broadcast chan *upackage.BPackage

	//manage client list
	cList map[int32]Client

	//safe lock for ctx/clist
	lock *sync.Mutex

	//use ctx for listening server status
	ctx context.Context

	//client exit signal
	signal chan int32

	//client id seed
	seed int32

	wg *sync.WaitGroup
}

//SetConnecter set your connecters obj factory
func (c *ConnectionMgr) SetConnecter(of ObjFactory) {

	c.of = of
}

//Init if you wanna extend `ConnectionMgr` class .you must call this method
func (c *ConnectionMgr) Init() {

	c.broadcast = make(chan *upackage.BPackage, 100)
	c.lock = &sync.Mutex{}
	c.enableBroadcast()
	c.cList = make(map[int32]Client)
	c.signal = make(chan int32)
	c.wg = &sync.WaitGroup{}
}

//HandleConnections ...
func (c *ConnectionMgr) HandleConnections(ctx context.Context, conn *net.TCPConn) {

	if nil == c.ctx {

		c.lock.Lock()
		if nil == c.ctx {

			c.ctx = ctx
			go c.help()
		}
		c.lock.Unlock()
	}

	client := c.of()

	if reflect.Ptr != reflect.TypeOf(client).Kind() {

		log.Panicf("ObjFactory func return value must be ptr")
	}

	if nil == client {

		log.Panicf("Value of clientFactory.NewSocketClinet must be non-nil")
	}

	c.seed++

	c.lock.Lock()
	c.cList[c.seed] = client
	log.Println("online users: ", len(c.cList))
	c.lock.Unlock()

	if sub, ok := client.(HandleMessage); ok {

		client.Redirect2Sub(sub)
	}
	//add the failed log.

	client.Initialize(ctx, c.signal, conn, c.seed, 0)
}

//Done safe exit service
func (c *ConnectionMgr) Done() {
	c.wg.Wait()
	close(c.broadcast)
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

func (c *ConnectionMgr) help() {

	c.wg.Add(1)
	go func() {

		willExit := false
		for {

			exit := false
			select {

			case id := <-c.signal:

				c.removeClient(id)
				if willExit && 0 == len(c.cList) {

					exit = true
				}
			case <-c.ctx.Done():

				if 0 == len(c.cList) {

					exit = true
				} else {

					willExit = true
					c.closeAllClient()
				}
			}

			if exit {
				break
			}
		}
		c.wg.Done()
	}()
}

func (c *ConnectionMgr) removeClient(id int32) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.cList, id)
}

func (c *ConnectionMgr) closeAllClient() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, v := range c.cList {
		v.Close()
	}
}
