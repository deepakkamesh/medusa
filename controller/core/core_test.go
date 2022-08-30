package core

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

var core *Core

func init() {
	var err error
	core, err = NewCore("127.0.0.1:3344", "core.cfg.test.json")
	if err != nil {
		fmt.Printf("Error initializing core:%v", err)
		os.Exit(1)
	}
	core.StartCore()
	// Add some delay for the service to listen & get ready.
	time.Sleep(100 * time.Millisecond)
}

func TestRelayConfig(t *testing.T) {
	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:3344")
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		t.Errorf("Could not dial server:%v", err)
	}

	defer conn.Close()

	data := []struct {
		req  []byte
		resp []byte
	}{
		{[]byte{0xAA, 1, 2, 3, 4, 5, 6}, []byte{0xAB, 0xA, 1, 2, 3, 4, 6, 7, 8, 9, 10, 6, 6, 6, 6, 115}},
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
