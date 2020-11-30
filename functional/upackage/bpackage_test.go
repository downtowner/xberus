package upackage

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestPackage(t *testing.T) {

	//build a package
	idpkg := NewBPackage()
	idpkg.AddPackageID(1)

	idpkg.AddInt8(8)
	idpkg.AddInt16(16)
	idpkg.AddInt32(32)
	idpkg.AddInt64(64)

	idpkg.AddUint8(18)
	idpkg.AddUint16(116)
	idpkg.AddUint32(132)
	idpkg.AddUint64(164)

	idpkg.AddBool(false)

	idpkg.AddFloat32(1.2345678, 7)
	idpkg.AddFloat64(1.234567890, 9)

	idpkg.AddString("not with length")
	idpkg.AddStringL("with length", 20)

	idpkg.AddInt64(6464)
	idpkg.Done()

	log.Println(string(idpkg.GetData()))

	//create a read package
	rpkg := NewBPackage()
	rpkg.AddBytes(idpkg.GetData())

	log.Println("mark: ", rpkg.ReadHeaderMark())
	log.Println("id: ", rpkg.ReadPackageID())
	log.Println("size: ", rpkg.ReadDataLength())

	log.Println("", rpkg.ReadInt8())
	log.Println("", rpkg.ReadInt16())
	log.Println("", rpkg.ReadInt32())
	log.Println("", rpkg.ReadInt64())

	log.Println("", rpkg.ReadUint8())
	log.Println("", rpkg.ReadUint16())
	log.Println("", rpkg.ReadUint32())
	log.Println("", rpkg.ReadUint64())

	log.Println("", rpkg.ReadBool())

	log.Println(rpkg.ReadFloat32())
	log.Println(rpkg.ReadFloat64())

	log.Println("", rpkg.ReadString())
	log.Println("", rpkg.ReadStringL(20))

	log.Println("", rpkg.ReadInt64())

	testContext()
}

func testContext() {
	//tree()
	//testcontext1()
	testcontext2()
}

func tree() {
	ctx1 := context.Background()
	ctx2, _ := context.WithCancel(ctx1)

	go func(ctx context.Context) {
		log.Println("ctx2")
		select {

		case <-ctx.Done():

			log.Println("ctx2 exit")
		}
	}(ctx2)

	ctx3, _ := context.WithTimeout(ctx2, time.Second*5)
	go func(ctx context.Context) {
		log.Println("ctx3")
		select {

		case <-ctx.Done():

			log.Println("ctx3 exit")
		}
	}(ctx3)

	ctx4, _ := context.WithTimeout(ctx3, time.Second*3)
	go func(ctx context.Context) {
		log.Println("ctx4")
		select {

		case <-ctx.Done():

			log.Println("ctx4 exit")
		}
	}(ctx4)

	ctx5, cancle := context.WithTimeout(ctx3, time.Second*6)
	go func(ctx context.Context) {
		log.Println("ctx5")
		select {

		case <-ctx.Done():

			log.Println("ctx5 exit")
		}
	}(ctx5)

	cancle()

	ctx6 := context.WithValue(ctx5, "userID", 12)
	go func(ctx context.Context) {
		log.Println("ctx6")
		select {

		case <-ctx.Done():

			log.Println("ctx6 exit", ctx.Err())
		}
	}(ctx6)

	time.Sleep(time.Second * 10)

}

func testcontext2() {

	var wg sync.WaitGroup
	app := testApp{}

	wg.Add(1)

	go func() {

		ch := make(chan os.Signal, 1)

		signal.Notify(ch,
			// kill -SIGINT XXXX 或 Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			os.Kill,
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-ch:

			println("shutdown...")

			timeout := 5 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			app.Shutdown(ctx)
			wg.Done()
			return
		}

	}()

	app.Run()
	wg.Wait()

	log.Println("wtf")
}

type testApp struct {
	ctx    context.Context
	cancel context.CancelFunc
	opts   int64
}

func (t *testApp) Run() {

	t.ctx, t.cancel = context.WithCancel(context.TODO())

	var wg sync.WaitGroup

	if nil == t.ctx {
		log.Println("no context")
	} else {

	}

	handle := func(f *int64) {

		t := time.Tick(500 * time.Millisecond)
		for k := range t {

			log.Println("k", k)

			if atomic.LoadInt64(f) > 0 {

				break
			}
		}

		wg.Done()
	}

	wg.Add(1)
	go handle(&t.opts)

	wg.Wait()
}

func (t *testApp) Shutdown(ctx context.Context) {
	t.ctx = ctx

	select {
	case <-t.ctx.Done():
		log.Println("exit reason:", t.ctx.Err())
		atomic.AddInt64(&t.opts, 1)
	}
}
