package pkg

import (
	"testing"
	"time"
)

type fakeEventStream struct {
	ret chan Event
}

func (f *fakeEventStream) Next(chan struct{}) chan Event {
	return f.ret
}

func (f *fakeEventStream) trigger() {
	go func() {
		f.ret <- Event("asdf")
	}()
}

// TestPeriodicReconcilerRun attempts to validate the behaviour of the central Run
// loop of the PeriodicReconciler, particularly its response to events
func TestPeriodicReconcilerRun(t *testing.T) {
	fes := &fakeEventStream{make(chan Event)}
	called := make(chan struct{})
	rec := func() {
		go func() {
			called <- struct{}{}
		}()
	}
	pr := &reconciler{
		ival:    time.Hour, // Implausibly high reconcile interval we never expect to reach
		rFunc:   rec,
		eStream: fes,
	}
	// launch the PeriodicReconciler in the background
	prDone := make(chan bool)
	stop := make(chan bool)
	go func() {
		pr.Run(stop)
		prDone <- true
	}()
	// no reconcile yet expected
	select {
	case <-called:
		t.Fatalf("rFunc() called unexpectedly!")
	default:
	}
	// now, send an event on the EventStream and ensure rFunc occurs
	fes.trigger()
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatalf("rFunc() not called after trigger!")
	}
	// assert rFunc was only called once
	select {
	case <-called:
		t.Fatalf("rFunc() called unexpectedly!")
	default:
	}
	// another event should work OK
	fes.trigger()
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatalf("rFunc() not called after trigger!")
	}
	// again, assert rFunc was only called once
	select {
	case <-called:
		t.Fatalf("rFunc() called unexpectedly!")
	default:
	}
	// stop the loop
	close(stop)
	// now, sending an event should do nothing
	fes.trigger()
	select {
	case <-called:
		t.Fatalf("rFunc() called unexpectedly!")
	default:
	}
	// and the PeriodicReconciler should have shut down
	select {
	case <-prDone:
	case <-time.After(time.Second):
		t.Fatalf("PeriodicReconciler.Run did not return after stop signal!")
	}
}
