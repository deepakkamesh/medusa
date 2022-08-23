package main

import (
	"flag"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

func main() {
	var (
		httpHostPort = flag.String("http_port", ":8080", "host port for http server")
		hostPort     = flag.String("host_port", ":3334", "host port for medusa server")
	)

	flag.Parse()

	// Log flush Routine.
	go func() {
		for {
			glog.Flush()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	glog.Info("Starting Medusa")

	// Startup routines.
	core := core.NewCore(*httpHostPort, *hostPort)
	core.TempInit()
	core.StartPacketHandlers()
	core.StartHTTP()
}
