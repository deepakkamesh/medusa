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

				events, err := c.eventDB.GetEvent("ping", b.Room, b.Name, c.pingTimeout)
				if err != nil {
					glog.Errorf("Error reading from events log DB:%v", err)
					continue
				}

				if len(events) != 0 {
					continue
				}

				c.ha.SendAvail(b.Room, b.Name, "offline")
			}
		}
	}
}
