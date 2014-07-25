package registry

import (
	"testing"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/event"
)

func TestPipe(t *testing.T) {
	etcdchan := make(chan *etcd.Result)

	translate1 := func(resp *etcd.Result) *event.Event {
		return &event.Event{"TranslateTest1", resp.Action, nil}
	}

	translate2 := func(resp *etcd.Result) *event.Event {
		return &event.Event{"TranslateTest2", resp.Action, "foo"}
	}

	filters := []func(resp *etcd.Result) *event.Event{translate1, translate2}

	eventchan := make(chan *event.Event)
	stopchan := make(chan bool)

	send := func(ev *event.Event) {
		eventchan <- ev
	}

	go pipe(etcdchan, filters, send, stopchan)

	resp := etcd.Result{Action: "TestAction", Node: &etcd.Node{Key: "/", ModifiedIndex: 0}}
	etcdchan <- &resp

	ev1 := <-eventchan
	ev2 := <-eventchan

	close(stopchan)

	if ev1.Type != "TranslateTest1" {
		t.Fatalf("Expected ev1.Type \"TranslateTest1\" but got %q", ev1.Type)
	}

	if ev1.Payload.(string) != "TestAction" {
		t.Fatalf("Expected ev1.Payload \"TestAction\", but got something else")
	}

	if ev1.Context != nil {
		t.Fatalf("Expected ev1.Context be nil")
	}

	if ev2.Type != "TranslateTest2" {
		t.Fatalf("Expected ev2.Type \"TranslateTest2\" but got %q", ev2.Type)
	}

	payload := ev2.Payload.(string)
	if payload != "TestAction" {
		t.Fatalf("Expected ev2.Payload \"TestAction\", but got %q", payload)
	}

	if ev2.Context == nil {
		t.Fatalf("Expected ev2.Context to be non-nil")
	}

	ctx := ev2.Context.(string)
	if ctx != "foo" {
		t.Fatalf("Expected ev2.Context value \"foo\", got %q", ctx)
	}
}
