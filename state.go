package deisctl

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/coreos/fleet/job"
)

type testJob func(j *job.Job) bool

func testJobStateLoaded(j *job.Job) bool {
	return j == nil || j.State == nil || *(j.State) != job.JobStateLoaded
}

func testJobStateLaunched(j *job.Job) bool {
	return j == nil || j.State == nil || *(j.State) != job.JobStateLaunched
}

func testJobStateInactive(j *job.Job) bool {
	return j == nil || j.State == nil || *(j.State) != job.JobStateInactive
}

func testUnitStateActive(j *job.Job) bool {
	return j == nil || j.UnitState == nil || j.UnitState.ActiveState != "active"
}

// TODO: refactor to separate presentation to io.Writer from status polling

// waitForJobStates polls each of the indicated jobs until each of their
// states is equal to that which the caller indicates, or until the
// polling operation times out. waitForJobStates will retry forever, or
// up to maxAttempts times before timing out if maxAttempts is greater
// than zero. Returned is an error channel used to communicate when
// timeouts occur. The returned error channel will be closed after all
// polling operation is complete.
func waitForJobStates(jobs []string, test testJob, maxAttempts int, out io.Writer) chan error {
	errchan := make(chan error)
	var wg sync.WaitGroup
	for _, name := range jobs {
		wg.Add(1)
		go checkJobState(name, test, maxAttempts, out, &wg, errchan)
	}
	go func() {
		wg.Wait()
		close(errchan)
	}()
	return errchan
}

func checkJobState(jobName string, test testJob, maxAttempts int, out io.Writer, wg *sync.WaitGroup, errchan chan error) {
	defer wg.Done()
	sleep := 100 * time.Millisecond
	if maxAttempts < 1 {
		for {
			if assertJobState(jobName, test, out) {
				return
			}
			time.Sleep(sleep)
		}
	} else {
		for attempt := 0; attempt < maxAttempts; attempt++ {
			if assertJobState(jobName, test, out) {
				return
			}
			time.Sleep(sleep)
		}
		errchan <- fmt.Errorf("timed out waiting for job %s to report state", jobName)
	}
}

func assertJobState(name string, test testJob, out io.Writer) (ret bool) {
	j, err := cAPI.Job(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving Job(%s) from Registry: %v", name, err)
		return
	}
	if test(j) {
		if j.State != nil && j.UnitState != nil && j.UnitState.ActiveState != "" && j.UnitState.SubState != "" {
			fmt.Fprintf(out, "\033[0;33m%v:\033[0m %v, %v (%v)           \r", j.Name, *(j.State), j.UnitState.ActiveState, j.UnitState.SubState)
		}
		return
	}
	ret = true

	var msg string
	if j.State != nil && j.UnitState != nil && j.UnitState.ActiveState != "" && j.UnitState.SubState != "" {
		msg = fmt.Sprintf("\033[0;33m%v:\033[0m %v, %v (%v)", name, *(j.State), j.UnitState.ActiveState, j.UnitState.SubState)
	} else {
		msg = fmt.Sprintf("\033[0;33m%v:\033[0m %v", name, *(j.State))
	}

	machines, err := cAPI.Machines()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed retrieving list of Machines from Registry: %v", err)
	}
	for _, ms := range machines {
		if ms.ID != j.TargetMachineID {
			continue
		}
		msg = fmt.Sprintf("%s on %s", msg, machineFullLegend(ms, false))
		break
	}

	fmt.Fprintln(out, msg)
	return
}
