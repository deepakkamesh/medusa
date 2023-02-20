package core

import (
	"bytes"
	"io/ioutil"
	"net"

	"github.com/flynn/json5"
	"golang.org/x/exp/slices"
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
	PAddr6      []byte // Pipe Address for the virtual board.
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
	Actions     []byte // Supported actions.
}

// IsActionCapable returns true if the board is capable of the action.
func (b *Board) IsActionCapable(actionID byte) bool {
	return slices.Contains(b.Actions, actionID)
}

type Config struct {
	Relays []*Relay
	Boards []*Board
}

func NewConfig(filepath string) (*Config, error) {
	cfg, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = json5.Unmarshal([]byte(cfg), c)
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

func (f *Config) getRelayByPAddr(paddr []byte) *Relay {
	for _, v := range f.Relays {
		if bytes.Compare(v.PAddr0, paddr) == 0 {
			return v
		}
		if bytes.Compare(v.PAddr1, paddr) == 0 {
			return v
		}
		if bytes.Compare(v.PAddr2, paddr) == 0 {
			return v
		}
		if bytes.Compare(v.PAddr3, paddr) == 0 {
			return v
		}
		if bytes.Compare(v.PAddr4, paddr) == 0 {
			return v
		}
		if bytes.Compare(v.PAddr5, paddr) == 0 {
			return v
		}
		if bytes.Compare(v.PAddr6, paddr) == 0 {
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

func (f *Config) getBoardByAddr(addr []byte) *Board {
	for _, v := range f.Boards {
		if bytes.Compare(v.Addr, addr) == 0 {
			return v
		}
	}
	return nil
}

func (f *Config) getBoardByName(name string) *Board {
	for _, v := range f.Boards {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// getBoardByRoom returns all the boards. If room is "all"
// it returns all the boards.
func (f *Config) getBoardByRoom(room string) []Board {
	boards := []Board{}
	for _, v := range f.Boards {
		if v.Room == room || room == "all" {
			boards = append(boards, *v)
		}
	}
	return boards
}
