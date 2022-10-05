package controller_test

import (
	"testing"
	"time"

	"github.com/deepakkamesh/medusa/controller"
	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/deepakkamesh/medusa/controller/mocks"
	"github.com/golang/mock/gomock"
)

func TestMotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockMedusaCore(ctrl)

	c, e := controller.NewController(m, ":3344")
	if e != nil {
		t.Errorf("Failed init Controller %v", e)
	}

	// Startup Core.
	m.EXPECT().StartCore()
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
	//gomock.InOrder(
	m.EXPECT().GetBoardByAddr([]byte{1, 1, 1}).AnyTimes().Return(&core.Board{Room: "living"})
	m.EXPECT().GetBoardByAddr([]byte{2, 2, 2}).AnyTimes().Return(&core.Board{Room: "hallway"})
	//)
	c.Run()

}
