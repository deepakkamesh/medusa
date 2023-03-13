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

// Topics to listen from HA.
const topicSubscribe string = "giant/+/+/+/set"

var (
	templTopicState   string = "giant/%v/%v/%v/state"             // room/board_name/action.
	templTopicAvail   string = "giant/%v/%v/avail"                // room/board_name/action.
	templTopicConfig  string = "homeassistant/%v/%v_%v_%v/config" //Sensor type/room/board_name/action
	templTopicCommand string = "giant/%v/%v/%v/set"               // room/board_name/action
	templEntityName   string = "%v %v"
	templEntityUniqID string = "%v_%v_%v"
	templEntityObjID  string = "%v_%v_%v"
	templDeviceName   string = "%v_%v"
	payloadOnline     string = "online"
	payloadOffline    string = "offline"
)

// HA represents HomeAssistant Interface.
type HA interface {
	Connect() error
	HAMessage() <-chan HAMsg
	SendAvail(room string, name string, avail string) error
	SendMotion(room string, name string, motion bool) error
	SendTemp(room string, name string, temp, humidity float32) error
	SendLight(room, name string, light float32) error
	SendDoor(room, name string, open bool) error
	SendVolt(room, name string, light float32) error
	SendMQTTDiscoveryConfig(clean bool) error
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
	AvailTopic  string            `json:"availability_topic"`
}

// MQSensorConfig represents a HA Sensor.
type MQSensorConfig struct {
	Name        string            `json:"name"`
	ObjectID    string            `json:"object_id"`
	DeviceClass string            `json:"device_class"`
	StateTopic  string            `json:"state_topic"`
	UniqueID    string            `json:"unique_id"`
	Device      map[string]string `json:"device"`
	ValueTempl  string            `json:"value_template,omitempty"`
	UnitMeasure string            `json:"unit_of_measurement,omitempty"`
	ForceUpdate bool              `json:"force_update"`
	AvailTopic  string            `json:"availability_topic"`
}

// MQSirenConfig represents the HA Siren.
type MQSirenConfig struct {
	Name         string            `json:"name"`
	ObjectID     string            `json:"object_id"`
	CommandTopic string            `json:"command_topic"`
	UniqueID     string            `json:"unique_id"`
	Device       map[string]string `json:"device"`
	PayloadOn    string            `json:"payload_on"`
	PayloadOff   string            `json:"payload_off"`
	AvailTopic   string            `json:"availability_topic"`
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
		Action:    data[3],
		State:     state,
	}
	glog.Infof("Got HA message - Room: %v Board:%v Action:%v State: %v", b.Room, b.BoardName, b.Action, b.State)
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

// SetAvail sets the availability on all the entities on the dev board.
func (m *HomeAssistant) SendAvail(room string, name string, avail string) error {

	for _, brd := range m.CoreCfg.Boards {
		if brd.Name != name {
			continue
		}

		// Send availability topic to offline for all entities for the device. This topic
		// is shared between all entities on the device.
		topic := fmt.Sprintf(templTopicAvail, room, name)
		return m.sendSensorData(topic, 0, false, avail)
	}
	return nil
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

// SendLight sends the light level in lux to HA.
func (m *HomeAssistant) SendLight(room, name string, light float32) error {
	le := float32(math.Floor(float64(light)*100) / 100)

	_, actStr := core.ActionLookup(core.ActionLight, "")
	if actStr == "" {
		return fmt.Errorf("action string not found for action %v", core.ActionLight)
	}

	topic := fmt.Sprintf(templTopicState, room, name, actStr)
	return m.sendSensorData(topic, 0, false, fmt.Sprintf("%v", le))
}

// SendDoor sends contact sensor to HA.
func (m *HomeAssistant) SendDoor(room, name string, open bool) error {
	state := "OFF"
	if open {
		state = "ON"
	}

	_, actStr := core.ActionLookup(core.ActionDoor, "")
	if actStr == "" {
		return fmt.Errorf("action string not found for action %v", core.ActionDoor)
	}
	topic := fmt.Sprintf(templTopicState, room, name, actStr)
	return m.sendSensorData(topic, 0, false, state)
}

// SendVolt sends the voltage reading of board battery to HA.
func (m *HomeAssistant) SendVolt(room, name string, volt float32) error {
	v := float32(math.Floor(float64(volt)*100) / 100)

	_, actStr := core.ActionLookup(core.ActionVolt, "")
	if actStr == "" {
		return fmt.Errorf("action string not found for action %v", core.ActionVolt)
	}

	topic := fmt.Sprintf(templTopicState, room, name, actStr)
	return m.sendSensorData(topic, 0, false, fmt.Sprintf("%v", v))
}

// Sends Data on the specified topic.
func (m *HomeAssistant) sendSensorData(topic string, pri byte, retain bool, msg string) error {
	if m.MQTTClient == nil {
		return fmt.Errorf("mqtt broker %v not connected", m.mqttHost)
	}

	glog.Infof("Publish to mqtt-HA %v : %v", topic, msg)
	token := m.MQTTClient.Publish(topic, pri, retain, msg)
	token.Wait()
	return nil
}

// SendMQTTDiscoveryConfig sends sensors and entities via auto discovery.
func (m *HomeAssistant) SendMQTTDiscoveryConfig(clean bool) error {

	if m.MQTTClient == nil {
		return fmt.Errorf("mqtt broker %v not connected", m.mqttHost)
	}

	// TODO: add actions in this list to generate config.
	binarySensors := []byte{core.ActionMotion, core.ActionDoor}
	sirens := []byte{core.ActionBuzzer}
	sensors := []byte{core.ActionTemp, core.ActionLight, core.ActionVolt}
	buttons := []byte{core.ActionReset}
	_ = buttons

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
			availTopic := fmt.Sprintf(templTopicAvail, brd.Room, brd.Name) // availTopic is shared by all entities.

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
					AvailTopic:  availTopic,
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
					AvailTopic:  availTopic,
					UniqueID:    uniqueID,
					Device:      device,
					ForceUpdate: true,
				}

				// Any action Specific template changes.
				switch actionID {
				case core.ActionTemp:
					// Set value template since Medusa sends both temp and humidity together.
					sensorConfig.ValueTempl = "{{ value_json.temperature }}"
					sensorConfig.UnitMeasure = "°F"

				case core.ActionVolt:
					sensorConfig.UnitMeasure = "V"

				case core.ActionLight:
					sensorConfig.UnitMeasure = "lx"
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
						AvailTopic:  availTopic,
						UniqueID:    fmt.Sprintf(templEntityUniqID, brd.Room, brd.Name, "humidity"),
						Device:      device,
						ValueTempl:  "{{ value_json.humidity }}",
						UnitMeasure: "%",
						ForceUpdate: true,
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
					AvailTopic:   availTopic,
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

// Returns HA device class for each ActionID. Applies only to sensor or binary sensor.
func DeviceClass(actionID byte) string {
	m := map[byte]string{
		0x01: "motion",
		0x02: "temperature",
		0x03: "illuminance",
		0x04: "door",
		0x05: "voltage",
	}

	if v, ok := m[actionID]; ok {
		return v
	}
	return ""
}
