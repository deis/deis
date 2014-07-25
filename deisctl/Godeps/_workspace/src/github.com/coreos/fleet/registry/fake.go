package registry

import (
	"errors"
	"sync"

	"github.com/coreos/fleet/third_party/github.com/coreos/go-semver/semver"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/unit"
)

func NewFakeRegistry() *FakeRegistry {
	return &FakeRegistry{
		machines:     []machine.MachineState{},
		jobStates:    map[string]*unit.UnitState{},
		jobs:         map[string]job.Job{},
		units:        []unit.Unit{},
		version:      nil,
		bids:         map[string][]job.JobBid{},
		targetStates: map[string]job.JobState{},
		schedule:     map[string]string{},
	}
}

type FakeRegistry struct {
	// Not all methods of required by the Registry interface are implemented
	// by the TestRegistry. Any calls to these unimplemented methods will
	// result in a panic.
	Registry
	sync.RWMutex

	machines     []machine.MachineState
	jobStates    map[string]*unit.UnitState
	jobs         map[string]job.Job
	units        []unit.Unit
	version      *semver.Version
	bids         map[string][]job.JobBid
	targetStates map[string]job.JobState
	schedule     map[string]string
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

func (f *FakeRegistry) SetUnitStates(jobStates map[string]*unit.UnitState) {
	f.Lock()
	defer f.Unlock()

	f.jobStates = jobStates
}

func (f *FakeRegistry) SetUnits(units []unit.Unit) {
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

func (f *FakeRegistry) Jobs() ([]job.Job, error) {
	f.RLock()
	defer f.RUnlock()

	jobs := make([]job.Job, 0, len(f.jobs))
	for _, j := range f.jobs {
		jobs = append(jobs, j)
	}
	return jobs, nil
}

func (f *FakeRegistry) Job(name string) (*job.Job, error) {
	f.RLock()
	defer f.RUnlock()

	j, ok := f.jobs[name]
	if !ok {
		return nil, nil
	}

	j.UnitState = f.jobStates[name]
	return &j, nil
}

func (f *FakeRegistry) CreateJob(j *job.Job) error {
	f.Lock()
	defer f.Unlock()

	_, ok := f.jobs[j.Name]
	if ok {
		return errors.New("Job already exists")
	}

	f.jobs[j.Name] = *j
	return nil
}

func (f *FakeRegistry) DestroyJob(name string) error {
	f.Lock()
	defer f.Unlock()

	delete(f.jobs, name)
	return nil
}

func (f *FakeRegistry) Bids(jo *job.JobOffer) ([]job.JobBid, error) {
	f.RLock()
	defer f.RUnlock()

	return f.bids[jo.Job.Name], nil
}

func (f *FakeRegistry) SubmitJobBid(jb *job.JobBid) {
	f.Lock()
	defer f.Unlock()

	_, ok := f.bids[jb.JobName]
	if !ok {
		f.bids[jb.JobName] = []job.JobBid{}
	}
	f.bids[jb.JobName] = append(f.bids[jb.JobName], *jb)
}

func (f *FakeRegistry) SetJobTargetState(name string, target job.JobState) error {
	f.Lock()
	defer f.Unlock()

	f.targetStates[name] = target
	return nil
}

func (f *FakeRegistry) JobTargetState(jobName string) (*job.JobState, error) {
	f.RLock()
	defer f.RUnlock()

	ts, ok := f.targetStates[jobName]
	if !ok {
		return nil, nil
	}
	return &ts, nil
}

func (f *FakeRegistry) ScheduleJob(name string, machID string) error {
	f.Lock()
	defer f.Unlock()

	f.schedule[name] = machID
	return nil
}

func (f *FakeRegistry) JobTarget(name string) (string, error) {
	f.RLock()
	defer f.RUnlock()

	return f.schedule[name], nil
}

func (f *FakeRegistry) SaveUnitState(jobName string, unitState *unit.UnitState) {
	f.Lock()
	defer f.Unlock()

	f.jobStates[jobName] = unitState
}

func (f *FakeRegistry) LatestVersion() (*semver.Version, error) {
	f.RLock()
	defer f.RUnlock()

	return f.version, nil
}
