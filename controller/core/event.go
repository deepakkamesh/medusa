package core

type Event interface {
	Addr() []byte
	PAddr() []byte
	HWAddr() []byte
}

type pktInfo struct {
	addr   []byte // Board address.
	paddr  []byte // Pipe address.
	hwaddr []byte // Hardware address.
}

func (f pktInfo) Addr() []byte {
	return f.addr
}

func (f pktInfo) PAddr() []byte {
	return f.paddr
}

func (f pktInfo) HWAddr() []byte {
	return f.hwaddr
}

// The following implement the Event Interface.
type Ping struct {
	pktInfo
}

type Motion struct {
	pktInfo
}

type Temp struct {
	pktInfo
	Temp     byte
	Humidity byte
}

type Error struct {
	pktInfo
	ErrCode byte
}
