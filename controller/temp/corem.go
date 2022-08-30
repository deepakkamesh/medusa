package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
)

func main() {

	core, err := core.NewCore("127.0.0.1:3344", "../core/core.cfg.test.json")

	if err != nil {
		fmt.Println(err)
	}
	core.StartCore()

	for {

		RemoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3344")

		conn, err := net.DialUDP("udp", nil, RemoteAddr)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
		log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

		defer conn.Close()

		// write a message to server

		message := []byte{0xAA, 1, 2, 3, 4, 5, 6}
		_, err = conn.Write(message)

		if err != nil {
			log.Println(err)
		}

		// receive message from server
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)

		for i := 0; i < n; i++ {
			fmt.Printf("%0X ", buffer[i])
		}
		fmt.Println("UDP Server : ", addr)
		//fmt.Println("Received from UDP server : ", string(buffer[:n]))
		time.Sleep(1 * time.Second)
	}
}
