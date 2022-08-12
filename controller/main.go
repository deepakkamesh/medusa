package main

import (
	"flag"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

func main() {
	var (
		httpHostPort = flag.String("http_port", ":8080", "host port for http server")
		hostPort     = flag.String("host_port", ":6999", "host port for medusa server")
	)
	flag.Parse()
	glog.Info("Starting Medusa")
	core := core.NewCore(*httpHostPort, *hostPort)

	core.StartHTTP()
}
