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
	"sync"

	"github.com/coreos/fleet/pkg"
)

func NewFakeUnitManager() *FakeUnitManager {
	return &FakeUnitManager{u: map[string]bool{}}
}

type FakeUnitManager struct {
	sync.RWMutex
	u map[string]bool
}

func (fum *FakeUnitManager) Load(name string, u UnitFile) error {
	fum.Lock()
	defer fum.Unlock()

	fum.u[name] = false
	return nil
}

func (fum *FakeUnitManager) ReloadUnitFiles() error {
	return nil
}

func (fum *FakeUnitManager) Unload(name string) {
	fum.Lock()
	defer fum.Unlock()

	delete(fum.u, name)
}

func (fum *FakeUnitManager) TriggerStart(string) {}
func (fum *FakeUnitManager) TriggerStop(string)  {}

func (fum *FakeUnitManager) Units() ([]string, error) {
	fum.RLock()
	defer fum.RUnlock()

	lst := make([]string, 0, len(fum.u))
	for name, _ := range fum.u {
		lst = append(lst, name)
	}
	return lst, nil
}

func (fum *FakeUnitManager) GetUnitState(name string) (us *UnitState, err error) {
	fum.RLock()
	defer fum.RUnlock()

	if _, ok := fum.u[name]; ok {
		us = &UnitState{
			LoadState:   "loaded",
			ActiveState: "active",
			SubState:    "running",
		}
	}
	return
}

func (fum *FakeUnitManager) GetUnitStates(filter pkg.Set) (map[string]*UnitState, error) {
	fum.RLock()
	defer fum.RUnlock()

	states := make(map[string]*UnitState)
	for _, name := range filter.Values() {
		if _, ok := fum.u[name]; ok {
			states[name] = &UnitState{"loaded", "active", "running", "", "", name}
		}
	}

	return states, nil
}

func (fum *FakeUnitManager) MarshalJSON() ([]byte, error) {
	return nil, nil
}
