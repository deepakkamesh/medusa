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
		{[]byte{0xAA, 67, 68, 69, 70, 55, 21}, []byte{0xAB, 10, 10, 2, 3, 4, 60, 7, 8, 9, 10, 61, 62, 63, 64, 115, 67, 1, 1}},
		{[]byte{0xAA, 31, 32, 33, 34, 35, 36}, []byte{0xAB, 16, 7, 8, 9, 10, 71, 7, 8, 9, 10, 72, 73, 74, 75, 115, 77, 1, 1}},
	}

	// Send config request to server.
	for _, v := range data {
		_, err = conn.Write(v.req)

		if err != nil {
			t.Error(err)
		}
		t.Logf("Config request:%v", v.req)
		//  Receive config response.
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("%v", err)
		}
		buffer = buffer[:n]
		t.Logf("Config response:%v", buffer)
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
	if err := core.Temp([]byte{4, 5, 3}); err != nil {
		t.Fatalf("Failed to call temp: %v", err)
	}
	if err := readandCompare(conn, []byte{0xAD, 60, 7, 8, 9, 10, PktTypeActionReq, 4, 5, 3, ActionTemp}); err != nil {
		t.Error(err)
	}
	// LedOn.
	if err := core.LEDOn([]byte{4, 5, 3}, true); err != nil {
		t.Fatalf("Failed to call ledon: %v", err)
	}
	if err := readandCompare(conn, []byte{0xAD, 60, 7, 8, 9, 10, PktTypeActionReq, 4, 5, 3, ActionLED, 1}); err != nil {
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
		{[]byte{0xAD, 2, 2, 2, 2, 2, PktTypePing, 0xA, 0xA, 0xA},
			"Ping", []byte{0xA, 0xA, 0xA}, []byte{2, 2, 2, 2, 2}, 0, []byte{}},
		{[]byte{0xAD, 2, 2, 2, 2, 2, PktTypeData, 0xA, 0xA, 0xA, ActionTemp, 0, 35, 40},
			"Temp", []byte{0xA, 0xA, 0xA}, []byte{2, 2, 2, 2, 2}, ActionTemp, []byte{35, 40}},
		{[]byte{0xAD, 2, 2, 2, 2, 2, PktTypeData, 0xA, 0xA, 0xA, ActionMotion, 0},
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
				case event = <-core.Event:
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
			_ = f

		case Temp:
			if f.Temp != v.data[0] {
				t.Errorf("got %v, want %v", f.Temp, v.data[0])
			}

			if f.Humidity != v.data[1] {
				t.Errorf("got %v, want %v", f.Humidity, v.data[1])
			}

		case Motion:
			_ = f

		default:
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

	hwaddr := []byte{67, 68, 69, 70, 55, 21}
	defBrdAddr := []byte{0xff, 0xff, 0xff}

	// Validate relay config mode.
	if err := core.SetRelayConfigMode(hwaddr, true); err != nil {
		t.Error(err)
	}
	if err := readandCompare(conn, []byte{PktTypeRelayCfgResp, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 60, 7, 8, 9, 10, 61, 62, 63, 64, 115, 67, 1, 1}); err != nil {
		t.Errorf("Failed: %v", err)
	}

	// Validate board config.
	if err := core.SetBoardConfig(defBrdAddr, defPipeAdress, []byte{4, 5, 3}, hwaddr); err != nil {
		t.Error(err)
	}
	if err := readandCompare(conn, []byte{PktTypeRelayBoardData, 0x68, 0x65, 0x6C, 0x6C, 0x6F, PktTypeConfig, 0xff, 0xff, 0xff, 15, 1, 60, 7, 8, 9, 10, 4, 5, 3}); err != nil {
		t.Errorf("Failed: %v", err)
	}

	// Validate resetting relay config.
	if err := core.SetRelayConfigMode(hwaddr, false); err != nil {
		t.Error(err)
	}
	if err := readandCompare(conn, []byte{PktTypeRelayCfgResp, 10, 10, 2, 3, 4, 60, 7, 8, 9, 10, 61, 62, 63, 64, 115, 67, 1, 1}); err != nil {
		t.Errorf("Failed: %v", err)
	}

}
