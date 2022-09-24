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
	Motion bool
}

type Door struct {
	pktInfo
	Door bool
}

type Temp struct {
	pktInfo
	Temp     float32
	Humidity float32
}
type Light struct {
	pktInfo
	Light float32
}

type Volt struct {
	pktInfo
	Volt float32
}
type Error struct {
	pktInfo
	ErrCode byte
}
