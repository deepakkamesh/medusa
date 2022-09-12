package core

import (
	"encoding/binary"
	"fmt"
	"math"
)

const (
	// Packet Types.
	PktTypePing      = 0x02
	PktTypeData      = 0x01
	PktTypeActionReq = 0x10
	PktTypeConfig    = 0x05

	PktTypeRelayCfgReq    = 0xAA
	PktTypeRelayCfgResp   = 0xAB
	PktTypeRelayError     = 0xAC
	PktTypeRelayBoardData = 0xAD // Relay.

	// Actions.
	ActionMotion      = 0x01
	ActionTemp        = 0x02
	ActionLight       = 0x03
	ActionDoor        = 0x04
	ActionVolt        = 0x05
	ActionBuzzer      = 0x10
	ActionGarage      = 0x11
	ActionRelay       = 0x12
	ActionLED         = 0x13
	ActionReset       = 0x14
	ActionTest        = 0x16
	ActionFlushTXFIFO = 0x17 // Flush TX_FIFO (relay).

	// Error Codes.
	ErrNA             = 0x00
	ErrRadioInit      = 0x02
	ErrLoadAck        = 0x03
	ErrActionNotImpl  = 0x04
	ErrUnknownPktType = 0x05
	ErrPipeAddr404    = 0x06

	PktTypeRelayCfgReqLen    = 7
	PktTypeRelayBoardDataLen = 10 // min len of board data pkt.
)

// Default config pipe address.
var defPipeAdress []byte = []byte{0x68, 0x65, 0x6C, 0x6C, 0x6F}

func okPktTypeRelayCfgReq(buffer []byte) bool {
	if buffer[0] != PktTypeRelayCfgReq || len(buffer) < PktTypeRelayCfgReqLen {
		return false
	}
	return true
}

// getHWAddr returns the hardware address from the relay config request packet.
func getHWAddr(buffer []byte) []byte {
	return buffer[1:7]
}

// makePktTypeRelayCfgResp creates a relay config response packet.
// if defPAddr is true, sets PAddr0 to defPipeAddress ("hello").
func makePktTypeRelayCfgResp(r *Relay, defPaddr bool) []byte {
	pkt := []byte{}
	pkt = append(pkt, PktTypeRelayCfgResp)
	if defPaddr {
		pkt = append(pkt, defPipeAdress...)
	} else {
		pkt = append(pkt, r.PAddr0...)
	}
	pkt = append(pkt, r.PAddr1...)
	pkt = append(pkt, r.PAddr2[0])
	pkt = append(pkt, r.PAddr3[0])
	pkt = append(pkt, r.PAddr4[0])
	pkt = append(pkt, r.PAddr5[0])
	pkt = append(pkt, r.PAddr6[0])
	pkt = append(pkt, r.Channel)
	pkt = append(pkt, r.Addr...)
	return pkt
}

func makePktTypeConfig(addr []byte, paddr []byte, b *Board) []byte {
	pkt := []byte{}
	pkt = append(pkt, PktTypeRelayBoardData)
	pkt = append(pkt, paddr...)
	pkt = append(pkt, PktTypeConfig)
	pkt = append(pkt, addr...)
	pkt = append(pkt, b.ARD)
	pkt = append(pkt, b.PingInt)
	pkt = append(pkt, b.PAddr...)
	pkt = append(pkt, b.Addr...)

	return pkt
}

func makePktTypeActionReq(actionID byte, addr []byte, paddr []byte, data []byte) []byte {
	pkt := []byte{}
	pkt = append(pkt, PktTypeRelayBoardData)
	pkt = append(pkt, paddr...)
	pkt = append(pkt, PktTypeActionReq)
	pkt = append(pkt, addr...)
	pkt = append(pkt, actionID)
	pkt = append(pkt, data...)

	return pkt
}

// translatePacket converts the byte packet into an event.
// TODO make more readable and add some length validation.
func translatePacket(pkt []byte, hwaddr []byte) (Event, error) {

	p := pktInfo{
		hwaddr: hwaddr,
	}

	// Handle if relay error.
	if pkt[0] == PktTypeRelayError {
		return translateErrorPacket(p, pkt[1])
	}

	// Validate to ensure packet is valid.
	// TODO: Implement checksum.
	if len(pkt) < PktTypeRelayBoardDataLen {
		return nil, fmt.Errorf("bad packet len")
	}

	// Handle if relay board data.
	if pkt[0] != PktTypeRelayBoardData {
		return nil, fmt.Errorf("unknown relay packet type %v", pkt[0])
	}

	p.paddr = pkt[1:6]
	p.addr = pkt[7:10]
	switch pkt[6] {
	case PktTypePing:
		return translatePingPacket(p)

	case PktTypeData:
		errCode := pkt[11]
		action := pkt[10]
		data := pkt[12:]

		if errCode != ErrNA {
			return translateErrorPacket(p, errCode)
		}
		return translateActionPacket(p, action, data)

	default:
		return nil, fmt.Errorf("unknown board packet type %v", pkt[6])
	}
}

func translatePingPacket(p pktInfo) (Event, error) {
	return Ping{
		pktInfo: p,
	}, nil
}

func translateErrorPacket(p pktInfo, errCode byte) (Event, error) {
	return Error{
		pktInfo: p,
		ErrCode: errCode,
	}, nil
}

func translateActionPacket(p pktInfo, action byte, data []byte) (Event, error) {

	switch action {

	case ActionTemp:
		t, h := readTemp(data)
		return Temp{
			pktInfo:  p,
			Temp:     t,
			Humidity: h,
		}, nil

	case ActionMotion:
		motion := false
		if data[0] == 1 {
			motion = true
		}
		return Motion{
			pktInfo: p,
			Motion:  motion,
		}, nil

	case ActionVolt:
		var x uint
		x = uint(data[1])
		x = x<<8 | uint(data[0])
		return Volt{
			pktInfo: p,
			Volt:    (2.048 * 1023) / float32(x),
		}, nil
	default:
		return nil, fmt.Errorf("unknown action %v", action)
	}
}

func readTemp(data []byte) (temp float32, humidity float32) {

	_temp := []byte{data[0], data[1], data[2], data[3]}
	_humidity := []byte{data[4], data[5], data[6], data[7]}

	// Convert to 'F and return.
	return float32frombytes(_temp)*1.8 + 32, float32frombytes(_humidity)
}

func float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
