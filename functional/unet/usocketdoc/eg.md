# 使用实例

```go
package main

import (
	"log"
	"sync"
	"time"

	"git.vnnox.net/ncp/xframe/functional/unet"
)

func testServer() {

	m := &sync.WaitGroup{}

	u := unet.NewUSocket("127.0.0.1:10086", false)
	u.Connect(func(mark int8, id int32, cmd string, data []byte) error {

		log.Println("read: =================================")
		log.Println("mark:", mark)
		log.Println("id:", id)
		log.Println("cmd:", cmd)
		log.Println("data:", data)
		return nil
	})

	m.Add(1)
	go func() {

		for {
			time.Sleep(time.Second * 2)
			u.SendIDMessage(1, []byte("hello,world!"))
		}

		m.Done()
	}()

	m.Wait()
}

func main() {

	testServer()
}

```

