package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func handleUDPConnection(conn *net.UDPConn) {

	// here is where you want to do stuff like read or write to client

	buffer := make([]byte, 1024)

	n, addr, err := conn.ReadFromUDP(buffer)
	_ = addr
	fmt.Printf("%v Rcvd %db %v -- ", time.Now().Unix(), n, addr.IP)
	for i := 0; i < n; i++ {
		fmt.Printf("%X,", buffer[i])
	}
	fmt.Println()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	hostName := ""
	portNum := "6000"
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on port 6000")

	defer ln.Close()

	for {
		// wait for UDP client to connect
		handleUDPConnection(ln)
	}

}
