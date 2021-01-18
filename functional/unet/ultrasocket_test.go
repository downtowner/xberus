package unet

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestSocket(t *testing.T) {

	m := &sync.WaitGroup{}

	u := NewUSocket("127.0.0.1:10086", false)
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
			time.Sleep(time.Second * 2)
			u.SendIDMessage(1, []byte("hello,world!"))
		}

		m.Done()
	}()

	m.Wait()
}
