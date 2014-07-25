package event

import (
	"testing"
)

type TestListener struct {
	evchan chan Event
}

func (l *TestListener) HandleEventTypeOne(ev Event) {
	go func() { l.evchan <- ev }()
}

func TestEventBus(t *testing.T) {
	evchan := make(chan Event)

	bus := NewEventBus()
	bus.AddListener("test", &TestListener{evchan})

	ev := Event{"EventTypeOne", "payload", "Y"}
	bus.Dispatch(&ev)

	select {
	case recv := <-evchan:
		if recv.Payload.(string) != "payload" {
			t.Error("event payload is incorrect")
		}
		if recv.Context.(string) != "Y" {
			t.Error("event context is incorrect")
		}
	default:
		t.Fatalf("Failed to dispatch event")
	}
}

func TestEventBusNoDispatch(t *testing.T) {
	evchan := make(chan Event)

	bus := NewEventBus()
	bus.AddListener("test", &TestListener{evchan})

	go func() {
		ev := Event{"EventTypeTwo", "payload", "Y"}
		bus.Dispatch(&ev)
	}()

	go func() {
		ev := Event{"EventTypeOne", "payload", "Y"}
		bus.Dispatch(&ev)
	}()

	recv := <-evchan
	if recv.Type != "EventTypeOne" {
		t.Fatalf("handler received unexpected event")
	}
}
