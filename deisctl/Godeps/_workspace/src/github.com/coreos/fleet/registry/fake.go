package registry

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/coreos/fleet/Godeps/_workspace/src/github.com/coreos/go-semver/semver"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/unit"
)

func NewFakeRegistry() *FakeRegistry {
	return &FakeRegistry{
		machines:  []machine.MachineState{},
		jobStates: map[string]map[string]*unit.UnitState{},
		jobs:      map[string]job.Job{},
		units:     []unit.UnitFile{},
		version:   nil,
	}
}

type FakeRegistry struct {
	// Not all methods required by the Registry interface are implemented
	// by the FakeRegistry. Any calls to these unimplemented methods will
	// result in a panic.
	Registry
	sync.RWMutex

	machines  []machine.MachineState
	jobStates map[string]map[string]*unit.UnitState
	jobs      map[string]job.Job
	units     []unit.UnitFile
	version   *semver.Version
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

	f.jobStates = make(map[string]map[string]*unit.UnitState, len(states))
	for _, us := range states {
		us := us
		name := us.UnitName
		if _, ok := f.jobStates[name]; !ok {
			f.jobStates[name] = make(map[string]*unit.UnitState)
		}
		f.jobStates[name][us.MachineID] = &us
	}
}

func (f *FakeRegistry) SetUnits(units []unit.UnitFile) {
	f.Lock()
	defer f.Unlock()

	f.units = units
}

func (f *FakeRegistry) SetLatestVersion(v semver.Version) {
	f.Lock()
	defer f.Unlock()

	f.version = &v
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

	f.jobStates[jobName][unitState.MachineID] = unitState
}

func (f *FakeRegistry) RemoveUnitState(jobName string) error {
	delete(f.jobStates, jobName)
	return nil
}

func (f *FakeRegistry) UnitStates() ([]*unit.UnitState, error) {
	f.Lock()
	defer f.Unlock()

	var states []*unit.UnitState

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

func (f *FakeRegistry) LatestVersion() (*semver.Version, error) {
	f.RLock()
	defer f.RUnlock()

	return f.version, nil
}

func (f *FakeRegistry) UnitHeartbeat(name, machID string, ttl time.Duration) error {
	return nil
}

func (f *FakeRegistry) ClearUnitHeartbeat(string) {}
