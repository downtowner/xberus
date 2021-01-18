package unet

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

const (
	shutdownInterval = 500 * time.Millisecond
)

//NewTCPServer create a tcpserver obj
func NewTCPServer() *TCPServer {

	p := TCPServer{}
	p.wg = sync.WaitGroup{}
	p.core = 1
	p.flag = 0
	return &p
}

//Connections handler connnections
type Connections interface {
	//handle connection process
	HandleConnections(context.Context, *net.TCPConn)

	//graceful exit for each connection
	Done()
}

//TCPServer server structure
type TCPServer struct {
	//go
	wg sync.WaitGroup

	//run in number of core
	core int

	//atomic for shutdown
	flag int32

	//listener
	listener *net.TCPListener
}

//Run start run service.
/**
network: tcp protocol
address: listen net address
Connections：handle process of a connection
**/
func (t *TCPServer) Run(network, address string, hc Connections) {

	if nil == hc {

		panic("Connections interface must be no-nil.")
	}

	ctx, cancle := context.WithCancel(context.TODO())

	addr, err := net.ResolveTCPAddr(network, address)
	if nil != err {

		log.Printf("Resolve error: %s\n", err)
		cancle()
		return
	}

	t.listener, err = net.ListenTCP(network, addr)
	if nil != err {

		log.Println(err)
		cancle()
		return
	}

	log.Println("server running on ", addr)

	for {

		conn, err := t.listener.AcceptTCP()
		if nil != err {

			log.Println("Listen err:", err.Error())
			break
		}

		log.Printf("A new client [%s] connected...\n", conn.RemoteAddr().String())

		hc.HandleConnections(ctx, conn)

	}

	cancle()

	hc.Done()

	log.Println("server closed...")
}

//Shutdown close service
func (t *TCPServer) closed(ctx context.Context) {

	//follow job may be faied in high concurrency scenarios
	if atomic.CompareAndSwapInt32(&t.flag, 0, 1) {

		t.listener.Close()

		select {
		case <-ctx.Done():

			//can do somthing
		}
	}

}

//Shutdown close service
func (t *TCPServer) Shutdown(timeout time.Duration, sig ...os.Signal) {

	go func(timeout time.Duration, sig ...os.Signal) {

		ch := make(chan os.Signal, 1)

		signal.Notify(ch, sig...)

		select {
		case m := <-ch:

			log.Printf("server recv [%s] signal. will shutdown...\n", m)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			t.closed(ctx)
		}

	}(timeout, sig...)
}
