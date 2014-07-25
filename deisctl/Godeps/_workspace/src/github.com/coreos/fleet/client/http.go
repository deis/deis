package client

import (
	"encoding/base64"
	"net/http"

	"github.com/coreos/fleet/third_party/code.google.com/p/google-api-go-client/googleapi"
	"github.com/coreos/fleet/third_party/github.com/coreos/go-semver/semver"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

func NewHTTPClient(c *http.Client) (API, error) {
	svc, err := schema.New(c)
	if err != nil {
		return nil, err
	}
	return &HTTPClient{svc: svc}, nil
}

type HTTPClient struct {
	svc *schema.Service

	//NOTE(bcwaldon): This is only necessary until the API interface
	// is fully implemented by HTTPClient
	API
}

func (c *HTTPClient) Machines() ([]machine.MachineState, error) {
	machines := make([]machine.MachineState, 0)
	call := c.svc.Machines.List()
	for call != nil {
		page, err := call.Do()
		if err != nil {
			return nil, err
		}

		machines = append(machines, mapMachinePageToMachineStates(page.Machines)...)

		if len(page.NextPageToken) > 0 {
			call = c.svc.Machines.List()
			call.NextPageToken(page.NextPageToken)
		} else {
			call = nil
		}
	}
	return machines, nil
}

func mapMachinePageToMachineStates(entities []*schema.Machine) []machine.MachineState {
	machines := make([]machine.MachineState, len(entities))
	for i, _ := range entities {
		me := entities[i]

		ms := machine.MachineState{
			ID:       me.Id,
			PublicIP: me.PrimaryIP,
		}

		ms.Metadata = make(map[string]string, len(me.Metadata))
		for k, v := range me.Metadata {
			ms.Metadata[k] = v
		}

		machines[i] = ms
	}

	return machines
}

func (c *HTTPClient) Jobs() ([]job.Job, error) {
	machines, err := c.Machines()
	if err != nil {
		return nil, err
	}

	mm := make(map[string]*machine.MachineState, len(machines))
	for i, _ := range machines {
		m := machines[i]
		mm[m.ID] = &m
	}

	var jobs []job.Job
	call := c.svc.Units.List()
	for call != nil {
		page, err := call.Do()
		if err != nil {
			return nil, err
		}

		units, err := mapUnitPageToJobs(page.Units, mm)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, units...)

		if len(page.NextPageToken) > 0 {
			call = c.svc.Units.List()
			call.NextPageToken(page.NextPageToken)
		} else {
			call = nil
		}
	}
	return jobs, nil
}

func (c *HTTPClient) Job(name string) (*job.Job, error) {
	u, err := c.svc.Units.Get(name).Do()
	if err != nil {
		if is404(err) {
			err = nil
		}
		return nil, err
	}

	if u == nil {
		return nil, nil
	}

	machines, err := c.Machines()
	if err != nil {
		return nil, err
	}

	mm := make(map[string]*machine.MachineState, len(machines))
	for i, _ := range machines {
		m := machines[i]
		mm[m.ID] = &m
	}

	return mapUnitToJob(u, mm)
}

func (c *HTTPClient) JobTarget(name string) (string, error) {
	u, err := c.svc.Units.Get(name).Do()
	if err != nil {
		return "", err
	}

	return u.TargetMachineID, nil
}

func mapUnitPageToJobs(entities []*schema.Unit, mm map[string]*machine.MachineState) ([]job.Job, error) {
	jobs := make([]job.Job, len(entities))
	for i, _ := range entities {
		entity := entities[i]
		j, err := mapUnitToJob(entity, mm)
		if err != nil {
			return nil, err
		}
		if j != nil {
			jobs[i] = *j
		}
	}

	return jobs, nil
}

func mapUnitToJob(entity *schema.Unit, mm map[string]*machine.MachineState) (*job.Job, error) {
	contents, err := base64.StdEncoding.DecodeString(entity.FileContents)
	if err != nil {
		return nil, err
	}
	u, err := unit.NewUnit(string(contents))
	if err != nil {
		return nil, err
	}

	js := job.JobState(entity.CurrentState)
	j := job.Job{
		Name:     entity.Name,
		State:    &js,
		Unit:     *u,
		UnitHash: u.Hash(),
	}

	// populate a UnitState object only if the entity
	// is actually reporting relevant data
	if entity.Systemd != nil {
		j.UnitState = &unit.UnitState{
			LoadState:   entity.Systemd.LoadState,
			ActiveState: entity.Systemd.ActiveState,
			SubState:    entity.Systemd.SubState,
		}
		if len(entity.Systemd.MachineID) > 0 {
			m, ok := mm[entity.Systemd.MachineID]
			if ok {
				j.UnitState.MachineState = m
			} else {
				j.UnitState.MachineState = &machine.MachineState{ID: entity.Systemd.MachineID}
			}
		}
	}

	return &j, nil
}

func (c *HTTPClient) DestroyJob(name string) error {
	req := schema.DeletableUnit{
		Name: name,
	}
	return c.svc.Units.Delete(name, &req).Do()
}

func (c *HTTPClient) CreateJob(j *job.Job) error {
	req := schema.DesiredUnitState{
		Name:         j.Name,
		DesiredState: string(job.JobStateInactive),
		FileContents: base64.StdEncoding.EncodeToString([]byte(j.Unit.Raw)),
	}
	return c.svc.Units.Set(j.Name, &req).Do()
}

func (c *HTTPClient) SetJobTargetState(name string, state job.JobState) error {
	req := schema.DesiredUnitState{
		Name:         name,
		DesiredState: string(state),
	}
	return c.svc.Units.Set(name, &req).Do()
}

//NOTE(bcwaldon): This is only temporary until a better version negotiation mechanism is in place
func (c *HTTPClient) LatestVersion() (*semver.Version, error) {
	return semver.NewVersion("0.0.0")
}

func encodeUnitContents(u *unit.Unit) string {
	return base64.StdEncoding.EncodeToString([]byte(u.Raw))
}

func is404(err error) bool {
	googerr, ok := err.(*googleapi.Error)
	return ok && googerr.Code == http.StatusNotFound
}
