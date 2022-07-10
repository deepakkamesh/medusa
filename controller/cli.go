// Sender.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	hostPort := os.Args[2]
	RemoteAddr, err := net.ResolveUDPAddr("udp", hostPort)

	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	// note : you can use net.ResolveUDPAddr for LocalAddr as well
	//        for this tutorial simplicity sake, we will just use nil

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Established connection to %s \n", hostPort)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	defer conn.Close()

	pkt := strings.Split(os.Args[3], ",")
	fmt.Println(pkt)
	message := []byte{}
	for i := 0; i < len(pkt); i++ {
		v, _ := strconv.ParseUint(pkt[i], 16, 8)
		message = append(message, byte(v))
	}
	fmt.Println(message)
	// write a message to server
	_, err = conn.Write(message)

	if err != nil {
		log.Println(err)
	}

}
