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
		httpHostPort  = flag.String("http_port", ":8080", "host port for http server")
		hostPort      = flag.String("host_port", ":3334", "host port for medusa server")
		cfgFname      = flag.String("core_conf", "core.cfg.test.json", "config file for core hardware")
		mqttHost      = flag.String("mqtt_host", "homeassistant.local:1883", "hostport for home assistant")
		mqUser        = flag.String("mq_user", "mq", "username for mqtt")
		mqPass        = flag.String("mq_pass", "mqtt", "passwd for mqtt")
		sensorPollInt = flag.Int("poll_int", 90, "poll interval in seconds for sensors")
		pingTimeout   = flag.Int("ping_timeout", 300, "ping timeout in seconds before a board is declared dead")
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

	// HomeAssistant connector.
	ha := controller.NewHA(*mqttHost, *mqUser, *mqPass, core.CoreConfig())

	// Init EventDB
	eventDB, err := controller.NewEventDB()
	if err != nil {
		glog.Fatalf("Failed to init EventDB:%v", err)
	}

	//  Init Controller.
	ctrl, err := controller.NewController(core, ha, eventDB, *httpHostPort,
		time.Duration(*sensorPollInt)*time.Second, time.Duration(*pingTimeout)*time.Second)
	if err != nil {
		glog.Fatalf("Failed init controller %v", err)
	}

	// Add Event Processors.
	ctrl.AddEventProcessor(ctrl.EventProcessorPingCheck)

	// Finally startup controller.
	if err := ctrl.Startup(); err != nil {
		glog.Fatalf("Failed to startup Controller:%v", err)
	}
}
