package core

import (
	"bytes"
	"testing"
)

func TestConfigFileLoad(t *testing.T) {

	conf, err := NewConfig("core.cfg.test.json")
	if err != nil {
		t.Errorf("Error loading config %v", err)
	}
	data := []struct {
		got  []byte
		want []byte
	}{
		{got: conf.Boards[0].Addr, want: []byte{0x1, 0x1, 0x1}},
		{got: conf.Boards[0].PAddr, want: []byte{0x1, 0xc1, 0xd1, 0xe1, 0xf1}},
		{got: conf.Boards[0].Actions, want: []byte{1, 2, 3, 5, 0x14}},
		{got: conf.Relays[0].HWAddr, want: []byte{0xa1, 0xb1, 0xc1, 0xd1, 0xe1, 0xf1}},
		{got: conf.Relays[0].PAddr0, want: []byte{0x1, 0xc1, 0xd1, 0xe1, 0xf1}},
	}

	for _, v := range data {
		if bytes.Compare(v.got, v.want) != 0 {
			t.Errorf("Failed got:%v want:%v", v.got, v.want)
		}
	}

	if conf.Boards[0].Room != "living" {
		t.Errorf("Failed got:%v want:%v", conf.Boards[0].Room, "living")
	}

}
