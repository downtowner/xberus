package unet

import (
	"log"
	"sync"
	"time"
)

func testServer() {

	m := &sync.WaitGroup{}

	u := NewUSocket("127.0.0.1:9999", false)
	err := u.Connect(func(mark int8, id int32, cmd string, data []byte) error {

		log.Println("read: =================================")
		log.Println("mark:", mark)
		log.Println("id:", id)
		log.Println("cmd:", cmd)
		log.Println("data:", data)
		return nil
	})

	if nil != err {
		log.Println("connect err:", err)
		return
	}

	m.Add(1)
	go func() {

		for {
			log.Println("send:")
			time.Sleep(time.Second * 2)
			u.SendIDMessage(1, []byte("hello,world!"))
		}

		m.Done()
	}()

	m.Wait()
}

//main xxxxx
func main() {

	log.Println("start")

	for i := 0; i < 100; i++ {

		//xxxx
		go testServer()
		time.Sleep(time.Microsecond * 10)
	}

	testServer()

	log.Println("endl")
}
