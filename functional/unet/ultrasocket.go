package unet

import (
	"fmt"
	"log"
	"net"
	"sync"
	"xframe/functional/upackage"
)

//UltralSocket ...
type UltralSocket struct {
	//Synchronous or asynchronous
	sync bool

	//connection
	conn *net.TCPConn

	//server address
	sAdress string

	//handle message
	handle func(mark int8, id int32, cmd string, data []byte) error

	//goroutine group
	wg *sync.WaitGroup
}

//NewUSocket ...
func NewUSocket(addr string, syn bool) *UltralSocket {

	return &UltralSocket{sAdress: addr, sync: syn, wg: &sync.WaitGroup{}}
}

//Connect params future-use
func (u *UltralSocket) Connect(params ...interface{}) error {

	if !u.sync {

		var ok bool
		if u.handle, ok = params[0].(func(int8, int32, string, []byte) error); !ok {

			return fmt.Errorf("Asynchronous connection must specify the first parameter(func(int8, int32, string, []byte) error)")
		}
	}

	if nil == u.handle {

		return fmt.Errorf("Please mount message processing")
	}

	if nil != u.conn {

		u.Close()
	}

	addr, err := net.ResolveTCPAddr("tcp", u.sAdress)
	if err != nil {
		return err
	}

	cn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	cn.SetKeepAlive(true)

	u.conn = cn

	if !u.sync {

		u.wg.Add(1)
		go u.loop()
	}

	return nil
}

//SendIDMessage ...
func (u *UltralSocket) SendIDMessage(id int, data []byte) error {

	pkg := upackage.BPackage{}
	pkg.AddPackageID(int16(id))
	pkg.AddBytes(data)
	pkg.Done()

	wData := pkg.GetData()
	leftCount := len(wData)
	for leftCount > 0 {

		slen, err := u.conn.Write(wData)
		if nil != err {

			return err
		}

		wData = wData[slen:]
		leftCount -= slen
	}

	return nil
}

//SendCmdMessage ...
func (u *UltralSocket) SendCmdMessage(cmd string, data []byte) error {

	pkg := upackage.BPackage{}
	pkg.AddPackageCmd(cmd)
	pkg.AddBytes(data)
	pkg.Done()

	wData := pkg.GetData()
	leftCount := len(wData)
	for leftCount > 0 {

		slen, err := u.conn.Write(wData)
		if nil != err {

			return err
		}

		wData = wData[slen:]
		leftCount -= slen
	}

	return nil
}

//SendAndRecvIDMsg ...
func (u *UltralSocket) SendAndRecvIDMsg(id int, data []byte, handle func(int, []byte) error) error {

	if !u.sync {

		return fmt.Errorf("Asynchronous process does not support synchronous interface")
	}

	if err := u.SendIDMessage(id, data); nil != err {

		return err
	}
	
	_, id, _, body, err := u.receiveOnce()
	if nil != err {

		return err
	}

	if nil == handle {

		return fmt.Errorf("warning: Lack of necessary processing")
	}

	return handle(id, body)
}

//SendAndRecvCmdMsg ...
func (u *UltralSocket) SendAndRecvCmdMsg(cmd string, data []byte, handle func(string, []byte) error) error {

	if !u.sync {

		return fmt.Errorf("Asynchronous process does not support synchronous interface")
	}

	if err := u.SendCmdMessage(cmd, data); nil != err {

		return err
	}

	_, _, cmd, body, err := u.receiveOnce()
	if nil != err {

		return err
	}

	if nil == handle {

		return fmt.Errorf("warning: Lack of necessary processing")
	}

	return handle(cmd, body)
}

func (u *UltralSocket) loop() error {

	for {

		defer u.wg.Done()

		mark, id, cmd, body, err := u.receiveOnce()
		if nil != err {

			return err
		}

		if err := u.handle(mark, int32(id), cmd, body); nil != err {

			return err
		}
	}
}

func (u *UltralSocket) receiveOnce() (int8, int, string, []byte, error) {

	recvFixedData := func(buf []byte, length int) error {

		recvSize := 0
		for {

			tmpBuf := make([]byte, length-recvSize)
			if nil == tmpBuf {

				log.Panic("out of memery")
			}

			cn, err := u.conn.Read(tmpBuf)

			if nil != err {

				log.Println("read err:", err)
				return err
			}

			recvSize += cn
			buf = append(buf, buf...)
			if recvSize >= length {

				break
			}
		}

		return nil
	}

	mark := make([]byte, 1)
	if nil == mark {

		log.Panic("out of memery")
	}

	cn, err := u.conn.Read(mark)
	if err != nil {

		log.Println("read mark error:", err)
		return 0, 0, "", nil, err
	}

	if 0 == cn || upackage.HeaderMarkUninit == upackage.HeaderType(mark[0]) {

		return 0, 0, "", nil, fmt.Errorf("mark uninit")
	}

	headerlen := 0
	switch upackage.HeaderType(mark[0]) {

	case upackage.HeaderMarkID:

		headerlen = 6
	case upackage.HeaderMarkCmd:

		headerlen = 12
	}

	headerBuf := []byte{}
	if err := recvFixedData(headerBuf, headerlen); nil != err {

		return 0, 0, "", nil, err
	}

	pkgHeader := upackage.NewBPackage(headerBuf)
	id := 0
	cmd := ""
	switch upackage.HeaderType(mark[0]) {

	case upackage.HeaderMarkID:

		id = int(pkgHeader.ReadPackageID())
	case upackage.HeaderMarkCmd:

		cmd = pkgHeader.ReadPackageCmd()
	}

	bodyBuf := []byte{}
	if err := recvFixedData(bodyBuf, int(pkgHeader.ReadDataLength())); nil != err {

		return 0, 0, "", nil, err
	}

	return int8(mark[0]), id, cmd, bodyBuf, nil
}

//Close ...
func (u *UltralSocket) Close() {
	u.conn.Close()
	u.wg.Wait()
}
