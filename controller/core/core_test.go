package core

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"
	"time"
)

var core *Core
var hostPort string = "127.0.0.1:3344"

func init() {
	var err error
	core, err = NewCore(hostPort, "core.cfg.test.json")
	if err != nil {
		fmt.Printf("Error initializing core:%v", err)
		os.Exit(1)
	}
	core.StartCore()
	// Add some delay for the service to listen & get ready.
	time.Sleep(100 * time.Millisecond)
}

func TestRelayConfig(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", hostPort)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		t.Fatalf("Could not dial server:%v", err)
	}

	defer conn.Close()

	data := []struct {
		req  []byte
		resp []byte
	}{
		{[]byte{PktTypeRelayCfgReq, 0xa1, 0xb1, 0xc1, 0xd1, 0xe1, 0xf1},
			[]byte{PktTypeRelayCfgResp, 0x1, 0xc1, 0xd1, 0xe1, 0xf1, 0x2, 0xc1, 0xd1, 0xe1, 0xf1, 0x3, 0x4, 0x5, 0x6, 0x7, 0x73, 0xd1, 0xe1, 0xf1}},
		{[]byte{PktTypeRelayCfgReq, 0xa2, 0xb2, 0xc2, 0xd2, 0xe2, 0xf2},
			[]byte{PktTypeRelayCfgResp, 0x1, 0xc2, 0xd2, 0xe2, 0xf2, 0x2, 0xc2, 0xd2, 0xe2, 0xf2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x73, 0xd2, 0xe2, 0xf2}},
	}

	// Send config request to server.
	for _, v := range data {
		_, err = conn.Write(v.req)

		if err != nil {
			t.Error(err)
		}
		//  Receive config response.
		buffer := make([]byte, 255)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("%v", err)
		}
		buffer = buffer[:n]
		if bytes.Compare(buffer, v.resp) != 0 {
			t.Errorf("want %v got %v", v.resp, buffer)
		}
	}
}

// TestAction tests if the action function results in the correct relay data packet.
func TestAction(t *testing.T) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", hostPort)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		t.Fatalf("Failed to dial server:%v", err)
	}
	// Delay needed to stabilize listen on server.
	time.Sleep(100 * time.Millisecond)

	// Temp.
	if err := core.Temp([]byte{0x1, 0x1, 0x1}); err != nil {
		t.Fatalf("Failed to call temp: %v", err)
	}
	if err := readandCompare(conn, []byte{PktTypeRelayBoardData, 0x1, 0xc1, 0xd1, 0xe1, 0xf1, PktTypeActionReq, 0x1, 0x1, 0x1, ActionTemp}); err != nil {
		t.Error(err)
	}

	// LedOn.
	if err := core.LEDOn([]byte{0xd1, 0xe1, 0xf1}, true); err != nil {
		t.Fatalf("Failed to call ledon: %v", err)
	}
	if err := readandCompare(conn, []byte{PktTypeRelayBoardData, 0x7, 0xc1, 0xd1, 0xe1, 0xf1, PktTypeActionReq, 0xd1, 0xe1, 0xf1, ActionLED, 1}); err != nil {
		t.Error(err)
	}

	// Buzzer.
	if err := core.BuzzerOn([]byte{0xd1, 0xe1, 0xf1}, true, 500); err != nil {
		t.Fatalf("Failed to call ledon: %v", err)
	}
	if err := readandCompare(conn, []byte{PktTypeRelayBoardData, 0x7, 0xc1, 0xd1, 0xe1, 0xf1, PktTypeActionReq, 0xd1, 0xe1, 0xf1, ActionBuzzer, 1, 0x01, 0xf4}); err != nil {
		t.Error(err)
	}

}

// readandCompare reads from conn and returns nil if want matches it.
func readandCompare(conn net.Conn, want []byte) error {
	buf := make([]byte, 255)
	n, err := conn.Read(buf)
	buf = buf[:n]
	if err != nil {
		return fmt.Errorf("Write to server failed:%v", err.Error())
	}
	if bytes.Compare(want, buf) != 0 {
		return fmt.Errorf("Failed: got %v want %v", buf, want)
	}
	return nil
}

