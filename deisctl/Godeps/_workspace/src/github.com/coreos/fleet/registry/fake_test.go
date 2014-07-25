package registry

import (
	"reflect"
	"testing"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/unit"
)

func TestFakeRegistryJobLifecycle(t *testing.T) {
	reg := NewFakeRegistry()

	jobs, err := reg.Jobs()
	if err != nil {
		t.Fatalf("Received error while calling Jobs: %v", err)
	}
	if !reflect.DeepEqual([]job.Job{}, jobs) {
		t.Fatalf("Expected no jobs, got %v", jobs)
	}

	u, _ := unit.NewUnit("")
	j1 := job.NewJob("job1.service", *u)
	err = reg.CreateJob(j1)
	if err != nil {
		t.Fatalf("Received error while calling CreateJob: %v", err)
	}

	jobs, err = reg.Jobs()
	if err != nil {
		t.Fatalf("Received error while calling Jobs: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("Expected 1 Job, got %v", jobs)
	}
	if jobs[0].Name != "job1.service" {
		t.Fatalf("Expected Job with name \"job1.service\", got %q", jobs[0].Name)
	}

	err = reg.ScheduleJob("job1.service", "XXX")
	if err != nil {
		t.Fatalf("Received error while calling ScheduleJob: %v", err)
	}

	tgt, err := reg.JobTarget("job1.service")
	if err != nil {
		t.Fatalf("Received error while calling JobTarget: %v", err)
	}
	if tgt != "XXX" {
		t.Fatalf("Job should be scheduled to XXX, got %v", tgt)
	}

	err = reg.DestroyJob("job1.service")
	if err != nil {
		t.Fatalf("Received error while calling DestroyJob: %v", err)
	}

	jobs, err = reg.Jobs()
	if err != nil {
		t.Fatalf("Received error while calling Jobs: %v", err)
	}
	if !reflect.DeepEqual([]job.Job{}, jobs) {
		t.Fatalf("Expected no jobs, got %v", jobs)
	}
}
