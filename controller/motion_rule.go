package controller

import (
	"fmt"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

func (c *Controller) motionRule(in chan core.Event) {
	for {
		ev := <-in
		room := c.core.GetBoardByAddr(ev.Addr()).Room

		log, e := c.eventDB.GetEvent("motion", "living", 150*time.Millisecond)
		if e != nil {
			glog.Errorf("Failed to query eventDB:%v", e)
		}
		if len(log) > 0 {
			fmt.Println("high confidence motion", room)
			continue
		}
		fmt.Println("low confidence motion", room)
	}
}
