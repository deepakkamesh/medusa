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
			glog.Infof("Event Temp - Addr:%v Paddr:%v Hwaddr:%v temp:%v humi:%v\n", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()), f.Temp, f.Humidity)

		case core.Motion:
			glog.Infof("Event Motion addr:%v Paddr:%v Hwaddr:%v %t\n", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()), f.Motion)

		case core.Door:
			glog.Infof("Event Door addr:%v Paddr:%v Hwaddr:%v %t\n", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()), f.Door)

		case core.Volt:
			glog.Infof("Event Volt - addr:%v Paddr:%v Hwaddr:%v volts:%v", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()), f.Volt)

		case core.Light:
			glog.Infof("Event Light - addr:%v Paddr:%v Hwaddr:%v light:%v", core.PP2(f.Addr()), core.PP2(f.PAddr()), core.PP2(f.HWAddr()), f.Light)
		}

	}
}
