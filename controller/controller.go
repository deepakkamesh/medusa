package controller

import (
	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

type Controller struct {
	httpPort string
	core     *core.Core
}

func NewController(c *core.Core, httpPort string) *Controller {
	return &Controller{
		core:     c,
		httpPort: httpPort,
	}
}

// Startup Controller.
func (c *Controller) Startup() error {
	c.core.StartCore()
	return nil
}

// Run main loop.
func (c *Controller) Run() {
	for {
		event := <-c.core.Event

		switch f := event.(type) {
		case core.Ping:
			glog.Infof("Event Ping -  Addr:%v Paddr:%v HWaddr:%v\n", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()))

		case core.Temp:
			glog.Infof("Event Temp - Addr:%v paddr:%v hwaddr:%v temp:%v humi:%v\n", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()), f.Temp, f.Humidity)

		case core.Motion:
			glog.Infof("Event Motion addr:%v addr:%v addr:%v\n", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()))

		case core.Volt:
			glog.Infof("Event Volt - Addr:%v paddr:%v hwaddr:%v volts:%v", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()), f.Volt)
		}

	}
}
