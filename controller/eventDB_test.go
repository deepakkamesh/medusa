package controller_test

import (
	"testing"
	"time"

	"github.com/deepakkamesh/medusa/controller"
)

func TestGetLastEvent(t *testing.T) {
	eventDB, err := controller.NewEventDB()
	if err != nil {
		t.Errorf("failed to create eventDB:%v", err)
	}

	// Log 4 events.
	if err := eventDB.LogEvent(controller.EventLog{time.Now().Add(-200 * time.Millisecond), "ping", 1, "living", "b1", []byte{1, 1, 1}}); err != nil {
		t.Errorf("Failed to log event:%v", err)
	}
	if err := eventDB.LogEvent(controller.EventLog{time.Now().Add(-100 * time.Millisecond), "ping", 1, "living", "b1", []byte{1, 1, 1}}); err != nil {
		t.Errorf("Failed to log event:%v", err)
	}

	if err := eventDB.LogEvent(controller.EventLog{time.Now().Add(-50 * time.Millisecond), "temperature", 50, "living", "b1", []byte{1, 1, 1}}); err != nil {
		t.Errorf("Failed to log event:%v", err)
	}

	if err := eventDB.LogEvent(controller.EventLog{time.Now(), "ping", 1, "living", "b1", []byte{1, 1, 1}}); err != nil {
		t.Errorf("Failed to log event:%v", err)
	}

	events, err := eventDB.GetLastEvent("ping", "living", "b1", 2)
	if err != nil {
		t.Errorf("Unable to get events from eventDB:%v", err)
	}

	// Validate we get only 2 events.
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %v", len(events))
	}

	// Validate we got the 2 ping events.
	if events[0].Metric != "ping" || events[1].Metric != "ping" {
		t.Errorf("Expected ping got %v", events[0].Metric)
	}
}

func TestPurgeDB(t *testing.T) {
	eventDB, err := controller.NewEventDB()
	if err != nil {
		t.Errorf("failed to create eventDB:%v", err)
	}

	// Log 2 events.
	if err := eventDB.LogEvent(controller.EventLog{time.Now().Add(-200 * time.Millisecond), "ping", 1, "living", "b1", []byte{1, 1, 1}}); err != nil {
		t.Errorf("Failed to log event:%v", err)
	}
	if err := eventDB.LogEvent(controller.EventLog{time.Now().Add(-100 * time.Millisecond), "ping", 1, "living", "b1", []byte{1, 1, 1}}); err != nil {
		t.Errorf("Failed to log event:%v", err)
	}

	events, err := eventDB.GetLastEvent("ping", "living", "b1", 2)
	if err != nil {
		t.Errorf("Unable to get events from eventDB:%v", err)
	}

	// Validate we get only 2 events.
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %v", len(events))
	}

	if err := eventDB.PurgeDB(); err != nil {
		t.Errorf("Failed to purgeDB: %v", err)
	}

	// Validate we get no events.
	events, err = eventDB.GetLastEvent("ping", "living", "b1", 2)
	if err != nil {
		t.Errorf("Unable to get events from eventDB:%v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %v", len(events))
	}
}
