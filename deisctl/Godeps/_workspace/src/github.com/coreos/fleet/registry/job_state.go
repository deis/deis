package registry

import (
	"path"
	"time"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/job"
)

func (r *EtcdRegistry) determineJobState(jobName string) *job.JobState {
	state := job.JobStateInactive

	tgt, _ := r.JobTarget(jobName)
	if tgt == "" {
		return &state
	}

	if r.getUnitState(jobName) == nil {
		return &state
	}

	state = job.JobStateLoaded

	agent, pulse := r.CheckJobPulse(jobName)
	if !pulse || agent != tgt {
		return &state
	}

	state = job.JobStateLaunched
	return &state
}

func (r *EtcdRegistry) JobHeartbeat(jobName, agentMachID string, ttl time.Duration) error {
	req := etcd.Set{
		Key:   r.jobHeartbeatPath(jobName),
		Value: agentMachID,
		TTL:   ttl,
	}
	_, err := r.etcd.Do(&req)
	return err
}

func (r *EtcdRegistry) CheckJobPulse(jobName string) (string, bool) {
	req := etcd.Get{
		Key: r.jobHeartbeatPath(jobName),
	}
	resp, err := r.etcd.Do(&req)
	if err != nil {
		return "", false
	}

	return resp.Node.Value, true
}

func (r *EtcdRegistry) ClearJobHeartbeat(jobName string) {
	req := etcd.Delete{
		Key: r.jobHeartbeatPath(jobName),
	}
	r.etcd.Do(&req)
}

func (r *EtcdRegistry) jobHeartbeatPath(jobName string) string {
	return path.Join(r.keyPrefix, jobPrefix, jobName, "job-state")
}
