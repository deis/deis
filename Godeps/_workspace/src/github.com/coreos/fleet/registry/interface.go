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
	"time"

	"github.com/coreos/go-semver/semver"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/unit"
)

type Registry interface {
	ClearUnitHeartbeat(name string)
	CreateUnit(*job.Unit) error
	DestroyUnit(string) error
	UnitHeartbeat(name, machID string, ttl time.Duration) error
	Machines() ([]machine.MachineState, error)
	RemoveMachineState(machID string) error
	RemoveUnitState(jobName string) error
	SaveUnitState(jobName string, unitState *unit.UnitState, ttl time.Duration)
	ScheduleUnit(name, machID string) error
	SetUnitTargetState(name string, state job.JobState) error
	SetMachineState(ms machine.MachineState, ttl time.Duration) (uint64, error)
	UnscheduleUnit(name, machID string) error

	UnitRegistry
}

type UnitRegistry interface {
	Schedule() ([]job.ScheduledUnit, error)
	ScheduledUnit(name string) (*job.ScheduledUnit, error)
	Unit(name string) (*job.Unit, error)
	Units() ([]job.Unit, error)
	UnitStates() ([]*unit.UnitState, error)
}

type ClusterRegistry interface {
	LatestDaemonVersion() (*semver.Version, error)

	// EngineVersion returns the current version of the cluster. If the
	// cluster does not yet have a version, zero will be returned. If
	// any error occurs, an error object will be returned. In this case,
	// the returned version number should be ignored.
	EngineVersion() (int, error)

	// UpdateEngineVersion attempts to compare-and-swap the cluster version
	// from one value to another. Any failures in this process will be
	// indicated by the returned error object. A nil value will be returned
	// on success.
	UpdateEngineVersion(from, to int) error
}
