package controller

import (
	"fmt"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

type Controller struct {
	httpPort    string                                    // http port.
	core        core.MedusaCore                           // medusa Core struct.
	eventDB     *EventDB                                  // Event log struct.
	ha          HA                                        // HomeAssistant
	eventProcs  map[chan core.Event]func(chan core.Event) // eventProcs is the array of event processors.
	pollInt     time.Duration                             // Poll interval for sensors in seconds.
	pingTimeout time.Duration                             // Ping timeout before a board is considered dead.
}

// NewController returns a new initialized controller.
func NewController(c core.MedusaCore, ha HA, eventDB *EventDB, httpPort string, pollInt, pingTimeout time.Duration) (*Controller, error) {
	// Initialize controller.
	ct := &Controller{
		core:        c,
		httpPort:    httpPort,
		ha:          ha,
		eventDB:     eventDB,
		eventProcs:  make(map[chan core.Event]func(chan core.Event)),
		pollInt:     pollInt,
		pingTimeout: pingTimeout,
	}
	return ct, nil
}

// Startup starts the Controller.
func (c *Controller) Startup() error {
	c.core.StartCore()

	// Startup EventProcessors.
	for c, f := range c.eventProcs {
		go f(c)
	}

	// Connect to MQTT broker for Home Assistant.
	if err := c.ha.Connect(); err != nil {
		return err
	}

	// Startup Event Trigger
	c.SensorDataReq(c.pollInt)

	// Startup Handlers.
	go c.CoreMsgHandler()
	go c.HAMsgHandler()
	if err := c.StartHTTP(); err != nil {
		return fmt.Errorf("failed to start http server %v", err)
	}

	return nil
}

// AddEventProcessor adds event processor routines to be run.
func (c *Controller) AddEventProcessor(f func(chan core.Event)) {
	// Handlers can be used to process events before sending them their way.
	// eg. prepocessing motion events or aggregating data etc.
	// Example below (defined in event_processesor.go)
	//ct.handlers[make(chan core.Event)] = ct.motionRule
	c.eventProcs[make(chan core.Event)] = f
}

// SensorDataReq requests events from sensors at regular interval.
func (c *Controller) SensorDataReq(dur time.Duration) chan bool {
	done := make(chan bool)

	go func() {
		tick := time.NewTicker(dur)
		for {
			select {
			case <-tick.C:
				// TODO: Add other actions.
				for _, brd := range c.core.GetBoardByRoom("all") {
					switch {
					case brd.IsActionCapable(core.ActionTemp):
						if err := c.core.Temp(brd.Addr); err != nil {
							glog.Errorf("Failed to get temp %v", err)
						}

					case brd.IsActionCapable(core.ActionLight):
					}
				}

			case <-done:
				return
			}
		}
	}()

	return done
}

// CoreMsgHandler main loop. This needs to be started separately (not in startup) so its
// accessible in tests.
func (c *Controller) CoreMsgHandler() {
	for {
		event, ok := <-c.core.Event()
		if !ok {
			return
		}
		addr := event.Addr()
		paddr := event.PAddr()
		hwaddr := event.HWAddr()
		tmstmp := time.Now()

		// Handle any relay error events first, since they wont have any addr associated with them.
		if coreError, ok := event.(core.Error); ok {
			glog.Errorf("Event Error - Addr:%v Paddr:%v HWaddr:%v ErrorCode:%v \n",
				core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr), coreError.ErrCode)
			continue
		}

		board := c.core.GetBoardByAddr(addr)

		if board == nil {
			glog.Warningf("Unable to locate board pkt Addr:%v Paddr:%v HWaddr:%v\n",
				core.PP2(addr), core.PP2(paddr), core.PP2(hwaddr))
			continue
		}
		room := board.Room
		name := board.Name

		// Log Event.
		switch f := event.(type) {
		case core.Ping:
			glog.Infof("Event Ping -  Addr:%v\n", core.PP2(addr))
			if e := c.eventDB.LogEvent(EventLog{tmstmp, "ping", 0, room, name, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Temp:
			glog.Infof("Event Temp - Addr:%v temp:%v humi:%v\n", core.PP2(addr), f.Temp, f.Humidity)
			if e := c.eventDB.LogEvent(EventLog{tmstmp, "temperature", f.Temp, room, name, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}
			if e := c.eventDB.LogEvent(EventLog{tmstmp, "humidity", f.Humidity, room, name, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}
			if err := c.ha.SendTemp(board.Room, board.Name, f.Temp, f.Humidity); err != nil {
				glog.Errorf("Failed to send temp event to HA:%v", err)
			}

		case core.Motion:
			glog.Infof("Event Motion addr:%v %v %t\n", core.PP2(addr), room, f.Motion)
			var motion float32
			if f.Motion {
				motion = 1
			}
			if e := c.eventDB.LogEvent(EventLog{tmstmp, "motion", motion, room, name, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}
			if err := c.ha.SendMotion(board.Room, board.Name, f.Motion); err != nil {
				glog.Errorf("Failed to send motion event to HA:%v", err)
			}

		case core.Door:
			glog.Infof("Event Door addr:%v %v %t\n", core.PP2(addr), room, f.Door)
			var door float32
			if f.Door {
				door = 1
			}
			if e := c.eventDB.LogEvent(EventLog{tmstmp, "door", door, room, name, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Volt:
			glog.Infof("Event Volt - addr:%v volts:%v", core.PP2(addr), f.Volt)
			if e := c.eventDB.LogEvent(EventLog{tmstmp, "voltage", f.Volt, room, name, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}

		case core.Light:
			glog.Infof("Event Light - addr:%v light:%v", core.PP2(addr), f.Light)
			if e := c.eventDB.LogEvent(EventLog{tmstmp, "light", f.Light, room, name, addr}); e != nil {
				glog.Errorf("Failed to log to eventDB:%v", e)
			}
		}

		// Send event to all rules engines only after logging into eventDB first.
		for ch, f := range c.eventProcs {
			_ = f
			ch <- event
		}
	}
}

// HAMsgHandler main loop.
func (c *Controller) HAMsgHandler() {
	for {
		msg, ok := <-c.ha.HAMessage()
		if !ok {
			return
		}

		actionID, _ := core.ActionLookup(0, msg.Action)

		brds := c.core.GetBoardByRoom(msg.Room)

		for _, brd := range brds {
			if !brd.IsActionCapable(actionID) {
				continue
			}

			switch actionID {
			case core.ActionBuzzer:
				c.core.BuzzerOn(brd.Addr, msg.State, 100)
			}

		}
	}
}
