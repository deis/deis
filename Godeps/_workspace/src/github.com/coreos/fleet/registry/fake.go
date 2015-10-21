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
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/pkg/lease"
	"github.com/coreos/fleet/unit"
	"github.com/coreos/go-semver/semver"
)

func NewFakeRegistry() *FakeRegistry {
	return &FakeRegistry{
		machines:      []machine.MachineState{},
		jobStates:     map[string]map[string]*unit.UnitState{},
		jobs:          map[string]job.Job{},
		daemonVersion: nil,
	}
}

type FakeRegistry struct {
	// Not all methods required by the Registry interface are implemented
	// by the FakeRegistry. Any calls to these unimplemented methods will
	// result in a panic.
	Registry
	sync.RWMutex

	machines      []machine.MachineState
	jobStates     map[string]map[string]*unit.UnitState
	jobs          map[string]job.Job
	daemonVersion *semver.Version
}

func (f *FakeRegistry) SetMachines(machines []machine.MachineState) {
	f.Lock()
	defer f.Unlock()

	f.machines = machines
}

func (f *FakeRegistry) SetJobs(jobs []job.Job) {
	f.Lock()
	defer f.Unlock()

	f.jobs = make(map[string]job.Job, len(jobs))
	for _, j := range jobs {
		f.jobs[j.Name] = j
	}
}

func (f *FakeRegistry) SetUnitStates(states []unit.UnitState) {
	f.Lock()
	defer f.Unlock()

	f.jobStates = make(map[string]map[string]*unit.UnitState)
	for _, us := range states {
		us := us
		name := us.UnitName
		if _, ok := f.jobStates[name]; !ok {
			f.jobStates[name] = make(map[string]*unit.UnitState)
		}
		f.jobStates[name][us.MachineID] = &us
	}
}

func (f *FakeRegistry) Machines() ([]machine.MachineState, error) {
	f.RLock()
	defer f.RUnlock()

	return f.machines, nil
}

func (f *FakeRegistry) Units() ([]job.Unit, error) {
	f.RLock()
	defer f.RUnlock()

	var sorted sort.StringSlice
	for _, j := range f.jobs {
		sorted = append(sorted, j.Name)
	}
	sorted.Sort()

	units := make([]job.Unit, len(f.jobs))
	for i, jName := range sorted {
		j := f.jobs[jName]
		u := job.Unit{
			Name:        j.Name,
			Unit:        j.Unit,
			TargetState: j.TargetState,
		}
		units[i] = u
	}

	return units, nil
}

func (f *FakeRegistry) Schedule() ([]job.ScheduledUnit, error) {
	f.RLock()
	defer f.RUnlock()

	var sorted sort.StringSlice
	for _, j := range f.jobs {
		sorted = append(sorted, j.Name)
	}
	sorted.Sort()

	sUnits := make([]job.ScheduledUnit, 0, len(f.jobs))
	for _, jName := range sorted {
		j := f.jobs[jName]
		su := job.ScheduledUnit{
			Name:            j.Name,
			State:           j.State,
			TargetMachineID: j.TargetMachineID,
		}
		sUnits = append(sUnits, su)
	}

	return sUnits, nil
}

func (f *FakeRegistry) Unit(name string) (*job.Unit, error) {
	f.RLock()
	defer f.RUnlock()

	j, ok := f.jobs[name]
	if !ok {
		return nil, nil
	}

	u := job.Unit{
		Name:        j.Name,
		Unit:        j.Unit,
		TargetState: j.TargetState,
	}
	return &u, nil
}

func (f *FakeRegistry) ScheduledUnit(name string) (*job.ScheduledUnit, error) {
	f.RLock()
	defer f.RUnlock()

	j, ok := f.jobs[name]
	if !ok {
		return nil, nil
	}

	su := job.ScheduledUnit{
		Name:            j.Name,
		State:           j.State,
		TargetMachineID: j.TargetMachineID,
	}
	return &su, nil
}

func (f *FakeRegistry) CreateUnit(u *job.Unit) error {
	f.Lock()
	defer f.Unlock()

	_, ok := f.jobs[u.Name]
	if ok {
		return errors.New("unit already exists")
	}

	j := job.Job{
		Name: u.Name,
		Unit: u.Unit,
	}

	f.jobs[u.Name] = j
	return f.unsafeSetUnitTargetState(u.Name, u.TargetState)
}

func (f *FakeRegistry) DestroyUnit(name string) error {
	f.Lock()
	defer f.Unlock()

	delete(f.jobs, name)
	return nil
}

func (f *FakeRegistry) SetUnitTargetState(name string, target job.JobState) error {
	f.Lock()
	defer f.Unlock()

	return f.unsafeSetUnitTargetState(name, target)
}

func (f *FakeRegistry) unsafeSetUnitTargetState(name string, target job.JobState) error {
	j, ok := f.jobs[name]

	if !ok {
		return errors.New("job does not exist")
	}

	j.TargetState = target
	f.jobs[name] = j

	return nil
}

