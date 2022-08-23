package core

import (
	"fmt"
	"net"

	"github.com/golang/glog"
)

// StartPacketHandlers starts up the packet handlers.
func (c *Core) StartPacketHandlers() {
	go c.boardPktHandler()
	go c.relayConfigHandler()
}

// boardPktHandler handles the TCP connection for board data exchange.
func (c *Core) boardPktHandler() {
	l, err := net.Listen("tcp", c.hostPort)
	if err != nil {
		glog.Fatalf("Error listening %v", err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			glog.Fatalf("Error accepting tcp connection: %v", err)
		}
		go c.handleRequest(conn)
	}
}

// Handles incoming requests.
func (c *Core) handleRequest(conn net.Conn) {
	// Save TCP connection details.
	c.relays["0"].conn = conn
	defer conn.Close()

	for {
		buf := make([]byte, 255)
		sz, err := conn.Read(buf)
		if err != nil {
			glog.Errorf("Error reading connection: %v", err)
			return
		}
		glog.Infof(PrintPkt("tcp pkt:", buf, sz))
	}
}

// Start UDP connection for config requests.
func (c *Core) relayConfigHandler() {

	lu, err := net.ListenPacket("udp", c.hostPort)
	if err != nil {
		glog.Fatalf("Error listening UDP: %v", err)
	}
	defer lu.Close()

	for {
		c.sendRelayConfig(lu)
	}
}

// sendRelayConfig handles new UDP packets.
func (c *Core) sendRelayConfig(conn net.PacketConn) {

	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFrom(buffer)
	if err != nil {
		glog.Fatalf("Failed reading from UDP: %v", err)
	}

	// Save connection details.
	// TODO: Add after checking the MAC address.
	c.relays["0"].connUDP = conn
	c.relays["0"].addr = addr

	preamble := fmt.Sprintf("Got UDP %db %v -- ", n, addr.String())
	glog.Infof(PrintPkt(preamble, buffer, n))
}

// PrintPkt returns a string of packet formatted properly.
func PrintPkt(preamble string, b []byte, s int) string {
	logMsg := preamble
	for i := 0; i < s; i++ {
		logMsg = logMsg + fmt.Sprintf("%X ", b[i])
	}
	return logMsg
}
