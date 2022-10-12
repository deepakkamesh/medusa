package controller_test

import (
	"testing"
	"time"

	"github.com/deepakkamesh/medusa/controller"
	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/deepakkamesh/medusa/controller/mocks"
	"github.com/golang/mock/gomock"
)

/*
func TestConfigGen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockMedusaCore(ctrl)

	c, e := controller.NewController(m, ":3344" )
	if e != nil {
		t.Fatalf("Failed init Controller %v", e)

	}

	cr, e := core.NewConfig("core/core.cfg.json")
	if e != nil {
		t.Errorf("failed to open %v %v", "core.dev", e)
	}

	// Startup Core.
	m.EXPECT().StartCore()
	m.EXPECT().CoreConfig().Return(cr)
	if err := c.Startup(); err != nil {
		t.Fatalf("Failed to startup Core: %v", err)
	}

	// Send some events.
	eventChan := make(chan core.Event)

	p1 := core.Motion{
		PktInfo: core.PktInfo{
			BoardAddr:    []byte{1, 1, 1},
			PipeAddr:     []byte{},
			HardwareAddr: []byte{},
		},
		Motion: true,
	}
	p2 := core.Motion{
		PktInfo: core.PktInfo{
			BoardAddr:    []byte{1, 1, 1},
			PipeAddr:     []byte{},
			HardwareAddr: []byte{},
		},
		Motion: false,
	}

	pkts := []core.Event{p1, p2}

	go func() {
		i := 0
		for i = 0; i < len(pkts); i++ {

			time.Sleep(500 * time.Millisecond)
			eventChan <- pkts[i]
		}
		// To kill the forever loop.
		time.Sleep(1 * time.Second)
		close(eventChan)
	}()

	m.EXPECT().Event().AnyTimes().Return(eventChan)
	m.EXPECT().GetBoardByAddr([]byte{1, 1, 1}).AnyTimes().Return(&core.Board{Room: "dining", Name: "b1"})
	//m.EXPECT().GetBoardByAddr([]byte{2, 2, 2}).AnyTimes().Return(&core.Board{Room: "hallway-down"})
	c.Run()

}*/

func TestMotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockMedusaCore(ctrl)
	h := mocks.NewMockHA(ctrl)

	c, e := controller.NewController(m, h, ":3344")
	if e != nil {
		t.Errorf("Failed init Controller %v", e)
	}

	// Startup Core.
	m.EXPECT().StartCore()
	h.EXPECT().Connect()
	if err := c.Startup(); err != nil {
		t.Errorf("Failed to startup Core: %v", err)
	}

	// Send some events.
	eventChan := make(chan core.Event)

	p1 := core.Motion{
		PktInfo: core.PktInfo{
			BoardAddr:    []byte{1, 1, 1},
			PipeAddr:     []byte{},
			HardwareAddr: []byte{},
		},
		Motion: true,
	}
	p2 := core.Motion{
		PktInfo: core.PktInfo{
			BoardAddr:    []byte{2, 2, 2},
			PipeAddr:     []byte{},
			HardwareAddr: []byte{},
		},
		Motion: true,
	}

	pkts := []core.Event{p1, p2}

	go func() {
		i := 0
		for i = 0; i < len(pkts); i++ {

			time.Sleep(100 * time.Millisecond)
			eventChan <- pkts[i]
		}
		// To kill the forever loop.
		time.Sleep(1 * time.Second)
		close(eventChan)
	}()

	m.EXPECT().Event().AnyTimes().Return(eventChan)
	h.EXPECT().SendSensorData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	m.EXPECT().GetBoardByAddr([]byte{1, 1, 1}).AnyTimes().Return(&core.Board{Room: "living", Name: "b1"})
	m.EXPECT().GetBoardByAddr([]byte{2, 2, 2}).AnyTimes().Return(&core.Board{Room: "hallway-down", Name: "b1"})
	c.Run()
}
