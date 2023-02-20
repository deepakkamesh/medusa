package controller

import (
	"encoding/json"
	"fmt"
	"math"
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

var (
	templTopicState   string = "giant/%v/%v/%v/state"             // room/board_name/action.
	templTopicConfig  string = "homeassistant/%v/%v_%v_%v/config" //Sensor type/room/board_name/action
	templTopicCommand string = "giant/%v/%v/%v/set"               // room/board_name/action
	templEntityName   string = "%v %v"
	templEntityUniqID string = "%v_%v_%v"
	templEntityObjID  string = "%v_%v_%v"
	templDeviceName   string = "%v_%v"
)

// HA represents HomeAssistant Interface.
type HA interface {
	Connect() error
	HAMessage() <-chan HAMsg
	SendMotion(room string, name string, motion bool) error
	SendTemp(room string, name string, temp, humidity float32) error
	SendSensorConfig(clean bool) error
}

// MQState represents the state json from HA MQ message.
type MQState struct {
	State string `json:"state"`
}

// MQTempHumidity represents a temp, humidity message to HA.
type MQTempHumidity struct {
	Temp     float32 `json:"temperature"`
	Humidity float32 `json:"humidity"`
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

// MQSensorConfig represents a HA Sensor.
type MQSensorConfig struct {
	Name        string            `json:"name"`
	ObjectID    string            `json:"object_id"`
	DeviceClass string            `json:"device_class"`
	StateTopic  string            `json:"state_topic"`
	UniqueID    string            `json:"unique_id"`
	Device      map[string]string `json:"device"`
	ValueTempl  string            `json:"value_template"`
	UnitMeasure string            `json:"unit_of_measurement"`
	ForceUpdate bool              `json:"force_update"`
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
	MQMsg     mqtt.Message // Full topic.
	Room      string
	BoardName string
	Action    string
	State     bool
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
		MQMsg:     msg,
		Room:      data[1],
		BoardName: data[2],
		Action:    data[2],
		State:     state,
	}
	m.HAMsg <- b
}

// mqttConnectHandler is the callback once the mqtt connection is established.
func (m *HomeAssistant) mqttConnectHandler(client mqtt.Client) {
	glog.Infof("MQTT connected:%v", m.mqttHost)
}

// mqttConnLostHandler is the callback when connection to mqtt server is lost.
func (m *HomeAssistant) mqttConnLostHandler(client mqtt.Client, err error) {
	glog.Warningf("MQTT Connection Lost. Retrying: %v", err)
	m.MQTTClient = nil
	time.Sleep(1 * time.Second)

	// Try every 5 seconds for 1 hr.
	for i := 0; i < 720; i++ {
		if err := m.Connect(); err == nil {
			glog.Infof("MQTT reconnected after %v attempts", i)
			return
		}
		glog.Warningf("MQTT Connection Lost. Retrying attempt #%v:%v", i, err)
		time.Sleep(5 * time.Second)
	}

	glog.Fatalf("Giving up on MQTT connection. Exiting..")
}

// SendMotion sends motion event to HA.
func (m *HomeAssistant) SendMotion(room string, name string, motion bool) error {
	state := "OFF"
	if motion {
		state = "ON"
	}

	_, actStr := core.ActionLookup(core.ActionMotion, "")
	if actStr == "" {
		return fmt.Errorf("action string not found for action %v", core.ActionMotion)
	}
	topic := fmt.Sprintf(templTopicState, room, name, actStr)
	return m.sendSensorData(topic, 0, false, state)
}

// SendTemp sends Temp and Humidity to HA.
func (m *HomeAssistant) SendTemp(room, name string, temp, humidity float32) error {
	te := float32(math.Floor(float64(temp)*10) / 10)
	hu := float32(math.Floor(float64(humidity)*10) / 10)
	t := MQTempHumidity{te, hu}

	a, e := json.Marshal(t)
	if e != nil {
		return e
	}

	_, actStr := core.ActionLookup(core.ActionTemp, "")
	if actStr == "" {
		return fmt.Errorf("action string not found for action %v", core.ActionTemp)
	}

	topic := fmt.Sprintf(templTopicState, room, name, actStr)
	return m.sendSensorData(topic, 0, false, string(a))
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

// SendSensorConfig sends sensors via auto discovery.
func (m *HomeAssistant) SendSensorConfig(clean bool) error {

	if m.MQTTClient == nil {
		return fmt.Errorf("mqtt broker %v not connected", m.mqttHost)
	}

	binarySensors := []byte{core.ActionMotion, core.ActionDoor}
	sirens := []byte{core.ActionBuzzer}
	sensors := []byte{core.ActionTemp}

	for _, brd := range m.CoreCfg.Boards {
		for _, actionID := range brd.Actions {

			_, actionStr := core.ActionLookup(actionID, "")
			if actionStr == "" {
				glog.Warningf("Action not found for value %v", actionID)
				continue
			}

			// Common config for HA entity.
			name := fmt.Sprintf(templEntityName, brd.Room, actionStr)
			objectID := fmt.Sprintf(templEntityObjID, brd.Room, brd.Name, actionStr)
			uniqueID := fmt.Sprintf(templEntityUniqID, brd.Room, brd.Name, actionStr)
			stateTopic := fmt.Sprintf(templTopicState, brd.Room, brd.Name, actionStr)
			commandTopic := fmt.Sprintf(templTopicCommand, brd.Room, brd.Name, actionStr)
			device := map[string]string{
				"identifiers":    core.PP2(brd.Addr),
				"suggested_area": brd.Room,
				"name":           fmt.Sprintf(templDeviceName, brd.Room, brd.Name),
				"manufacturer":   "Medusa",
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

				if err := m.packAndSendEntityDiscovery(clean, sensorConfig, "binary_sensor", fmt.Sprintf(templTopicConfig, "binary_sensor", brd.Room, brd.Name, actionStr)); err != nil {
					return err
				}

			case slices.Contains(sensors, actionID):
				sensorConfig := MQSensorConfig{
					Name:        name,
					ObjectID:    objectID,
					DeviceClass: DeviceClass(actionID),
					StateTopic:  stateTopic,
					UniqueID:    uniqueID,
					Device:      device,
					ForceUpdate: true,
				}

				// Set value template since Medusa sends both temp and humidity together.
				if actionID == core.ActionTemp {
					sensorConfig.ValueTempl = "{{ value_json.temperature }}"
					sensorConfig.UnitMeasure = "°F"
				}

				if err := m.packAndSendEntityDiscovery(clean, sensorConfig, "sensor", fmt.Sprintf(templTopicConfig, "sensor", brd.Room, brd.Name, actionStr)); err != nil {
					return err
				}

				// Handle special case of humidity entity which needs to be created  in HA for every temperature device in Medusa 'cause  medusa a temperature action retrieves both temp and humidity.
				if actionID == core.ActionTemp {
					sensorConfig := MQSensorConfig{
						Name:        fmt.Sprintf(templEntityName, brd.Room, "humidity"),
						ObjectID:    fmt.Sprintf(templEntityObjID, brd.Room, brd.Name, "humidity"),
						DeviceClass: "humidity",
						StateTopic:  stateTopic, // State topic is shared with temperature entity.
						UniqueID:    fmt.Sprintf(templEntityUniqID, brd.Room, brd.Name, "humidity"),
						Device:      device,
						ValueTempl:  "{{ value_json.humidity }}",
						UnitMeasure: "%",
					}

					if err := m.packAndSendEntityDiscovery(clean, sensorConfig, "sensor", fmt.Sprintf(templTopicConfig, "sensor", brd.Room, brd.Name, "humidity")); err != nil {
						return err
					}
				}

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

				if err := m.packAndSendEntityDiscovery(clean, deviceConfig, "siren", fmt.Sprintf(templTopicConfig, "siren", brd.Room, brd.Name, actionStr)); err != nil {
					return err
				}
			} // End of switch.
		}
	}

	return nil
}

func (m *HomeAssistant) packAndSendEntityDiscovery(clean bool, v interface{}, mqttType string, discTopic string) error {
	// Marshall string.
	a, err := json.Marshal(v)
	if err != nil {
		return err
	}
	discMsg := string(a)

	// Send empty message to delete HA entity.
	if clean {
		discMsg = ""
	}
	token := m.MQTTClient.Publish(discTopic, 0, true, discMsg)
	token.Wait()
	time.Sleep(300 * time.Millisecond)
	glog.Infof("%v:%v", discTopic, discMsg)

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
