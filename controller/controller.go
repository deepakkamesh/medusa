package controller

import (
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

type Controller struct {
	httpPort string                                    // http port.
	core     core.MedusaCore                           // medusa Core struct.
	eventDB  *EventDB                                  // Event log struct.
	rules    map[chan core.Event]func(chan core.Event) // rules is the array of rules functions.
}

func NewController(c core.MedusaCore, httpPort string) (*Controller, error) {

	// Create new EventDB to hold recent events.
	eventDB, err := NewEventDB()
	if err != nil {
		return nil, err
	}

	// Initialize controller.
	ct := &Controller{
		core:     c,
		httpPort: httpPort,
		eventDB:  eventDB,
		rules:    make(map[chan core.Event]func(chan core.Event)),
	}

	// Add rule engines.
	ct.rules[make(chan core.Event)] = ct.motionRule

	return ct, nil
}

// Startup Controller.
func (c *Controller) Startup() error {
	c.core.StartCore()

	// Startup rules engines.
	for c, f := range c.rules {
		go f(c)
	}
	return nil
}

// Run main loop.
func (c *Controller) Run() {
	for {
		event, ok := <-c.core.Event()
		if !ok {
			return
		}
		addr := event.Addr()
		paddr := event.PAddr()
		hwaddr := event.HWAddr()
		tmstmp := time.Now()

		board := c.core.GetBoardByAddr(addr)
		room := "unknown"
		if board != nil {
			room = board.Room
		}

		// Log Event.
		switch f := event.(type) {
		case core.Ping:
			glog.Infof("Event Ping -  Addr:%v Paddr:%v HWaddr:%v\n", core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr))
			if e := c.eventDB.LogEvent(eventLog{tmstmp, "ping", 0, room, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Temp:
			glog.Infof("Event Temp - Addr:%v Paddr:%v Hwaddr:%v temp:%v humi:%v\n", core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr), f.Temp, f.Humidity)
			if e := c.eventDB.LogEvent(eventLog{tmstmp, "temp", f.Temp, room, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}
			if e := c.eventDB.LogEvent(eventLog{tmstmp, "humidity", f.Humidity, room, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Motion:
			glog.Infof("Event Motion addr:%v Paddr:%v Hwaddr:%v %t\n", core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr), f.Motion)
			var motion float32
			if f.Motion {
				motion = 1
			}
			if e := c.eventDB.LogEvent(eventLog{tmstmp, "motion", motion, room, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Door:
			glog.Infof("Event Door addr:%v Paddr:%v Hwaddr:%v %t\n", core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr), f.Door)
			var door float32
			if f.Door {
				door = 1
			}
			if e := c.eventDB.LogEvent(eventLog{tmstmp, "door", door, room, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Volt:
			glog.Infof("Event Volt - addr:%v Paddr:%v Hwaddr:%v volts:%v", core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr), f.Volt)
			if e := c.eventDB.LogEvent(eventLog{tmstmp, "volt", f.Volt, room, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Light:
			glog.Infof("Event Light - addr:%v Paddr:%v Hwaddr:%v light:%v", core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr), f.Light)
			if e := c.eventDB.LogEvent(eventLog{tmstmp, "light", f.Light, room, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}
		}

		// Send event to all rules engines only after logging into eventDB first.
		for ch, f := range c.rules {
			_ = f
			ch <- event
		}

	}
}
