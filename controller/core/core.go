package core

import (
	"fmt"
	"net"

	"github.com/golang/glog"
)

// Core is the main struct for the Medusa Core handling.
type Core struct {
	hostPort string // IP Port for TCP & UDP bindings.
	config   *Config
	Pkt      chan Packet // Channel to receive the packet.
}

// NewCore returns an initialized Core.
func NewCore(hostPort string, cfgFname string) (*Core, error) {
	config, err := newConfig(cfgFname)
	if err != nil {
		return nil, err
	}

	return &Core{
		hostPort: hostPort,
		config:   config,
		Pkt:      make(chan Packet),
	}, nil
}

// StartCore starts up the packet handlers.
func (c *Core) StartCore() {
	go c.boardPktHandler()
	go c.relayConfigHandler()
}

// RequestAction sends an action request to the board addr.
func (c *Core) RequestAction(addr []byte, actionID byte, data []byte) {

}

// SendRawPacket sends the raw packet to the board.
func (c *Core) SendRawPacket(pkt []byte) error {
	//	_, err := c.relays["0"].conn.Write(pkt)
	return nil
}

// SendManualRelayCfg sends the manual config response for relay.
func (c *Core) SendManualRelayCfg(pkt []byte) error {
	//_, err := c.relays["0"].connUDP.WriteTo(pkt, c.relays["0"].IP)
	return nil
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

		// Save TCP connection details.
		ip := conn.RemoteAddr().(*net.TCPAddr).IP
		relay := c.config.getRelaybyIP(ip)
		if relay == nil {
			glog.Errorf("Relay with IP not registered:%v", ip)
			continue
		}
		relay.conn = conn

		// Now handle the request.
		go c.handleRequest(conn)
	}
}

// Handles incoming requests.
func (c *Core) handleRequest(conn net.Conn) {
	// Save TCP connection details.
	//	c.relays["0"].conn = conn
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
	preamble := fmt.Sprintf("%v - PktTypeRelayCfgReq:", addr.String())
	glog.Infof(PrintPkt(preamble, buffer, n))

	// Validate the config request packet.
	if !okPktTypeRelayCfgReq(buffer) {
		glog.Error("Unknown packet type in UDP connection or bad pkt len")
		return
	}

	// Get relay associated with the mac.
	mac := getHWAddr(buffer)
	relay := c.config.getRelayByHWAddr(mac)
	if relay == nil {
		glog.Errorf("HWAddr not found:%v", mac)
		return
	}
	relay.IP = addr.(*net.UDPAddr).IP

	// Send Relay config.
	relayCfg := makePktTypeRelayCfgResp(relay)
	preamble = fmt.Sprintf("%v - PktTypeRelayCfgResp:", addr.String())
	glog.Infof(PrintPkt(preamble, relayCfg, len(relayCfg)))
	if _, err := conn.WriteTo(relayCfg, addr); err != nil {
		glog.Errorf("Failed to send relay config:%v", err)
	}
}