// TestDataPacket tests if the packet send by relay is translated properly into an event.
// (needs to be run after TestRelayConfig or IP will not be registered)
func TestDataPacket(t *testing.T) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", hostPort)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		t.Fatalf("Failed to dial server:%v", err)
	}

	// Test Data.
	data := []struct {
		req     []byte
		pktType string
		addr    []byte
		hwaddr  []byte
		action  byte
		data    []byte
	}{
		{[]byte{0xA, 0xAD, 2, 2, 2, 2, 2, PktTypePing, 0xA, 0xA, 0xA},
			"Ping", []byte{0xA, 0xA, 0xA}, []byte{2, 2, 2, 2, 2}, 0, []byte{}},
		{[]byte{0x14, 0xAD, 2, 2, 2, 2, 2, PktTypeData, 0xA, 0xA, 0xA, ActionTemp, 0, 0xA0, 0x71, 0xD9, 0x41, 0x6B, 0xDE, 0x30, 0x42},
			"Temp", []byte{0xA, 0xA, 0xA}, []byte{2, 2, 2, 2, 2}, ActionTemp, []byte{80, 44}},
		{[]byte{0xD, 0xAD, 2, 2, 2, 2, 2, PktTypeData, 0xA, 0xA, 0xA, ActionMotion, 0, 1},
			"Motion", []byte{0xA, 0xA, 0xA}, []byte{2, 2, 2, 2, 2}, ActionMotion, []byte{}},
	}

	for _, v := range data {
		_, err = conn.Write(v.req)
		if err != nil {
			t.Fatalf("Failed to write %v", err)
		}

		// Set a timeout so we are not blocked on channel.
		var event Event
		func() {
			ticker := time.NewTicker(1 * time.Second)
			for {
				select {
				case event = <-core.Event():
					return
				case <-ticker.C:
					t.Fatalf("Timeout on channel")
				}
			}
		}()

		// Check Packet type.
		x := reflect.TypeOf(event)
		if x.Name() != v.pktType {
			t.Errorf("got %v, want %v", x.Name(), v.pktType)
		}

		// Check Addr.
		if bytes.Compare(event.Addr(), v.addr) != 0 {
			t.Errorf("got %v want %v", event.Addr(), v.addr)
		}

		// Check HWAddr.
		if bytes.Compare(event.PAddr(), v.hwaddr) != 0 {
			t.Errorf("got %v want %v", event.PAddr(), v.hwaddr)
		}

		// Check Packet specific stuff.
		switch f := event.(type) {
		case Ping:

		case Temp:
			if byte(f.Temp) != v.data[0] {
				t.Errorf("got %v, want %v", f.Temp, v.data[0])
			}

			if byte(f.Humidity) != v.data[1] {
				t.Errorf("got %v, want %v", f.Humidity, v.data[1])
			}

		case Motion:
			if !f.Motion {
				t.Errorf("need motion")
			}

		}
	}
}

func TestConfigMode(t *testing.T) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", hostPort)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		t.Fatalf("Failed to dial server:%v", err)
	}
	// Delay needed to stabilize listen on server.
	time.Sleep(100 * time.Millisecond)

	hwaddr := []byte{0xa1, 0xb1, 0xc1, 0xd1, 0xe1, 0xf1}
	defBrdAddr := []byte{0xff, 0xff, 0xff}

	// Validate relay config mode.
	if err := core.RelayConfigMode(hwaddr, true); err != nil {
		t.Error(err)
	}
	if err := readandCompare(conn,
		[]byte{PktTypeRelayCfgResp, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x2, 0xc1, 0xd1, 0xe1, 0xf1, 0x3, 0x4, 0x5, 0x6, 0x7, 0x73, 0xd1, 0xe1, 0xf1}); err != nil {
		t.Errorf("Failed: %v", err)
	}

	// Validate board config.
	if err := core.BoardConfig(defBrdAddr, defPipeAdress, hwaddr, []byte{0x1, 0x1, 0x1}); err != nil {
		t.Error(err)
	}
	if err := readandCompare(conn,
		[]byte{PktTypeRelayBoardData, 0x68, 0x65, 0x6C, 0x6C, 0x6F, PktTypeConfig, 0xff, 0xff, 0xff, 0xB, 0x1, 0x1, 0xc1, 0xd1, 0xe1, 0xf1, 0x1, 0x1, 0x1}); err != nil {
		t.Errorf("Failed: %v", err)
	}

	// Validate resetting relay config.
	if err := core.RelayConfigMode(hwaddr, false); err != nil {
		t.Error(err)
	}
	if err := readandCompare(conn,
		[]byte{PktTypeRelayCfgResp, 0x1, 0xc1, 0xd1, 0xe1, 0xf1, 0x2, 0xc1, 0xd1, 0xe1, 0xf1, 0x3, 0x4, 0x5, 0x6, 0x7, 0x73, 0xd1, 0xe1, 0xf1}); err != nil {
		t.Errorf("Failed: %v", err)
	}

}
