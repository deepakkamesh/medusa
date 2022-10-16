package controller_test

import (
	"testing"
	"time"

	"github.com/deepakkamesh/medusa/controller"
	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/deepakkamesh/medusa/controller/mocks"
	"github.com/golang/mock/gomock"
)

func TestHAMsgHandler(t *testing.T) {
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

	// Channel to send events from HA.
	msgChan := make(chan controller.HAMsg)

	// HA messages to send.
	msgs := []controller.HAMsg{
		{
			Room:   "living",
			Action: "buzzer",
			State:  true,
		},
	}

	go func() {
		for _, msg := range msgs {
			time.Sleep(100 * time.Millisecond)
			msgChan <- msg
		}
		// To kill the forever loop.
		time.Sleep(1 * time.Second)
		close(msgChan)
	}()

	h.EXPECT().HAMessage().AnyTimes().Return(msgChan)
	m.EXPECT().GetBoardByRoom("living").Times(1).Return([]core.Board{
		{
			Addr:    []byte{1, 1, 1},
			Actions: []byte{0x10},
		},
		{
			Addr:    []byte{1, 1, 2},
			Actions: []byte{0x10},
		},
	})
	m.EXPECT().BuzzerOn([]byte{1, 1, 1}, true, 100)
	m.EXPECT().BuzzerOn([]byte{1, 1, 2}, true, 100)
	c.HAMsgHandler()
}

func TestCoreMsgHandler(t *testing.T) {
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

	// Test Motion Events.
	m.EXPECT().Event().AnyTimes().Return(eventChan)
	h.EXPECT().SendMotion("living", true)
	h.EXPECT().SendMotion("hallway-down", true)
	m.EXPECT().GetBoardByAddr([]byte{1, 1, 1}).AnyTimes().Return(&core.Board{Room: "living", Name: "b1"})
	m.EXPECT().GetBoardByAddr([]byte{2, 2, 2}).AnyTimes().Return(&core.Board{Room: "hallway-down", Name: "b1"})

	// TODO: Add tests for other events.

	c.CoreMsgHandler()
}
