package pkg

import (
	"time"

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
// ival, or in response to anything emitted from the channel returned by
// trigFunc
func NewPeriodicReconciler(interval time.Duration, recFunc func(), eStream EventStream) PeriodicReconciler {
	return &reconciler{
		ival:    interval,
		rFunc:   recFunc,
		eStream: eStream,
	}
}

type reconciler struct {
	ival    time.Duration
	rFunc   func()
	eStream EventStream
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

	ticker := time.After(r.ival)
	for {
		select {
		case <-stop:
			log.V(1).Info("Reconciler exiting due to stop signal")
			return
		case <-ticker:
			ticker = time.After(r.ival)
			log.V(1).Info("Reconciler tick")
			r.rFunc()
		case <-trigger:
			ticker = time.After(r.ival)
			log.V(1).Info("Reconciler triggered")
			r.rFunc()
		}
	}

}
