package controller

import (
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

/*
func (c *Controller) motionRule(in chan core.Event) {

	adj := make(map[string][]string)

	// Define adjancey for rooms
	adj["living"] = []string{"hallway-down"}
	adj["hallway-down"] = []string{"living"}

	for {
		ev := <-in
		if reflect.TypeOf(ev).String() != "core.Motion" {
			continue
		}
		board := c.core.GetBoardByAddr(ev.Addr())
		if board == nil { // board not found.
			continue
		}
		m, ok := ev.(core.Motion)
		if !ok {
			glog.Error("Cast of event to core.Motion failed")
		}

		if err := c.ha.SendMotion(board.Room, board.Name, m.Motion); err != nil {
			glog.Warningf("Failed to send motion event to HA:%v", err)
		}

		// Check if adjancey motion detected.
				highC := false
				for _, v := range adj[room] {
					log, e := c.eventDB.GetEvent("motion", v, 150*time.Millisecond)
					if e != nil {
						glog.Errorf("Failed to query eventDB:%v", e)
					}

					if len(log) > 0 {
						highC = true
						break
					}
				}

				if highC {
					fmt.Println("high confidence motion", room)
					return
				}
				fmt.Println("low confidence motion", room)

	}
}*/

// EventProcessorPingCheck calls HA and sets entities offline if ping is not received within timeout.
func (c *Controller) EventProcessorPingCheck(in chan core.Event) {
	tick := time.NewTicker(c.pingTimeout)
	for {
		select {
		case _, ok := <-in:
			if !ok {
				return
			}

		case <-tick.C:
			cfg := c.core.CoreConfig()
			for _, b := range cfg.Boards {

				pings, err := c.eventDB.GetEvent("ping", b.Room, b.Name, c.pingTimeout)
				if err != nil {
					glog.Errorf("Error reading from events log DB:%v", err)
					continue
				}

				lastAvail, err := c.eventDB.GetLastEvent("availability", b.Room, b.Name, 1)
				if err != nil {
					glog.Errorf("Error reading from events log DB:%v", err)
					continue
				}

				switch {
				// Got Pings, set availability to online if no previous availability event or previous availability was offline.
				case len(pings) > 0:
					if len(lastAvail) == 0 || lastAvail[0].Value == 0 {
						if err := c.eventDB.LogEvent(EventLog{time.Now(), "availability", 1, b.Room, b.Name, b.Addr}); err != nil {
							glog.Errorf("Failed to log to eventDB:%v", err)
						}
						c.ha.SendAvail(b.Room, b.Name, payloadOnline)
					}

				// No pings received within timeout, set availability to offline if no previous availability event or previous availability was online.
				case len(pings) == 0:
					if len(lastAvail) == 0 || lastAvail[0].Value == 1 {
						if err := c.eventDB.LogEvent(EventLog{time.Now(), "availability", 0, b.Room, b.Name, b.Addr}); err != nil {
							glog.Errorf("Failed to log to eventDB:%v", err)
						}
						c.ha.SendAvail(b.Room, b.Name, payloadOffline)
					}
				}
			}
		}
	}
}
