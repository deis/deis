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

package client

import (
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/registry"
	"github.com/coreos/fleet/schema"
)

type RegistryClient struct {
	registry.Registry
}

func (rc *RegistryClient) Units() ([]*schema.Unit, error) {
	rUnits, err := rc.Registry.Units()
	if err != nil {
		return nil, err
	}

	sUnits, err := rc.Registry.Schedule()
	if err != nil {
		return nil, err
	}

	sUnitMap := make(map[string]*job.ScheduledUnit, len(sUnits))
	for _, sUnit := range sUnits {
		sUnit := sUnit
		sUnitMap[sUnit.Name] = &sUnit
	}

	units := make([]*schema.Unit, len(rUnits))
	for i, ru := range rUnits {
		units[i] = schema.MapUnitToSchemaUnit(&ru, sUnitMap[ru.Name])
	}

	return units, nil
}

func (rc *RegistryClient) Unit(name string) (*schema.Unit, error) {
	rUnit, err := rc.Registry.Unit(name)
	if err != nil || rUnit == nil {
		return nil, err
	}

	var sUnit *job.ScheduledUnit
	// Only non-global units have an associated Schedule
	if !rUnit.IsGlobal() {
		sUnit, err = rc.Registry.ScheduledUnit(name)
		if err != nil {
			return nil, err
		}
	}

	return schema.MapUnitToSchemaUnit(rUnit, sUnit), nil
}

func (rc *RegistryClient) CreateUnit(u *schema.Unit) error {
	rUnit := job.Unit{
		Name:        u.Name,
		Unit:        *schema.MapSchemaUnitOptionsToUnitFile(u.Options),
		TargetState: job.JobStateInactive,
	}

	if len(u.DesiredState) > 0 {
		ts, err := job.ParseJobState(u.DesiredState)
		if err != nil {
			return err
		}

		rUnit.TargetState = ts
	}

	return rc.Registry.CreateUnit(&rUnit)
}

func (rc *RegistryClient) UnitStates() ([]*schema.UnitState, error) {
	rUnitStates, err := rc.Registry.UnitStates()
	if err != nil {
		return nil, err
	}

	states := make([]*schema.UnitState, len(rUnitStates))
	for i, rus := range rUnitStates {
		states[i] = schema.MapUnitStateToSchemaUnitState(rus)
	}

	return states, nil
}

func (rc *RegistryClient) SetUnitTargetState(name, target string) error {
	return rc.Registry.SetUnitTargetState(name, job.JobState(target))
}
