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

	c, e := controller.NewController(m, h, ":3344", 1)
	if e != nil {
		t.Errorf("Failed init Controller %v", e)
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

func TestEventTrigger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockMedusaCore(ctrl)
	h := mocks.NewMockHA(ctrl)

	c, e := controller.NewController(m, h, ":3344", 1)
	if e != nil {
		t.Errorf("Failed init Controller %v", e)
	}

	boards := []core.Board{
		{
			Addr:    []byte{1, 2, 3},
			Actions: []byte{0x02},
		},
		{
			Addr:    []byte{1, 1, 1},
			Actions: []byte{1},
		},
	}

	m.EXPECT().GetBoardByRoom("all").AnyTimes().Return(boards)
	m.EXPECT().Temp([]byte{1, 2, 3}).AnyTimes()
	d := c.EventTrigger(300 * time.Millisecond)
	time.Sleep(1 * time.Second)
	d <- true

}

// Tests is messages from Medusa are routed to HA.
func TestCoreMsgHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockMedusaCore(ctrl)
	h := mocks.NewMockHA(ctrl)

	c, e := controller.NewController(m, h, ":3344", 1)
	if e != nil {
		t.Errorf("Failed init Controller %v", e)
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
	p3 := core.Temp{
		PktInfo: core.PktInfo{
			BoardAddr:    []byte{2, 2, 2},
			PipeAddr:     []byte{},
			HardwareAddr: []byte{},
		},
		Temp:     71.1,
		Humidity: 50.5,
	}

	pkts := []core.Event{p1, p2, p3}

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
	h.EXPECT().SendMotion("living", "b1", true)
	h.EXPECT().SendMotion("hallway-down", "b1", true)
	h.EXPECT().SendTemp("hallway-down", "b1", float32(71.1), float32(50.5))
	m.EXPECT().GetBoardByAddr([]byte{1, 1, 1}).AnyTimes().Return(&core.Board{Room: "living", Name: "b1"})
	m.EXPECT().GetBoardByAddr([]byte{2, 2, 2}).AnyTimes().Return(&core.Board{Room: "hallway-down", Name: "b1"})

	c.CoreMsgHandler()
}