func (f *FakeRegistry) ScheduleUnit(name string, machID string) error {
	f.Lock()
	defer f.Unlock()

	j, ok := f.jobs[name]

	if !ok {
		return errors.New("unit does not exist")
	}

	j.TargetMachineID = machID
	f.jobs[name] = j

	return nil
}

func (f *FakeRegistry) SaveUnitState(jobName string, unitState *unit.UnitState, ttl time.Duration) {
	f.Lock()
	defer f.Unlock()

	if _, ok := f.jobStates[jobName]; !ok {
		f.jobStates[jobName] = make(map[string]*unit.UnitState)
	}
	f.jobStates[jobName][unitState.MachineID] = unitState
}

func (f *FakeRegistry) RemoveUnitState(jobName string) error {
	delete(f.jobStates, jobName)
	return nil
}

func (f *FakeRegistry) UnitStates() ([]*unit.UnitState, error) {
	f.Lock()
	defer f.Unlock()

	states := make([]*unit.UnitState, 0)

	// Sort by unit name, then by machineID
	sUnitNames := make([]string, 0)
	for name := range f.jobStates {
		sUnitNames = append(sUnitNames, name)
	}
	sort.Strings(sUnitNames)

	for _, name := range sUnitNames {
		sMIDs := make([]string, 0)
		for machineID := range f.jobStates[name] {
			sMIDs = append(sMIDs, machineID)
		}
		sort.Strings(sMIDs)
		for _, mID := range sMIDs {
			states = append(states, f.jobStates[name][mID])
		}
	}
	return states, nil
}

func (f *FakeRegistry) UnitHeartbeat(name, machID string, ttl time.Duration) error {
	return nil
}

func (f *FakeRegistry) ClearUnitHeartbeat(string) {}

func NewFakeClusterRegistry(dVersion *semver.Version, eVersion int) *FakeClusterRegistry {
	return &FakeClusterRegistry{
		dVersion: dVersion,
		eVersion: eVersion,
	}
}

type FakeClusterRegistry struct {
	dVersion *semver.Version
	eVersion int
}

func (fc *FakeClusterRegistry) LatestDaemonVersion() (*semver.Version, error) {
	return fc.dVersion, nil
}

func (fc *FakeClusterRegistry) EngineVersion() (int, error) {
	return fc.eVersion, nil
}

func (fc *FakeClusterRegistry) UpdateEngineVersion(from, to int) error {
	if fc.eVersion != from {
		return errors.New("version mismatch")
	}

	fc.eVersion = to
	return nil
}

func (fl *FakeLeaseRegistry) SetLease(name, machID string, ver int, ttl time.Duration) *fakeLease {
	l := &fakeLease{
		name:   name,
		machID: machID,
		ver:    ver,
		ttl:    ttl,
		reg:    fl,
	}

	fl.leaseMap[name] = l
	return l
}

type fakeLease struct {
	name   string
	machID string
	ver    int
	ttl    time.Duration
	reg    *FakeLeaseRegistry
}

func (l *fakeLease) MachineID() string {
	return l.machID
}

func (l *fakeLease) Version() int {
	return l.ver
}

func (l *fakeLease) TimeRemaining() time.Duration {
	return l.ttl
}

func (l *fakeLease) Index() uint64 {
	return 0
}

func (l *fakeLease) Renew(ttl time.Duration) error {
	if l.reg == nil {
		return errors.New("already released")
	}

	l.ttl = ttl
	return nil
}

func (l *fakeLease) Release() error {
	if l.reg == nil {
		return errors.New("already released")
	}

	delete(l.reg.leaseMap, l.name)
	l.reg = nil
	return nil
}

func NewFakeLeaseRegistry() *FakeLeaseRegistry {
	return &FakeLeaseRegistry{
		leaseMap: make(map[string]lease.Lease),
	}
}

type FakeLeaseRegistry struct {
	leaseMap map[string]lease.Lease
}

func (fl *FakeLeaseRegistry) GetLease(name string) (lease.Lease, error) {
	return fl.leaseMap[name], nil
}

func (fl *FakeLeaseRegistry) AcquireLease(name, machID string, ver int, ttl time.Duration) (lease.Lease, error) {
	if _, ok := fl.leaseMap[name]; ok {
		return nil, errors.New("already exists")
	}

	l := &fakeLease{
		name:   name,
		machID: machID,
		ver:    ver,
		ttl:    ttl,
		reg:    fl,
	}

	fl.leaseMap[name] = l
	return l, nil
}

func (fl *FakeLeaseRegistry) StealLease(name, machID string, ver int, ttl time.Duration, idx uint64) (lease.Lease, error) {
	if idx != 0 {
		panic("unable to test StealLease with index other than zero")
	}

	_, ok := fl.leaseMap[name]
	if !ok {
		return nil, errors.New("does not exist")
	}

	l := &fakeLease{
		name:   name,
		machID: machID,
		ver:    ver,
		ttl:    ttl,
		reg:    fl,
	}

	fl.leaseMap[name] = l
	return l, nil
}
