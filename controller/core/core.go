package core

import (
	"fmt"
	"net"
	"time"

	"github.com/golang/glog"
)

//go:generate mockgen -destination=../mocks/core_mock.go -package=mocks github.com/deepakkamesh/medusa/controller/core MedusaCore

// Interface definition.
type MedusaCore interface {
	Action(addr []byte, actionID byte, data []byte) error
	Light(addr []byte) error
	Temp(addr []byte) error
	Volt(addr []byte) error
	Reset(addr []byte) error
	LEDOn(addr []byte, on bool) error
	BuzzerOn(addr []byte, on bool, d int) error
	BoardConfig(addr []byte, paddr []byte, hwaddr []byte, naddr []byte) error
	RelayConfigMode(hwaddr []byte, yes bool) error
	StartCore()
	Event() <-chan Event
	GetBoardByAddr(b []byte) *Board
	GetBoardByName(name string) *Board
	GetBoardByRoom(room string) []Board
	GetRelaybyPAddr(paddr []byte) *Relay
	CoreConfig() *Config
}

// Core is the main struct for the Medusa Core handling.
type Core struct {
	hostPort string     // IP Port for TCP & UDP bindings.
	conf     *Config    // Config holds the hardware configuration.
	event    chan Event // Channel to send events.
}

// NewCore returns an initialized Core.
func NewCore(hostPort string, cfgFname string) (*Core, error) {
	config, err := NewConfig(cfgFname)
	if err != nil {
		return nil, err
	}

	return &Core{
		hostPort: hostPort,
		conf:     config,
		event:    make(chan Event, 50),
	}, nil
}

// CoreConfig returns the core config struct.
func (c *Core) CoreConfig() *Config {
	return c.conf
}

// GetBoardByName returns board info.
func (c *Core) GetBoardByName(name string) *Board {
	return c.conf.getBoardByName(name)
}

// GetBoardByAddr returns board info.
func (c *Core) GetBoardByAddr(addr []byte) *Board {
	return c.conf.getBoardByAddr(addr)
}

func (c *Core) GetBoardByRoom(room string) []Board {
	return c.conf.getBoardByRoom(room)
}

func (c *Core) GetRelaybyPAddr(paddr []byte) *Relay {
	return c.conf.getRelayByPAddr(paddr)
}

// Event returns the channel for events.
func (c *Core) Event() <-chan Event {
	return c.event
}

// Action sends an action request to the board addr.
func (c *Core) Action(addr []byte, actionID byte, data []byte) error {
	brd := c.conf.getBoardByAddr(addr)
	if brd == nil {
		return fmt.Errorf("address not found %v", addr)
	}
	relay := c.conf.getRelayByPAddr(brd.PAddr)
	if relay == nil {
		return fmt.Errorf("relay not found for pipe address %v", brd.PAddr)
	}
	if relay.conn == nil {
		return fmt.Errorf("relay not registered. hwaddr:%v", relay.HWAddr)
	}

	pkt := makePktTypeActionReq(actionID, brd.Addr, brd.PAddr, data)
	glog.Info(PP(pkt, "PktTypeActionReq:"))
	// TODO: Packets sent too fast to the same relay can corrupt. Need a delay or
	// some guarantee of staggered request.(quick fix is to sleep).
	_, err := relay.conn.Write(pkt)
	time.Sleep(300 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

// Light gets light level.
func (c *Core) Light(addr []byte) error {
	return c.Action(addr, ActionLight, []byte{})
}

// Temp - temp and humidity.
func (c *Core) Temp(addr []byte) error {
	return c.Action(addr, ActionTemp, []byte{})
}

// Volt gets the voltage of the board.
func (c *Core) Volt(addr []byte) error {
	return c.Action(addr, ActionVolt, []byte{})
}

// Reset resets the board.
func (c *Core) Reset(addr []byte) error {
	return c.Action(addr, ActionReset, []byte{})
}

// Sets buzzer on/off.
func (c *Core) BuzzerOn(addr []byte, on bool, d int) error {
	var data byte = 0
	if on {
		data = 1
	}
	hiD := byte(d >> 8)
	loD := byte(d & 0xFF)

	return c.Action(addr, ActionBuzzer, []byte{data, hiD, loD})
}

// LEDOn sets led on or off.
func (c *Core) LEDOn(addr []byte, on bool) error {
	var data byte = 0
	if on {
		data = 1
	}
	return c.Action(addr, ActionLED, []byte{data})
}

// BoardConfig sends the board configuration associated with naddr in
// config file to board address default addr and paddr.
func (c *Core) BoardConfig(addr []byte, paddr []byte, hwaddr []byte, naddr []byte) error {

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
	glog.Info(PP(pkt, "PktTypeConfig:"))

	_, err := relay.conn.Write(pkt)
	if err != nil {
		return err
	}
	return nil
}

// RelayConfigMode sets relay with hwaddr in config mode.
// if ok is false, unsets config mode.
func (c *Core) RelayConfigMode(hwaddr []byte, yes bool) error {
	relay := c.conf.getRelayByHWAddr(hwaddr)
	if relay == nil {
		return fmt.Errorf("relay not found for hwaddr %v", hwaddr)
	}
	if relay.conn == nil {
		return fmt.Errorf("relay not registered. hwaddr:%v", relay.HWAddr)
	}
	// Send Relay config.
	pkt := makePktTypeRelayCfgResp(relay, yes)
	glog.Info(PP(pkt, "PktTypeRelayCfgResp:"))

	_, err := relay.conn.Write(pkt)
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
		go c.handleRequest(conn, relay.HWAddr)
	}
}

// Handles incoming requests.
func (c *Core) handleRequest(conn net.Conn, hwaddr []byte) {
	defer conn.Close()

	for {
		buf := make([]byte, 255)
		n, err := conn.Read(buf)
		if err != nil {
			glog.Errorf("Error reading connection: %v", err)
			return
		}
		buf = buf[:n]
		glog.V(1).Info(PP(buf, "%v - Pkt(%v):", conn.RemoteAddr(), n))

		// Break buffer into packets.
		pkts, e := splitPackets(buf)
		if e != nil {
			glog.Errorf("Failed to split packets:%v", e)
			continue
		}
		for _, pkt := range pkts {
			// Translate packet to event and send to channel.
			event, err := translatePacket(pkt, hwaddr)
			if err != nil {
				glog.Errorf("Unable to translate packet:%v", err)
				continue
			}
			c.event <- event
		}
	}
}

func splitPackets(a []byte) (b [][]byte, err error) {
	st := 0
	l := 0
	for {
		if st >= len(a) {
			return
		}
		l = int(a[st])
		if len(a[st+1:]) < l {
			err = fmt.Errorf("bad len")
			return
		}
		b = append(b, a[st+1:st+l+1])
		st = st + l + 1
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
	glog.Infof(PP(buf, "%v - PktTypeRelayCfgReq:", addr.String()))

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
	if _, err := conn.WriteTo(relayCfg, addr); err != nil {
		glog.Errorf("Failed to send relay config:%v", err)
	}
	glog.Info(PP(relayCfg, "%v - PktTypeRelayCfgResp:", addr.String()))
}
