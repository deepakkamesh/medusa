package controller

import "github.com/deepakkamesh/medusa/controller/core"

type Controller struct {
	httpPort string
	core     *core.Core
}

func NewController(hostPort string, httpPort string) *Controller {
	return &Controller{
		core:     core.NewCore(hostPort),
		httpPort: httpPort,
	}
}

// Startup Controller.
func (c *Controller) Startup() error {
	c.core.StartCore()
	return nil
}
