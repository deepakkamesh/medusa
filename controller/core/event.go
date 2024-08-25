package core

type Event interface {
	Addr() []byte
	PAddr() []byte
	HWAddr() []byte
}

type PktInfo struct {
	BoardAddr    []byte // Board address.
	PipeAddr     []byte // Pipe address.
	HardwareAddr []byte // Hardware address.
}

func (f PktInfo) Addr() []byte {
	return f.BoardAddr
}

func (f PktInfo) PAddr() []byte {
	return f.PipeAddr
}

func (f PktInfo) HWAddr() []byte {
	return f.HardwareAddr
}

// The following implement the Event Interface.
type Ping struct {
	PktInfo
}

type Motion struct {
	PktInfo
	Motion bool
}

type Door struct {
	PktInfo
	Door bool
}

type Temp struct {
	PktInfo
	Temp     float32
	Humidity float32
}
type Light struct {
	PktInfo
	Light float32
}
type Gas struct {
	PktInfo
	Gas float32
}
type Pressure struct {
	PktInfo
	Pressure float32
}
type Altitude struct {
	PktInfo
	Altitude float32
}
type Volt struct {
	PktInfo
	Volt float32
}
type Error struct {
	PktInfo
	ErrCode byte
}
