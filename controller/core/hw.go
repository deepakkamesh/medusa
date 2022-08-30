package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
)

type Relay struct {
	conn        net.Conn
	IP          net.IP
	HWAddr      []byte
	PAddr0      []byte
	PAddr1      []byte
	PAddr2      []byte
	PAddr3      []byte
	PAddr4      []byte
	PAddr5      []byte
	Addr        []byte // Address of virt board on relay.
	Channel     byte
	Description string
	Room        string
	Name        string
}

type Board struct {
	Addr        []byte
	PAddr       []byte
	ARD         byte
	PingInt     byte
	Description string
	Room        string
	Name        string
}

type Config struct {
	Relays []*Relay
	Boards []*Board
}

func newConfig(filepath string) (*Config, error) {
	cfg, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = json.Unmarshal([]byte(cfg), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (f *Config) getRelaybyIP(ip net.IP) *Relay {
	for _, v := range f.Relays {
		if ip.Equal(v.IP) {
			return v
		}
	}
	return nil
}

func (f *Config) getRelayByHWAddr(hwaddr []byte) *Relay {
	for _, v := range f.Relays {
		if bytes.Compare(v.HWAddr, hwaddr) == 0 {
			return v
		}
	}
	return nil
}

// SendPKt sends the packet
func (f *Config) SendPkt(addr []byte, pkt []byte) error {
	return nil
}
