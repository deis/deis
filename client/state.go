package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/coreos/fleet/job"
)

type testJob func(j *job.Job) bool

func jobStateLoaded(j *job.Job) bool {
	if j == nil || j.State == nil {
		return false
	}
	return *(j.State) == job.JobStateLoaded
}

func jobStateLaunched(j *job.Job) bool {
	if j == nil || j.State == nil {
		return false
	}
	return *(j.State) == job.JobStateLaunched
}

func jobStateInactive(j *job.Job) bool {
	if j == nil || j.State == nil {
		return false
	}
	return *(j.State) == job.JobStateInactive
}

func unitStateActive(j *job.Job) bool {
	if j == nil || j.UnitState == nil {
		return false
	}
	return j.UnitState.ActiveState == "active"

}

// stateCheck defines how to monitor a job state
type stateCheck struct {
	test  testJob
	state chan *jobState
}

// newStateCheck returns a StateCheck struct with new channels for monitoring
func newStateCheck(test testJob) *stateCheck {
	state := make(chan *jobState)
	return &stateCheck{test, state}
}

// waitForJobStates polls each of the indicated jobs until each of their
// states is equal to that which the caller indicates via stateCheck test
func waitForJobStates(jobs []string, test testJob) error {
	var wg sync.WaitGroup
	errchan := make(chan error)
	check := newStateCheck(test)

	// check each job with the stateCheck
	for _, name := range jobs {
		wg.Add(1)
		go checkJobState(name, check, &wg, errchan)
	}

	// wait for all jobs to complete
	go func() {
		wg.Wait()
		close(errchan)
	}()

	// print output while jobs are transitioning
	defer fmt.Printf("\n")
	for {
		select {
		// read from state channel
		case state := <-check.state:
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
		case err := <-errchan:
			return err
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func checkJobState(jobName string, check *stateCheck, wg *sync.WaitGroup, errchan chan error) {
	defer wg.Done()
	sleep := 100 * time.Millisecond
	for {
		if assertJobState(jobName, check, errchan) {
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

func assertJobState(name string, check *stateCheck, errchan chan error) bool {
	j, err := cAPI.Job(name)
	if err != nil {
		errchan <- fmt.Errorf("Error retrieving Job(%s) from Registry: %v", name, err)
		return false
	}

	// send current state to the output channel
	check.state <- newJobState(name, j)

	// if state matches, close the channel and return
	if check.test(j) {
		close(check.state)
		return true
	}
	return false
}
