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

package registry

import (
	"reflect"
	"testing"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/unit"
)

func TestFakeRegistryUnitLifecycle(t *testing.T) {
	reg := NewFakeRegistry()

	units, err := reg.Units()
	if err != nil {
		t.Fatalf("Received error while calling Jobs: %v", err)
	}
	if !reflect.DeepEqual([]job.Unit{}, units) {
		t.Fatalf("Expected no units, got %v", units)
	}

	uf, _ := unit.NewUnitFile("")
	u1 := job.Unit{Name: "u1.service", Unit: *uf, TargetState: job.JobStateLoaded}
	err = reg.CreateUnit(&u1)
	if err != nil {
		t.Fatalf("Received error while calling CreateUnit: %v", err)
	}

	units, err = reg.Units()
	if err != nil {
		t.Fatalf("Received error while calling Units: %v", err)
	}
	if len(units) != 1 {
		t.Fatalf("Expected 1 Unit, got %v", units)
	}
	if !reflect.DeepEqual(u1, units[0]) {
		t.Fatalf("Expected unit %v, got %v", u1, units[0])
	}

	err = reg.ScheduleUnit("u1.service", "XXX")
	if err != nil {
		t.Fatalf("Received error while calling ScheduleUnit: %v", err)
	}

	su, err := reg.ScheduledUnit("u1.service")
	if err != nil {
		t.Fatalf("Received error while calling ScheduledUnit: %v", err)
	}
	if su.TargetMachineID != "XXX" {
		t.Fatalf("Unit should be scheduled to XXX, got %v", su.TargetMachineID)
	}

	err = reg.DestroyUnit("u1.service")
	if err != nil {
		t.Fatalf("Received error while calling DestroyUnit: %v", err)
	}

	units, err = reg.Units()
	if err != nil {
		t.Fatalf("Received error while calling Units: %v", err)
	}
	if !reflect.DeepEqual([]job.Unit{}, units) {
		t.Fatalf("Expected no units, got %v", units)
	}
}
