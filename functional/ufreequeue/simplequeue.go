package ufreequeue

import (
	"time"
)

//NewSimQueue ...
func NewSimQueue() *SimpleQueue {

	p := SimpleQueue{}
	p.msg = make(chan interface{}, MaxEleCount)

	return &p
}

//SimpleQueue ...
type SimpleQueue struct {
	msg chan interface{}
}

//Push push a msg to queue
func (s *SimpleQueue) Push(msg interface{}, timeout time.Duration) bool {

	select {

	case s.msg <- msg:

		return true
	case <-time.After(timeout):

		return false
	}
}

//Get get a msg from queue
func (s *SimpleQueue) Get(timeout time.Duration) interface{} {

	select {

	case data := <-s.msg:

		return data
	case <-time.After(timeout):

		return nil
	}
}
