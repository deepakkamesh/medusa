package main

import (
	"flag"
	"time"

	"github.com/deepakkamesh/medusa/controller"
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

	// Startup processes.
	ctrl := controller.NewController(*hostPort, *httpHostPort)
	if err := ctrl.Startup(); err != nil {
		glog.Fatalf("Failed to startup Controller:%v", err)
	}

	if err := ctrl.StartHTTP(); err != nil {
		glog.Fatalf("Failed to start http server %v", err)
	}
}
