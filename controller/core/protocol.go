package core

const (
	PktTypePing           = 0x02
	PktTypeData           = 0x01
	PktTypeRelayCfgReq    = 0xAA
	PktTypeRelayCfgResp   = 0xAB
	PktTypeRelayCfgReqLen = 7
)

type Packet interface {
	Addr() []byte // Returns board address.
}

type PingPacket struct {
	addr []byte
}

func (f *PingPacket) Addr() []byte {
	return f.addr
}

/********** utility functions for protocol *******************/

func okPktTypeRelayCfgReq(buffer []byte) bool {
	if buffer[0] != PktTypeRelayCfgReq || len(buffer) < PktTypeRelayCfgReqLen {
		return false
	}
	return true
}

func getHWAddr(buffer []byte) []byte {
	return buffer[1:7]
}

func makePktTypeRelayCfgResp(r *Relay) []byte {
	pkt := []byte{}
	pkt = append(pkt, PktTypeRelayCfgResp)
	pkt = append(pkt, r.PAddr0...)
	pkt = append(pkt, r.PAddr1...)
	pkt = append(pkt, r.PAddr2...)
	pkt = append(pkt, r.PAddr3...)
	pkt = append(pkt, r.PAddr4...)
	pkt = append(pkt, r.PAddr5...)
	pkt = append(pkt, r.Channel)
	pkt = append(pkt, r.Addr...)

	return pkt
}
