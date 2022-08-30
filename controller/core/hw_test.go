package core

import (
	"bytes"
	"testing"
)

func TestConfigFileLoad(t *testing.T) {

	conf, err := newConfig("core.cfg.json")
	if err != nil {
		t.Errorf("Error loading config %v", err)
	}
	data := []struct {
		got  []byte
		want []byte
	}{
		{got: conf.Boards[0].Addr, want: []byte{4, 5, 3, 5}},
		{got: conf.Boards[0].PAddr, want: []byte{6, 7, 8, 9, 10}},
		{got: conf.Relays[0].HWAddr, want: []byte{1, 2, 3, 4, 5}},
	}

	for _, v := range data {
		if bytes.Compare(v.got, v.want) != 0 {
			t.Errorf("Failed got:%v want:%v", v.got, v.want)
		}
	}

	if conf.Boards[0].Room != "Family" {
		t.Errorf("Failed got:%v want:%v", conf.Boards[0].Room, "Family")
	}
}
