package ufreequeue

import (
	"sync/atomic"
)

const (
	//MaxEleCount set size of queue
	MaxEleCount = 512
)

//NewAdQueue create obj
func NewAdQueue() *AdvancedQueue {

	p := AdvancedQueue{}
	p.data = make([]interface{}, MaxEleCount)
	return &p
}

//AdvancedQueue for read and write msg
type AdvancedQueue struct {
	//data feild
	data []interface{}

	//windex for wirte data
	windex int32

	//rindex for read data
	rindex int32
}

//Push push a msg to queue
func (a *AdvancedQueue) Push(msg interface{}) {

	//add judge logic for index of reading and writing

	for {

		tmp := atomic.LoadInt32(&a.windex)
		if atomic.CompareAndSwapInt32(&a.windex, tmp, (tmp+1)%MaxEleCount) {

			//index must be different
			a.data[a.windex] = msg
			break
		}

	}

}

//Get remove a msg from queue
func (a *AdvancedQueue) Get() interface{} {

	return nil
}
