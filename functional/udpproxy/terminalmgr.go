package udpproxy

type terminalMgr struct {
	toTerminals <-chan []byte
	toOrigin    chan<- []byte
	terminals   []string
	terminalObj []*terminaler
}

func NewTerminalerMgr(netaddr []string, r <-chan []byte, w chan<- []byte) *terminalMgr {

	p := terminalMgr{}
	p.terminals = netaddr
	p.toTerminals = r
	p.toOrigin = w

	return &p
}

func (t *terminalMgr) Run() {

	for _, v := range t.terminals {

		ter := NewTerminaler(v, t.toOrigin)
		ter.Run()

		t.terminalObj = append(t.terminalObj, ter)
	}

	t.listenMsg()
}

func (t *terminalMgr) listenMsg() {
	for {

		select {

		case data := <-t.toTerminals:

			for _, v := range t.terminalObj {

				v.SendMsg(data)
			}
		}
	}
}
