package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/coreos/fleet/job"
)

type testJob func(j *job.Job) bool

func testJobStateLoaded(j *job.Job) bool {
	if j == nil || j.State == nil {
		return false
	}
	return *(j.State) == job.JobStateLoaded
}

func testJobStateLaunched(j *job.Job) bool {
	if j == nil || j.State == nil {
		return false
	}
	return *(j.State) == job.JobStateLaunched
}

func testJobStateInactive(j *job.Job) bool {
	if j == nil || j.State == nil {
		return false
	}
	return *(j.State) == job.JobStateInactive
}

func testUnitStateActive(j *job.Job) bool {
	if j == nil || j.UnitState == nil {
		return false
	}
	return j.UnitState.ActiveState == "active"

}

// stateCheck defines how to monitor a job state
type stateCheck struct {
	test      testJob
	statechan chan *jobState
	errchan   chan error
}

// newStateCheck returns a StateCheck struct with new channels for monitoring
func newStateCheck(test testJob) *stateCheck {
	statechan := make(chan *jobState)
	errchan := make(chan error)
	return &stateCheck{test, statechan, errchan}
}

// waitForJobStates polls each of the indicated jobs until each of their
// states is equal to that which the caller indicates via stateCheck test
func waitForJobStates(jobs []string, check *stateCheck) error {
	var wg sync.WaitGroup

	// check each job with the stateCheck
	for _, name := range jobs {
		wg.Add(1)
		go checkJobState(name, check, &wg)
	}

	// wait for all jobs to complete
	go func() {
		wg.Wait()
		close(check.errchan)
	}()

	// print output while jobs are transitioning
	defer fmt.Printf("\n")
	for {
		select {
		// read from state channel
		case state := <-check.statechan:
			// return on closed channel
			if state == nil {
				return nil
			}
			// otherwise print output
			if state.loaded == "inactive" {
				fmt.Printf("\033[0;33m%v:\033[0m %v                                 \r",
					state.name, state.loaded)
			} else {
				fmt.Printf("\033[0;33m%v:\033[0m %v, %v (%v)                        \r",
					state.name, state.loaded, state.active, state.sub)
			}
		// read from error channel
		case err := <-check.errchan:
			return err
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func checkJobState(jobName string, check *stateCheck, wg *sync.WaitGroup) {
	defer wg.Done()
	sleep := 100 * time.Millisecond
	for {
		if assertJobState(jobName, check) {
			return
		}
		time.Sleep(sleep)
	}
}

type jobState struct {
	name   string
	loaded string
	active string
	sub    string
}

func newJobState(name string, j *job.Job) *jobState {
	var (
		loaded string
		active string
		sub    string
	)
	if j.State != nil {
		loaded = fmt.Sprintf("%v", *(j.State))
	}
	if j.UnitState != nil {
		active = j.UnitState.ActiveState
		sub = j.UnitState.SubState
	}
	return &jobState{name, loaded, active, sub}
}

func assertJobState(name string, check *stateCheck) bool {
	j, err := cAPI.Job(name)
	if err != nil {
		check.errchan <- fmt.Errorf("Error retrieving Job(%s) from Registry: %v", name, err)
		return false
	}

	// send current state to the output channel
	check.statechan <- newJobState(name, j)

	// if test function
	if check.test(j) {
		close(check.statechan)
		return true
	}
	return false
}
