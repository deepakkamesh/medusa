package core

type Core struct {
	hostPort     string
	httpHostPort string
	relays       map[string]*Relay
}

func NewCore(httpHostPort string, hostPort string) *Core {
	return &Core{
		hostPort:     hostPort,
		httpHostPort: httpHostPort,
		relays:       make(map[string]*Relay),
	}
}

func (c *Core) TempInit() {
	c.relays["0"] = &Relay{}
}
