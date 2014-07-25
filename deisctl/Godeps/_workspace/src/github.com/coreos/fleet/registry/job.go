package registry

import (
	"errors"
	"path"
	"strings"

	log "github.com/coreos/fleet/third_party/github.com/golang/glog"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/event"
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/unit"
)

const (
	jobPrefix = "job"
)

// Jobs lists all Jobs known by the Registry
func (r *EtcdRegistry) Jobs() ([]job.Job, error) {
	var jobs []job.Job

	req := etcd.Get{
		Key:       path.Join(r.keyPrefix, jobPrefix),
		Sorted:    true,
		Recursive: true,
	}

	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isKeyNotFound(err) {
			err = nil
		}
		return jobs, err
	}

	for _, dir := range resp.Node.Nodes {
		for _, node := range dir.Nodes {
			if !strings.HasSuffix(node.Key, "object") {
				continue
			}

			j := r.getJobFromJSON(node.Value)
			if j == nil {
				continue
			}

			jobs = append(jobs, *j)
		}
	}

	return jobs, nil
}

// GetJobTarget looks up where the given job is scheduled. If the job has
// been scheduled, the ID the target machine is returned. Otherwise, an
// empty string is returned.
func (r *EtcdRegistry) JobTarget(jobName string) (string, error) {
	req := etcd.Get{
		Key:       r.jobTargetAgentPath(jobName),
		Sorted:    false,
		Recursive: true,
	}

	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isKeyNotFound(err) {
			err = nil
		}
		return "", err
	}

	return resp.Node.Value, nil
}

func (r *EtcdRegistry) ClearJobTarget(jobName, machID string) error {
	req := etcd.Delete{
		Key:           r.jobTargetAgentPath(jobName),
		PreviousValue: machID,
	}

	_, err := r.etcd.Do(&req)
	if isKeyNotFound(err) {
		err = nil
	}

	return err
}

// Job looks for a Job of the given name in the Registry. It returns a fully
// hydrated Job on success, or nil on any kind of failure.
func (r *EtcdRegistry) Job(jobName string) (j *job.Job, err error) {
	req := etcd.Get{
		Key: path.Join(r.keyPrefix, jobPrefix, jobName, "object"),
	}

	resp, err := r.etcd.Do(&req)
	if err != nil {
		if isKeyNotFound(err) {
			err = nil
		}
		return
	}

	j = r.getJobFromJSON(resp.Node.Value)
	if j == nil {
		log.V(1).Infof("Error unmarshaling Job(%s): %v", jobName, err)
	}
	return
}

func (r *EtcdRegistry) getJobFromJSON(val string) *job.Job {
	var jm jobModel
	if err := unmarshal(val, &jm); err != nil {
		return nil
	}

	return r.getJobFromModel(jm)
}

func (r *EtcdRegistry) getJobFromModel(jm jobModel) *job.Job {
	var err error
	var unit *unit.Unit

	// New-style Jobs should have a populated UnitHash, and the contents of the Unit are stored separately in the Registry
	if !jm.UnitHash.Empty() {
		unit = r.getUnitByHash(jm.UnitHash)
		if unit == nil {
			log.Warningf("No Unit found in Registry for Job(%s)", jm.Name)
			return nil
		}
		if unit.Hash() != jm.UnitHash {
			log.Errorf("Unit Hash %s does not match expected %s for Job(%s)!", unit.Hash(), jm.UnitHash, jm.Name)
			return nil
		}
		log.V(2).Infof("Got Unit for Job(%s) from registry", jm.Name)
	} else {
		// Old-style Jobs had "Payloads" instead of Units, also stored separately in the Registry
		log.V(2).Infof("Legacy Job(%s) has no PayloadHash - looking for associated Payload", jm.Name)
		unit, err = r.getUnitFromLegacyPayload(jm.Name)
		if err != nil {
			log.Errorf("Error retrieving legacy payload for Job(%s)", jm.Name)
			return nil
		} else if unit == nil {
			log.Warningf("No Payload found in Registry for Job(%s)", jm.Name)
			return nil
		}

		log.Infof("Migrating legacy Payload(%s)", jm.Name)
		if err := r.storeOrGetUnit(*unit); err != nil {
			log.Warningf("Unable to migrate legacy Payload: %v", err)
		}
	}

	j := job.NewJob(jm.Name, *unit)

	j.UnitState = r.getUnitState(jm.Name)
	j.State = r.determineJobState(jm.Name)

	return j
}

// jobModel is used for serializing and deserializing Jobs stored in the Registry
type jobModel struct {
	Name     string
	UnitHash unit.Hash
}

// DestroyJob removes a Job object from the repository, along with any legacy
// associated Payload and SignatureSet. It does not yet remove underlying
// Units from the repository.
func (r *EtcdRegistry) DestroyJob(jobName string) error {
	req := etcd.Delete{
		Key:       path.Join(r.keyPrefix, jobPrefix, jobName),
		Recursive: true,
	}

	r.etcd.Do(&req)

	// TODO(jonboulle): add unit reference counting and actually destroying Units
	r.destroyLegacyPayload(jobName)
	r.destroySignatureSetOfJob(jobName)
	// TODO(jonboulle): handle errors

	return nil
}

