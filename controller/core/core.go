package core

type Core struct {
	hostPort     string
	httpHostPort string
	tempChan     chan []byte
}

func NewCore(httpHostPort string, hostPort string) *Core {
	return &Core{
		hostPort:     hostPort,
		httpHostPort: httpHostPort,
		tempChan:     make(chan []byte, 0),
	}
}
