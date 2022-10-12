package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -destination=./mocks/ha_mock.go -package=mocks github.com/deepakkamesh/medusa/controller HA

type HA interface {
	Connect() error
	SendSensorConfig(c *core.Config) error
	SendSensorData(topic string, pri byte, retain bool, msg string) error
}

type MQBinarySensorConfig struct {
	Name        string            `json:"name"`
	ObjectID    string            `json:"object_id"`
	DeviceClass string            `json:"device_class"`
	StateTopic  string            `json:"state_topic"`
	UniqueID    string            `json:"unique_id"`
	Device      map[string]string `json:"device"`
}

type HomeAssistant struct {
	mqttHost   string
	mqttClient mqtt.Client
}

func NewHA(mqttHost string) *HomeAssistant {
	return &HomeAssistant{
		mqttHost: mqttHost,
	}
}

func (m *HomeAssistant) Connect() error {
	options := mqtt.NewClientOptions()
	options.AddBroker("mqtt://" + m.mqttHost)
	options.OnConnectionLost = m.mqttConnectionLost
	options.SetUsername("mq")
	options.SetPassword("mqtt")

	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT connect error : %v", token.Error())
	}

	m.mqttClient = client
	return nil
}

func (m *HomeAssistant) mqttConnectionLost(client mqtt.Client, err error) {
	glog.Errorf("MQTT Connection Lost: %v", err)
	m.mqttClient = nil
}

// SendSensorConfig sends sensors via auto discovery.
func (m *HomeAssistant) SendSensorConfig(c *core.Config) error {

	if m.mqttClient == nil {
		return fmt.Errorf("mqtt broker %v not connected", m.mqttHost)
	}

	binarySensors := []byte{core.ActionMotion} //TODO move elsewhere.

	for _, brd := range c.Boards {

		for _, act := range brd.Actions {
			actionStr := core.ActionFriendlyName(act)
			if !slices.Contains(binarySensors, act) {
				continue
			}
			sensorConfig := MQBinarySensorConfig{
				Name:        actionStr,
				ObjectID:    fmt.Sprintf("%v_%v_%v", brd.Room, brd.Name, actionStr),
				DeviceClass: DeviceClass(act),
				StateTopic:  fmt.Sprintf("homeassistant/%v_%v_%v/state", brd.Room, brd.Name, actionStr),
				UniqueID:    fmt.Sprintf("%v_%v_%v", brd.Room, brd.Name, actionStr),
				Device: map[string]string{
					"identifiers":    core.PP2(brd.Addr),
					"suggested_area": brd.Room,
					"name":           brd.Name,
				},
			}

			// Marshall string.
			a, err := json.Marshal(sensorConfig)
			if err != nil {
				return err
			}
			configTopic := fmt.Sprintf("homeassistant/binary_sensor/%v_%v_%v/config", brd.Room, brd.Name, actionStr)

			token := m.mqttClient.Publish(configTopic, 0, false, fmt.Sprintf("%s", a))
			token.Wait()
			time.Sleep(time.Second)

		}
	}

	return nil
}

// Sends Data on the specified topic.
func (m *HomeAssistant) SendSensorData(topic string, pri byte, retain bool, msg string) error {
	if m.mqttClient == nil {
		return fmt.Errorf("mqtt broker %v not connected", m.mqttHost)
	}

	token := m.mqttClient.Publish(topic, pri, retain, msg)
	token.Wait()
	time.Sleep(time.Second)
	return nil
}

// Returns HA device class for each ActionID.
func DeviceClass(actionID byte) string {
	m := map[byte]string{
		0x01: "motion",
		0x02: "temperature",
		0x03: "light",
	}

	if v, ok := m[actionID]; ok {
		return v
	}
	return ""
}
