package registry

import (
	"time"

	"github.com/coreos/fleet/third_party/github.com/coreos/go-semver/semver"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/sign"
	"github.com/coreos/fleet/unit"
)

type Registry interface {
	Bids(jb *job.JobOffer) ([]job.JobBid, error)
	CheckJobPulse(jobName string) (string, bool)
	ClearJobHeartbeat(jobName string)
	ClearJobTarget(jobName, machID string) error
	CreateJob(j *job.Job) (err error)
	CreateJobOffer(jo *job.JobOffer) error
	CreateSignatureSet(ss *sign.SignatureSet) error
	DestroyJob(jobName string) error
	DestroySignatureSet(tag string)
	Job(jobName string) (j *job.Job, err error)
	JobHeartbeat(jobName, agentMachID string, ttl time.Duration) error
	Jobs() ([]job.Job, error)
	JobSignatureSet(name string) (*sign.SignatureSet, error)
	JobTarget(jobName string) (string, error)
	JobTargetState(jobName string) (*job.JobState, error)
	LatestVersion() (*semver.Version, error)
	LockJob(jobName, context string) *TimedResourceMutex
	LockJobOffer(jobName, context string) *TimedResourceMutex
	LockMachine(machID, context string) *TimedResourceMutex
	Machines() ([]machine.MachineState, error)
	RemoveMachineState(machID string) error
	RemoveUnitState(jobName string) error
	ResolveJobOffer(jobName string) error
	SaveUnitState(jobName string, unitState *unit.UnitState)
	ScheduleJob(jobName string, machID string) error
	SetJobTargetState(jobName string, state job.JobState) error
	SetMachineState(ms machine.MachineState, ttl time.Duration) (uint64, error)
	SubmitJobBid(jb *job.JobBid)
	UnresolvedJobOffers() []job.JobOffer
}
