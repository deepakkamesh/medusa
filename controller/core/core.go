package core

import (
	"fmt"
	"net"

	"github.com/golang/glog"
)

// Core is the main struct for the Medusa Core handling.
type Core struct {
	hostPort string // IP Port for TCP & UDP bindings.
	conf     *Config
	Event    chan Event // Channel to receive events.
}

// NewCore returns an initialized Core.
func NewCore(hostPort string, cfgFname string) (*Core, error) {
	config, err := newConfig(cfgFname)
	if err != nil {
		return nil, err
	}

	return &Core{
		hostPort: hostPort,
		conf:     config,
		Event:    make(chan Event),
	}, nil
}

// RequestAction sends an action request to the board addr.
func (c *Core) requestAction(addr []byte, actionID byte, data []byte) error {
	brd := c.conf.getBoardByAddr(addr)
	if brd == nil {
		return fmt.Errorf("address not found %v", addr)
	}
	relay := c.conf.getRelayByPAddr(brd.PAddr)
	if relay == nil {
		return fmt.Errorf("relay not found for pipe address %v", brd.PAddr)
	}
	pkt := makePktTypeActionReq(actionID, brd.Addr, brd.PAddr, data)

	if relay.conn == nil {
		return fmt.Errorf("relay not registered. hwaddr:%v", relay.HWAddr)
	}
	_, err := relay.conn.Write(pkt)
	if err != nil {
		return err
	}
	return nil
}

func (c *Core) Light(addr []byte) error {
	return c.requestAction(addr, ActionLight, []byte{})
}

func (c *Core) Temp(addr []byte) error {
	return c.requestAction(addr, ActionTemp, []byte{})
}

func (c *Core) LEDOn(addr []byte, on bool) error {
	var data byte = 0
	if on {
		data = 1
	}
	return c.requestAction(addr, ActionLED, []byte{data})
}

// BoardConfig sends the board configuration associated with naddr in
// config file to board address default addr and paddr.
func (c *Core) SetBoardConfig(addr []byte, paddr []byte, naddr []byte, hwaddr []byte) error {

	relay := c.conf.getRelayByHWAddr(hwaddr)
	if relay == nil {
		return fmt.Errorf("relay not found for hwaddr %v", hwaddr)
	}
	if relay.conn == nil {
		return fmt.Errorf("relay not registered. hwaddr:%v", relay.HWAddr)
	}

	brd := c.conf.getBoardByAddr(naddr)
	if brd == nil {
		return fmt.Errorf("address not found %v", addr)
	}

	pkt := makePktTypeConfig(addr, paddr, brd)

	_, err := relay.conn.Write(pkt)
	if err != nil {
		return err
	}
	return nil
}

// SetRelayConfigMode sets relay with hwaddr in config mode.
// if ok is false, unsets config mode.
func (c *Core) SetRelayConfigMode(hwaddr []byte, yes bool) error {
	relay := c.conf.getRelayByHWAddr(hwaddr)
	if relay == nil {
		return fmt.Errorf("relay not found for hwaddr %v", hwaddr)
	}
	if relay.conn == nil {
		return fmt.Errorf("relay not registered. hwaddr:%v", relay.HWAddr)
	}
	// Send Relay config.
	relayCfg := makePktTypeRelayCfgResp(relay, yes)
	_, err := relay.conn.Write(relayCfg)
	return err
}

// StartCore starts up the packet handlers.
func (c *Core) StartCore() {
	go c.boardPacketHandler()
	go c.relayConfigHandler()
}

// boardPacketHandler handles the TCP connection for board data exchange.
func (c *Core) boardPacketHandler() {
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

		// Save TCP conn based on IP info during relay registration.
		ip := conn.RemoteAddr().(*net.TCPAddr).IP
		relay := c.conf.getRelaybyIP(ip)
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
	defer conn.Close()

	for {
		buf := make([]byte, 255)
		n, err := conn.Read(buf)
		if err != nil {
			glog.Errorf("Error reading connection: %v", err)
			return
		}
		buf = buf[:n]
		preamble := fmt.Sprintf("%v - Pkt:", conn.RemoteAddr())
		glog.Infof(PrintPkt(preamble, buf, n))

		ip := conn.RemoteAddr().(*net.TCPAddr).IP

		relay := c.conf.getRelaybyIP(ip)

		// Translate packet to event and send to channel.
		event, err := translatePacket(buf, relay.HWAddr)
		if err != nil {
			glog.Errorf("Unable to translate packet:%v", err)
			continue
		}
		c.Event <- event
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

	buf := make([]byte, 255)
	n, addr, err := conn.ReadFrom(buf)
	if err != nil {
		glog.Fatalf("Failed reading from UDP: %v", err)
	}
	buf = buf[:n]
	preamble := fmt.Sprintf("%v - PktTypeRelayCfgReq:", addr.String())
	glog.Infof(PrintPkt(preamble, buf, n))

	// Validate the config request packet.
	if !okPktTypeRelayCfgReq(buf) {
		glog.Errorf("%v - Unknown packet type in UDP connection or bad pkt len", addr.String())
		return
	}

	// Get relay associated with the mac.
	mac := getHWAddr(buf)
	relay := c.conf.getRelayByHWAddr(mac)
	if relay == nil {
		glog.Errorf("%v - HWAddr not found:%v", addr.String(), mac)
		return
	}
	relay.IP = addr.(*net.UDPAddr).IP // Register IP from relay.

	// Send Relay config.
	relayCfg := makePktTypeRelayCfgResp(relay, false)
	preamble = fmt.Sprintf("%v - PktTypeRelayCfgResp:", addr.String())
	glog.Infof(PrintPkt(preamble, relayCfg, len(relayCfg)))
	if _, err := conn.WriteTo(relayCfg, addr); err != nil {
		glog.Errorf("Failed to send relay config:%v", err)
	}
}
