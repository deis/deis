// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
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
// loop of the PeriodicReconciler
func TestPeriodicReconcilerRun(t *testing.T) {
	ival := 5 * time.Hour
	fclock := clockwork.NewFakeClock()
	fes := &fakeEventStream{make(chan Event)}
	called := make(chan struct{})
	rec := func() {
		go func() {
			called <- struct{}{}
		}()
	}
	pr := &reconciler{
		ival:    ival,
		rFunc:   rec,
		eStream: fes,
		clock:   fclock,
	}
	// launch the PeriodicReconciler in the background
	prDone := make(chan bool)
	stop := make(chan bool)
	go func() {
		pr.Run(stop)
		prDone <- true
	}()
	// reconcile should have occurred once at start-up
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatalf("rFunc() not called at start-up as expected!")
	}
	// no further reconciles yet expected
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
	// now check that time changes have the expected effect
	fclock.Advance(2 * time.Hour)
	select {
	case <-called:
		t.Fatalf("rFunc() called unexpectedly!")
	default:
	}
	fclock.Advance(3 * time.Hour)
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatalf("rFunc() not called after time event!")
	}

	// stop the PeriodicReconciler
	close(stop)

	// now, sending an event should do nothing
	fes.trigger()
	select {
	case <-called:
		t.Fatalf("rFunc() called unexpectedly!")
	default:
	}
	// and nor should changes in time
	fclock.Advance(10 * time.Hour)
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
