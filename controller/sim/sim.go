package main

import (
	"log"
	"net"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
)

var hostPort string = "127.0.0.1:3334"
var hwaddr []byte = []byte{0xa1, 0xb1, 0xc1, 0xd1, 0xe1, 0xf1}

func main() {
	SendRelayConfigReq()
	time.Sleep(100 * time.Millisecond)
	TimedEvents()

}

func TimedEvents() {
	t2 := time.NewTicker(80 * time.Second)

	tcpAddr, _ := net.ResolveTCPAddr("tcp", hostPort)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Failed to dial server:%v", err)
	}

	go func() {
		for {
			pkt := make([]byte, 255)
			n, err := conn.Read(pkt)
			if err != nil {
				log.Printf("Read to server failed:%v", err.Error())
				return
			}
			pkt = pkt[:n]
			log.Printf(core.PP(pkt, "Pkt:"))
			if pkt[0] != core.PktTypeRelayBoardData {
				continue
			}
			switch pkt[6] {
			case core.PktTypeActionReq:
				ProcessAction(pkt, conn)
			}
		}
	}()

	for {
		select {
		case <-t2.C:
			pkt := []byte{0xAD, 2, 2, 2, 2, 2, core.PktTypePing, 0xA, 0xA, 0xA}
			if _, err := conn.Write(pkt); err != nil {
				log.Printf("Write to server failed:%v", err.Error())
			}
		}
	}
}

func ProcessAction(pkt []byte, conn net.Conn) {
	switch pkt[10] {
	case core.ActionTemp:
		pkt := []byte{0xAD, 2, 2, 2, 2, 2, core.PktTypeData, 0xA, 0xA, 0xA, core.ActionTemp, core.ErrNA, 30, 45}
		if _, err := conn.Write(pkt); err != nil {
			log.Printf("Write to server failed:%v", err.Error())
		}

	}

}

func SendRelayConfigReq() {
	addr, _ := net.ResolveUDPAddr("udp", hostPort)
	conn, err := net.DialUDP("udp", nil, addr)

	if err != nil {
		log.Fatalf("Could not dial server:%v", err)
	}
	pkt := []byte{core.PktTypeRelayCfgReq}
	pkt = append(pkt, hwaddr...)

	_, err = conn.Write(pkt)

	if err != nil {
		log.Print(err)
	}
	log.Printf(core.PP(pkt, "Config request:"))

	//  Receive config response.
	buffer := make([]byte, 255)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("%v", err)
	}
	buffer = buffer[:n]
	log.Printf("Config response:%v", buffer)
	defer conn.Close()
}
