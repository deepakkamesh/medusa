package controller_test

import (
	"testing"

	"github.com/deepakkamesh/medusa/controller"
	"github.com/deepakkamesh/medusa/controller/core"
	"github.com/deepakkamesh/medusa/controller/mocks"
	"github.com/golang/mock/gomock"
)

func TestHARecv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msg := mocks.NewMockMessage(ctrl)

	hachan := make(chan controller.HAMsg, 2) // So not to block PubHandler.

	ha := controller.HomeAssistant{
		HAMsg: hachan,
	}

	msg.EXPECT().Topic().Return("giant/living/b1/buzzer/set")
	msg.EXPECT().Payload().Return([]byte("{\"state\":\"true\"}"))

	ha.MQTTPubHandler(nil, msg)
	m := <-hachan

	if m.Action != "buzzer" {
		t.Errorf("got %v want buzzer", m.Action)
	}
	if m.BoardName != "b1" {
		t.Errorf("got %v want b1", m.BoardName)
	}
	if m.Room != "living" {
		t.Errorf("got %v want living", m.Room)
	}
	if m.State != true {
		t.Errorf("got %t want true", m.State)
	}
}

func TestHASend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockClient(ctrl)
	tk := mocks.NewMockToken(ctrl)
	ha := controller.HomeAssistant{
		MQTTClient: m,
		CoreCfg: &core.Config{
			Relays: []*core.Relay{},
			Boards: []*core.Board{
				{Room: "living", Name: "b1", Addr: []byte{1, 1, 1}, Actions: []byte{0x01, 0x02}},
			},
		},
	}

	// Test SendMotion.
	m.EXPECT().Publish("giant/living/b1/motion/state", gomock.Any(), false, "ON").Return(tk)
	tk.EXPECT().Wait()
	ha.SendMotion("living", "b1", true)

	// Test Temp Humidity. Rounds off to 1 precision point.
	m.EXPECT().Publish("giant/living/b1/temp/state", gomock.Any(), false, "{\"temperature\":70,\"humidity\":45.2}").Return(tk)
	tk.EXPECT().Wait()
	ha.SendTemp("living", "b1", 70.1, 45.3)

	// Test Volt. 2 decimal precision.
	m.EXPECT().Publish("giant/living/b1/volt/state", gomock.Any(), false, "3.23").Return(tk)
	tk.EXPECT().Wait()
	ha.SendVolt("living", "b1", 3.23234)

	// Test ping events.
	m.EXPECT().Publish("giant/living/b1/avail", gomock.Any(), false, "offline").Return(tk)
	tk.EXPECT().Wait()
	ha.SendAvail("living", "b1", "offline")
}

func TestMQTTDiscoveryConfig(t *testing.T) {
	cfg, e := core.NewConfig("./core/core.cfg.test.json")
	if e != nil {
		t.Errorf("failed to load config %v", e)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockClient(ctrl)
	tk := mocks.NewMockToken(ctrl)

	ha := controller.HomeAssistant{
		MQTTClient: m,
		CoreCfg:    cfg,
	}

	// Test if Entity MQTT config generation works.
	m.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tk)
	tk.EXPECT().Wait().AnyTimes()
	if err := ha.SendMQTTDiscoveryConfig(false); err != nil {
		t.Errorf("Test failed %v", err)
	}

	// Test if clean entity works.
	m.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(tk)
	tk.EXPECT().Wait().AnyTimes()
	if err := ha.SendMQTTDiscoveryConfig(true); err != nil {
		t.Errorf("Test failed %v", err)
	}

}
