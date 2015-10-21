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
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/coreos/fleet/log"
)

type Event string

// EventStream generates a channel which will emit an event as soon as one of
// interest occurs. Any background operation associated with the channel
// should be terminated when stop is closed.
type EventStream interface {
	Next(stop chan struct{}) chan Event
}

type PeriodicReconciler interface {
	Run(stop chan bool)
}

// NewPeriodicReconciler creates a PeriodicReconciler that will run recFunc at least every
// ival, or in response to anything emitted from EventStream.Next()
func NewPeriodicReconciler(interval time.Duration, recFunc func(), eStream EventStream) PeriodicReconciler {
	return &reconciler{
		ival:    interval,
		rFunc:   recFunc,
		eStream: eStream,
		clock:   clockwork.NewRealClock(),
	}
}

type reconciler struct {
	ival    time.Duration
	rFunc   func()
	eStream EventStream
	clock   clockwork.Clock
}

func (r *reconciler) Run(stop chan bool) {
	trigger := make(chan struct{})
	go func() {
		abort := make(chan struct{})
		for {
			select {
			case <-stop:
				close(abort)
				return
			case <-r.eStream.Next(abort):
				trigger <- struct{}{}
			}
		}
	}()

	ticker := r.clock.After(r.ival)

	// When starting up, reconcile once immediately
	log.Debug("Initial reconciliation commencing")
	r.rFunc()

	for {
		select {
		case <-stop:
			log.Debug("Reconciler exiting due to stop signal")
			return
		case <-ticker:
			ticker = r.clock.After(r.ival)
			log.Debug("Reconciler tick")
			r.rFunc()
		case <-trigger:
			ticker = r.clock.After(r.ival)
			log.Debug("Reconciler triggered")
			r.rFunc()
		}
	}

}
