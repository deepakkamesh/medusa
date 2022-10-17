package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/deepakkamesh/medusa/controller/core"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"golang.org/x/exp/slices"
)

//go:generate mockgen -destination=./mocks/ha_mock.go -package=mocks github.com/deepakkamesh/medusa/controller HA

// Mocks for mqtt package.
//go:generate mockgen -destination=./mocks/mqtt_mock.go -package=mocks  github.com/eclipse/paho.mqtt.golang Client
//go:generate mockgen -destination=./mocks/mqtt_mock_2.go -package=mocks  github.com/eclipse/paho.mqtt.golang Token
//go:generate mockgen -destination=./mocks/mqtt_mock_3.go -package=mocks  github.com/eclipse/paho.mqtt.golang Message

const topicSubscribe string = "giant/+/+/set"

// HA represents HomeAssistant Interface.
type HA interface {
	Connect() error
	HAMessage() <-chan HAMsg
	SendMotion(room string, motion bool) error
	SendSensorConfig(clean bool) error
}

// MQState represents the state json from HA MQ message.
type MQState struct {
	State string `json:"state"`
}

// MQBinarySensorConfig represents the HA Binary Sensor.
type MQBinarySensorConfig struct {
	Name        string            `json:"name"`
	ObjectID    string            `json:"object_id"`
	DeviceClass string            `json:"device_class"`
	StateTopic  string            `json:"state_topic"`
	UniqueID    string            `json:"unique_id"`
	Device      map[string]string `json:"device"`
}

// MQSirenConfig represents the HA Binary Sensor.
type MQSirenConfig struct {
	Name         string            `json:"name"`
	ObjectID     string            `json:"object_id"`
	CommandTopic string            `json:"command_topic"`
	UniqueID     string            `json:"unique_id"`
	Device       map[string]string `json:"device"`
	PayloadOn    string            `json:"payload_on"`
	PayloadOff   string            `json:"payload_off"`
}

// Struct to hold message from HA.
type HAMsg struct {
	MQMsg  mqtt.Message // Full topic.
	Room   string
	Action string
	State  bool
}

type HomeAssistant struct {
	mqttHost   string
	user       string
	passwd     string
	MQTTClient mqtt.Client  // MqTT client. Exportable for testing.
	HAMsg      chan HAMsg   // HA received message. Exportable for testing.
	CoreCfg    *core.Config // Core Config.
}

func NewHA(mqttHost string, user string, pass string, cfg *core.Config) *HomeAssistant {
	return &HomeAssistant{
		mqttHost: mqttHost,
		user:     user,
		passwd:   pass,
		HAMsg:    make(chan HAMsg),
		CoreCfg:  cfg,
	}
}

// Connect connects to MQ Server.
func (m *HomeAssistant) Connect() error {
	options := mqtt.NewClientOptions()
	options.AddBroker("mqtt://" + m.mqttHost)
	options.OnConnectionLost = m.mqttConnLostHandler
	options.SetDefaultPublishHandler(m.MQTTPubHandler)
	options.OnConnect = m.mqttConnectHandler
	options.SetUsername(m.user)
	options.SetPassword(m.passwd)

	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT connect error : %v", token.Error())
	}

	m.MQTTClient = client

	token = m.MQTTClient.Subscribe(topicSubscribe, 0, nil)
	token.Wait()

	return nil
}

func (m *HomeAssistant) HAMessage() <-chan HAMsg {
	return m.HAMsg
}

// Parse the HA message and pass back to controller.
func (m *HomeAssistant) MQTTPubHandler(client mqtt.Client, msg mqtt.Message) {
	// Parse topic to find room. // Format giant/<room>/<device_type>/set.
	data := strings.Split(msg.Topic(), "/")

	// Parse state.
	st := MQState{}
	if err := json.Unmarshal(msg.Payload(), &st); err != nil {
		glog.Warningf("Unable to unmarshall HA state from MQTT msg:%v", err)
		return
	}

	state, err := strconv.ParseBool(st.State)
	if err != nil {
		glog.Warningf("Unable to parsebool:%v", err)
	}

	b := HAMsg{
		MQMsg:  msg,
		Room:   data[1],
		Action: data[2],
		State:  state,
	}
	m.HAMsg <- b
}

