package controller

import (
	"fmt"
	"reflect"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

func (c *Controller) motionRule(in chan core.Event) {
	/*
		// Create a new pushover app with a token
		app := pushover.New("uQiRzpo4DXghDmr9QzzfQu27cmVRsG")

		// Create a new recipient
		recipient := pushover.NewRecipient("gznej3rKEVAvPUxu9vvNnqpmZpokzF")

		// Create the message to send
		message := pushover.NewMessage("Hello !")

		// Send the message to the recipient
		response, err := app.SendMessage(message, recipient)
		if err != nil {
			log.Panic(err)
		}*/

	adj := make(map[string][]string)

	// Define adjancey for rooms
	adj["living"] = []string{"hallway-down"}
	adj["hallway-down"] = []string{"living"}

	highC := false
	for {
		ev := <-in
		if reflect.TypeOf(ev).String() != "core.Motion" {
			continue
		}
		continue
		room := c.core.GetBoardByAddr(ev.Addr()).Room

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
}
