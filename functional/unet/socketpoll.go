package unet

import "sync"

//USocketPool ...
type USocketPool struct {
	serverAddr     string
	maxConnections int32
	idleSockets    []*UltralSocket
	busySockets    []*UltralSocket
	mutex          sync.Mutex
}

//NewPool ...
func NewPool(addr string, maxConn int32) *USocketPool {

	p := USocketPool{}
	p.maxConnections = maxConn
	p.serverAddr = addr
	p.idleSockets = []*UltralSocket{}
	p.busySockets = []*UltralSocket{}
	p.mutex = sync.Mutex{}

	return &p
}

//GetSocket ...
func (u *USocketPool) GetSocket() *UltralSocket {

	u.mutex.Lock()
	defer u.mutex.Unlock()

	if len(u.idleSockets) > 0 {

		c := u.idleSockets[0]
		u.idleSockets = u.idleSockets[1:]
		u.busySockets = append(u.busySockets, c)

		return c
	}

	if int32(len(u.idleSockets)+len(u.busySockets)) < u.maxConnections {

		c := NewUSocket(u.serverAddr, false)
		u.busySockets = append(u.busySockets, c)

		return c
	}

	return nil
}

//FreeSocket ...
func (u *USocketPool) FreeSocket(socket *UltralSocket) {

	u.mutex.Lock()
	defer u.mutex.Unlock()

	for i := 0; i < len(u.busySockets); i++ {

		if u.busySockets[i] == socket {

			u.idleSockets = append(u.idleSockets, socket)
			u.busySockets = append(u.busySockets[0:i], u.busySockets[i+1:]...)

		}
	}

}
