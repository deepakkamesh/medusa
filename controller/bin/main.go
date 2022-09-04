package main

import (
	"flag"
	"time"

	"github.com/deepakkamesh/medusa/controller"
	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/golang/glog"
)

func main() {
	var (
		httpHostPort = flag.String("http_port", ":8080", "host port for http server")
		hostPort     = flag.String("host_port", ":3334", "host port for medusa server")
		cfgFname     = flag.String("core_conf", "core.cfg.test.json", "config file for core hardware")
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

	// Init Core.
	core, err := core.NewCore(*hostPort, *cfgFname)
	if err != nil {
		glog.Fatalf("Failed to init core:%v", err)
	}

	//  Init Controller.
	ctrl := controller.NewController(core, *httpHostPort)
	if err := ctrl.Startup(); err != nil {
		glog.Fatalf("Failed to startup Controller:%v", err)
	}

	go ctrl.Run()

	if err := ctrl.StartHTTP(); err != nil {
		glog.Fatalf("Failed to start http server %v", err)
	}
}
