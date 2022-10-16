package controller_test

import (
	"testing"

	"github.com/deepakkamesh/medusa/controller"
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

	msg.EXPECT().Topic().Return("giant/living/buzzer/set")
	msg.EXPECT().Payload().Return([]byte("{\"state\":\"true\"}"))

	ha.MQTTPubHandler(nil, msg)
	m := <-hachan

	if m.Action != "buzzer" {
		t.Errorf("got %v want buzzer", m.Action)
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
	}

	// Test SendMotion.
	m.EXPECT().Publish("giant/living/motion/state", gomock.Any(), false, "ON").Return(tk)
	tk.EXPECT().Wait()
	ha.SendMotion("living", true)
}