// SendMotion sends motion event to HA.
func (m *HomeAssistant) SendMotion(room string, motion bool) error {
	state := "OFF"
	if motion {
		state = "ON"
	}
	topic := fmt.Sprintf("giant/%v/motion/state", room)
	return m.sendSensorData(topic, 0, false, state)
}

// SendTemp sends Temp and Humidity to HA.
func (m *HomeAssistant) SendTemp(room string, temp, humidity float32) error {
	return nil
}

// Sends Data on the specified topic.
func (m *HomeAssistant) sendSensorData(topic string, pri byte, retain bool, msg string) error {
	if m.MQTTClient == nil {
		return fmt.Errorf("mqtt broker %v not connected", m.mqttHost)
	}

	token := m.MQTTClient.Publish(topic, pri, retain, msg)
	token.Wait()
	time.Sleep(time.Second)
	return nil
}

func (m *HomeAssistant) mqttConnectHandler(client mqtt.Client) {
	glog.Infof("MQTT connected:%v", m.mqttHost)

	// Send configuration after connection successful.
	m.SendSensorConfig(false)
}

func (m *HomeAssistant) mqttConnLostHandler(client mqtt.Client, err error) {
	glog.Warningf("MQTT Connection Lost. Retrying: %v", err)
	m.MQTTClient = nil
	time.Sleep(1 * time.Second)
	// Attempt to reconnect.
	if err := m.Connect(); err != nil {
		glog.Errorf("Failed reconnecting to mqtt:%v", err)
	}
}

// SendSensorConfig sends sensors via auto discovery.
func (m *HomeAssistant) SendSensorConfig(clean bool) error {

	if m.MQTTClient == nil {
		return fmt.Errorf("mqtt broker %v not connected", m.mqttHost)
	}

	binarySensors := []byte{core.ActionMotion, core.ActionDoor} //TODO move elsewhere.
	sirens := []byte{core.ActionBuzzer}

	for _, brd := range m.CoreCfg.Boards {
		for _, actionID := range brd.Actions {
			_, actionStr := core.ActionLookup(actionID, "")

			mqttType := ""
			configMsg := ""

			// Common stuff.
			name := brd.Room + " " + actionStr
			objectID := fmt.Sprintf("%v_%v_%v", brd.Room, brd.Name, actionStr)
			uniqueID := fmt.Sprintf("%v_%v_%v", brd.Room, brd.Name, actionStr)
			stateTopic := fmt.Sprintf("giant/%v/%v/state", brd.Room, actionStr)
			commandTopic := fmt.Sprintf("giant/%v/%v/set", brd.Room, actionStr)
			device := map[string]string{
				"identifiers":    core.PP2(brd.Addr),
				"suggested_area": brd.Room,
				//		"name":           brd.Name,
				"manufacturer": "Medusa",
			}

			switch {
			case slices.Contains(binarySensors, actionID):
				sensorConfig := MQBinarySensorConfig{
					Name:        name,
					ObjectID:    objectID,
					DeviceClass: DeviceClass(actionID),
					StateTopic:  stateTopic,
					UniqueID:    uniqueID,
					Device:      device,
				}
				// Marshall string.
				a, err := json.Marshal(sensorConfig)
				if err != nil {
					return err
				}
				configMsg = string(a)
				mqttType = "binary_sensor"

			case slices.Contains(sirens, actionID):
				deviceConfig := MQSirenConfig{
					Name:         name,
					ObjectID:     objectID,
					CommandTopic: commandTopic,
					UniqueID:     uniqueID,
					PayloadOn:    "true",
					PayloadOff:   "false",
					Device:       device,
				}
				// Marshall string.
				a, err := json.Marshal(deviceConfig)
				if err != nil {
					return err
				}
				configMsg = string(a)
				mqttType = "siren"
			}

			// Send the discovery message.
			if configMsg != "" {
				configTopic := fmt.Sprintf("homeassistant/%v/%v_%v_%v/config", mqttType, brd.Room, brd.Name, actionStr)
				// Send empty message to delete HA entity.
				if clean {
					configMsg = ""
				}
				token := m.MQTTClient.Publish(configTopic, 0, false, configMsg)
				token.Wait()
				time.Sleep(1 * time.Second)
			}
		}
	}

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
