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
	ev, _ := controller.NewEventDB()

	c, e := controller.NewController(m, h, ev, ":3344", 1, 1)
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
	ev, _ := controller.NewEventDB()

	c, e := controller.NewController(m, h, ev, ":3344", 1, 1)
	if e != nil {
		t.Errorf("Failed init Controller %v", e)
	}

	boards := []core.Board{
		{
			Room:    "living",
			Name:    "b1",
			Addr:    []byte{1, 2, 3},
			PAddr:   []byte{1, 5, 5, 5, 5},
			Actions: []byte{0x02},
		},
		{
			Room:    "dining",
			Name:    "b1",
			Addr:    []byte{1, 1, 1},
			Actions: []byte{0x01},
		},
	}

	m.EXPECT().GetBoardByRoom("all").MinTimes(1).Return(boards)

	// Test if no availability data does not poll sensors.
	m.EXPECT().Temp([]byte{1, 2, 3}).Times(0)
	d := c.SensorDataReq(10 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	d <- true

	// Test if board is available, sensors are polled.
	_ = ev.LogEvent(controller.EventLog{time.Now(), "availability", 1, "living", "b1", []byte{1, 2, 3}})
	m.EXPECT().Temp([]byte{1, 2, 3}).Times(1)
	d = c.SensorDataReq(10 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	d <- true

	// Test if board is offline, no sensor poll.
	_ = ev.LogEvent(controller.EventLog{time.Now(), "availability", 0, "living", "b1", []byte{1, 2, 3}})
	m.EXPECT().Temp([]byte{1, 2, 3}).Times(0)
	d = c.SensorDataReq(10 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	d <- true

	// Test if board is offline after last poll, no sensor poll and relay reset.
	_ = ev.LogEvent(controller.EventLog{time.Now().Add(5 * time.Millisecond), "availability", 0, "living", "b1", []byte{1, 2, 3}})
	m.EXPECT().Temp([]byte{1, 2, 3}).Times(0)
	m.EXPECT().GetRelaybyPAddr([]byte{1, 5, 5, 5, 5}).Return(&core.Relay{Addr: []byte{4, 4, 4}})
	m.EXPECT().Reset([]byte{4, 4, 4})
	d = c.SensorDataReq(10 * time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	d <- true

}

// Tests is messages from Medusa are routed to HA.
func TestCoreMsgHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockMedusaCore(ctrl)
	h := mocks.NewMockHA(ctrl)
	ev, _ := controller.NewEventDB()

	c, e := controller.NewController(m, h, ev, ":3344", 1, 1)
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

// TestEventProcessorPingCheck checks if the PingCheck event processor sets the board offline if ping is late.
func TestEventProcessorPingCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockMedusaCore(ctrl)
	h := mocks.NewMockHA(ctrl)
	ev, _ := controller.NewEventDB()

	c, e := controller.NewController(m, h, ev, ":3344", 1, 50*time.Millisecond)
	if e != nil {
		t.Errorf("Failed init Controller %v", e)
	}

	m.EXPECT().CoreConfig().AnyTimes().Return(&core.Config{
		Relays: []*core.Relay{},
		Boards: []*core.Board{
			{Room: "living", Name: "b1", Addr: []byte{1, 1, 1}, Actions: []byte{0x01}},
		},
	})

	// Validate that HA offline message only triggers once.
	h.EXPECT().SendAvail("living", "b1", "offline").Times(1)
	ev.LogEvent(controller.EventLog{time.Now(), "ping", 0, "living", "b1", []byte{1, 1, 1}})
	eventChan := make(chan core.Event)
	go c.EventProcessorPingCheck(eventChan)
	// Ping timeout. Send HA offline once.
	time.Sleep(150 * time.Millisecond)
	close(eventChan)
	// Need a delay because close(eventChan) is not guaranteed to be realtime. Flushing DB
	// before goroutine returns will result in another call to offline, breaking test.
	time.Sleep(100 * time.Millisecond) // Delay to allow goroutine to finish.
	ev.PurgeDB()

	// Validate HA online message only triggers once.
	h.EXPECT().SendAvail("living", "b1", "online").Times(1)
	eventChan = make(chan core.Event)
	go c.EventProcessorPingCheck(eventChan)
	ev.LogEvent(controller.EventLog{time.Now().Add(10 * time.Millisecond), "ping", 0, "living", "b1", []byte{1, 1, 1}})
	time.Sleep(55 * time.Millisecond)
	ev.LogEvent(controller.EventLog{time.Now().Add(10 * time.Millisecond), "ping", 0, "living", "b1", []byte{1, 1, 1}})
	time.Sleep(55 * time.Millisecond)
	close(eventChan)
	time.Sleep(100 * time.Millisecond) // Delay to allow goroutine to finish.
	ev.PurgeDB()

	// Validate one HA online message followed by offline once.
	fn := h.EXPECT().SendAvail("living", "b1", "online").Times(1)
	h.EXPECT().SendAvail("living", "b1", "offline").Times(1).After(fn)
	eventChan = make(chan core.Event)
	go c.EventProcessorPingCheck(eventChan)
	ev.LogEvent(controller.EventLog{time.Now().Add(10 * time.Millisecond), "ping", 0, "living", "b1", []byte{1, 1, 1}})
	time.Sleep(110 * time.Millisecond)
	close(eventChan)
	time.Sleep(100 * time.Millisecond) // Delay to allow goroutine to finish.
	ev.PurgeDB()

	// Validate one HA offline message and followed by online.
	fn = h.EXPECT().SendAvail("living", "b1", "offline").Times(1)
	h.EXPECT().SendAvail("living", "b1", "online").Times(1).After(fn)
	eventChan = make(chan core.Event)
	go c.EventProcessorPingCheck(eventChan)
	time.Sleep(150 * time.Millisecond) // timeout exceeded, board offline.
	ev.LogEvent(controller.EventLog{time.Now().Add(10 * time.Millisecond), "ping", 0, "living", "b1", []byte{1, 1, 1}})
	time.Sleep(55 * time.Millisecond) // Board recovered. online.
	close(eventChan)
	time.Sleep(100 * time.Millisecond) // Delay to allow goroutine to finish.
	ev.PurgeDB()

}
