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

package unit

import (
	"encoding/json"
	"time"

	"github.com/coreos/fleet/log"
	"github.com/coreos/fleet/pkg"
)

type UnitStateHeartbeat struct {
	Name  string
	State *UnitState
}

func NewUnitStateGenerator(mgr UnitManager) *UnitStateGenerator {
	return &UnitStateGenerator{
		mgr:        mgr,
		subscribed: pkg.NewThreadsafeSet(),
	}
}

type UnitStateGenerator struct {
	mgr UnitManager

	subscribed     pkg.Set
	lastSubscribed pkg.Set
}

func (g *UnitStateGenerator) MarshalJSON() ([]byte, error) {
	data := struct {
		Subscribed []string
	}{
		Subscribed: g.subscribed.Values(),
	}

	return json.Marshal(data)
}

// Run periodically calls Generate and sends received *UnitStateHeartbeat
// objects to the provided channel.
func (g *UnitStateGenerator) Run(receiver chan<- *UnitStateHeartbeat, stop chan bool) {
	tick := time.Tick(time.Second)
	for {
		select {
		case <-stop:
			return
		case <-tick:
			beatchan, err := g.Generate()
			if err != nil {
				log.Errorf("Failed fetching current unit states: %v", err)
				continue
			}

			for ush := range beatchan {
				receiver <- ush
			}
		}
	}
}

// Generate returns and fills a channel with *UnitStateHeartbeat objects. Objects will
// only be returned for units to which this generator is currently subscribed.
func (g *UnitStateGenerator) Generate() (<-chan *UnitStateHeartbeat, error) {
	var lastSubscribed pkg.Set
	if g.lastSubscribed != nil {
		lastSubscribed = g.lastSubscribed.Copy()
	}

	subscribed := g.subscribed.Copy()
	g.lastSubscribed = subscribed

	reportable, err := g.mgr.GetUnitStates(subscribed)
	if err != nil {
		return nil, err
	}

	beatchan := make(chan *UnitStateHeartbeat)
	go func() {
		for name, us := range reportable {
			us := us
			beatchan <- &UnitStateHeartbeat{
				Name:  name,
				State: us,
			}
		}

		if lastSubscribed != nil {
			// For all units that were part of the subscription list
			// last time Generate ran, but are now not part of that
			// list, send nil-State heartbeats to signal removal
			for _, name := range lastSubscribed.Sub(subscribed).Values() {
				beatchan <- &UnitStateHeartbeat{
					Name: name,
				}
			}
		}

		close(beatchan)
	}()

	return beatchan, nil
}

// Subscribe adds a unit to the internal state filter
func (g *UnitStateGenerator) Subscribe(name string) {
	g.subscribed.Add(name)
}

// Unsubscribe removes a unit from the internal state filter
func (g *UnitStateGenerator) Unsubscribe(name string) {
	g.subscribed.Remove(name)
}