// destroyLegacyPayload removes an old-style Payload from the registry
func (r *EtcdRegistry) destroyLegacyPayload(payloadName string) {
	req := etcd.Delete{
		Key: path.Join(r.keyPrefix, payloadPrefix, payloadName),
	}
	r.etcd.Do(&req)
}

// CreateJob attempts to store a Job and its associated Unit in the registry
func (r *EtcdRegistry) CreateJob(j *job.Job) (err error) {
	if err := r.storeOrGetUnit(j.Unit); err != nil {
		return err
	}

	jm := jobModel{
		Name:     j.Name,
		UnitHash: j.UnitHash,
	}
	json, err := marshal(jm)
	if err != nil {
		return
	}

	req := etcd.Create{
		Key:   path.Join(r.keyPrefix, jobPrefix, j.Name, "object"),
		Value: json,
	}

	_, err = r.etcd.Do(&req)
	if err != nil && isNodeExist(err) {
		err = errors.New("job already exists")
	}

	return
}

func (r *EtcdRegistry) JobTargetState(jobName string) (*job.JobState, error) {
	req := etcd.Get{
		Key: r.jobTargetStatePath(jobName),
	}
	resp, err := r.etcd.Do(&req)
	if err != nil {
		log.Errorf("Unable to determine target-state of Job(%s): %v", jobName, err)
		return nil, err
	}

	return job.ParseJobState(resp.Node.Value), nil
}

func (r *EtcdRegistry) SetJobTargetState(jobName string, state job.JobState) error {
	req := etcd.Set{
		Key:   r.jobTargetStatePath(jobName),
		Value: string(state),
	}
	_, err := r.etcd.Do(&req)
	return err
}

func (es *EventStream) filterJobTargetStateChanges(resp *etcd.Result) *event.Event {
	if resp.Action != "set" {
		return nil
	}

	dir, baseName := path.Split(resp.Node.Key)
	if baseName != "target-state" {
		return nil
	}

	dir = strings.TrimSuffix(dir, "/")
	jobName := path.Base(dir)

	ts := job.ParseJobState(resp.Node.Value)
	if ts == nil {
		return nil
	}

	cs := es.registry.determineJobState(jobName)
	if *cs == *ts {
		return nil
	}

	var cType string
	switch *cs {
	case job.JobStateInactive:
		cType = "CommandLoadJob"
	case job.JobStateLoaded:
		if *ts == job.JobStateInactive {
			cType = "CommandUnloadJob"
		} else if *ts == job.JobStateLaunched {
			cType = "CommandStartJob"
		}
	case job.JobStateLaunched:
		if *ts == job.JobStateLoaded {
			cType = "CommandStopJob"
		} else if *ts == job.JobStateInactive {
			cType = "CommandUnloadJob"
		}
	}

	if cType == "" {
		return nil
	}

	agent, _ := es.registry.JobTarget(jobName)
	return &event.Event{cType, jobName, agent}
}

func (r *EtcdRegistry) ScheduleJob(jobName string, machID string) error {
	req := etcd.Create{
		Key:   r.jobTargetAgentPath(jobName),
		Value: machID,
	}
	_, err := r.etcd.Do(&req)
	return err
}

func (r *EtcdRegistry) LockJob(jobName, context string) *TimedResourceMutex {
	return r.lockResource("job", jobName, context)
}

func filterEventJobScheduled(resp *etcd.Result) *event.Event {
	if resp.Action != "create" {
		return nil
	}

	dir, baseName := path.Split(resp.Node.Key)
	if baseName != "target" {
		return nil
	}

	dir = strings.TrimSuffix(dir, "/")
	dir, jobName := path.Split(dir)

	dir = strings.TrimSuffix(dir, "/")
	dir, prefixName := path.Split(dir)

	if prefixName != jobPrefix {
		return nil
	}

	return &event.Event{"EventJobScheduled", jobName, resp.Node.Value}
}

func filterEventJobUnscheduled(resp *etcd.Result) *event.Event {
	if resp.Action != "delete" && resp.Action != "compareAndDelete" {
		return nil
	}

	dir, baseName := path.Split(resp.Node.Key)
	if baseName != "target" {
		return nil
	}

	dir = strings.TrimSuffix(dir, "/")
	dir, jobName := path.Split(dir)

	dir = strings.TrimSuffix(dir, "/")
	dir, prefixName := path.Split(dir)

	if prefixName != jobPrefix {
		return nil
	}

	if resp.PrevNode == nil {
		return nil
	}

	return &event.Event{"EventJobUnscheduled", jobName, resp.PrevNode.Value}
}

func filterEventJobDestroyed(resp *etcd.Result) *event.Event {
	if resp.Action != "delete" {
		return nil
	}

	dir, jobName := path.Split(resp.Node.Key)
	dir = strings.TrimSuffix(dir, "/")
	dir, prefixName := path.Split(dir)

	if prefixName != jobPrefix {
		return nil
	}

	return &event.Event{"EventJobDestroyed", jobName, nil}
}

func (r *EtcdRegistry) jobTargetAgentPath(jobName string) string {
	return path.Join(r.keyPrefix, jobPrefix, jobName, "target")
}

func (r *EtcdRegistry) jobTargetStatePath(jobName string) string {
	return path.Join(r.keyPrefix, jobPrefix, jobName, "target-state")
}
