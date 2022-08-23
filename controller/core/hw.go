package core

import "net"

type Relay struct {
	conn    net.Conn
	connUDP net.PacketConn
	addr    net.Addr
}
