package core

import (
	"fmt"
	"net"

	"github.com/golang/glog"
)

func (c *Core) StartPacketHandler() {
	l, err := net.Listen("tcp", c.hostPort)
	if err != nil {
		glog.Fatalf("Error listening %v", err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			glog.Fatalf("Error accepting:%v", err)
		}
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	for {
		select {
		case v := <-c.tempChan:
			_ = v

		default:
			buf := make([]byte, 255)
			sz, err := conn.Read(buf)
			if err != nil {
				glog.Warningf("Error reading connection: %v", err)
				conn.Close()
				return
			}
			fmt.Printf("got  %v %v \n", sz, buf)
			// Send a response back to person contacting us.
			conn.Write([]byte("Ack"))
			conn.Write([]byte("hell"))
			// Close the connection when you're done with it.
			//	conn.Close()
		}
	}
}
