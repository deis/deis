package fleet

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/coreos/fleet/schema"
)

// waitForUnitSubStates polls each of the indicated jobs until each of their
// systemd substates is equal to that which the caller indicates
func waitForUnitSubStates(units []string, desiredState string) (outchan chan *schema.Unit, errchan chan error) {
	var wg sync.WaitGroup
	errchan = make(chan error)
	outchan = make(chan *schema.Unit)

	// check each unit for desired state
	for _, name := range units {
		wg.Add(1)
		go checkUnitSubState(name, desiredState, outchan, errchan, &wg)
	}

	// wait for all jobs to complete
	go func() {
		wg.Wait()
		close(outchan)
		close(errchan)
	}()

	return outchan, errchan
}

func checkUnitSubState(name string, desiredState string, outchan chan *schema.Unit, errchan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if assertUnitSubState(name, desiredState, outchan, errchan) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func assertUnitSubState(name string, desiredState string, outchan chan *schema.Unit, errchan chan error) bool {
	u, err := cAPI.Unit(name)
	if err != nil {
		errchan <- fmt.Errorf("Error retrieving Job(%s) from Registry: %v", name, err)
		return false
	}

	// send unit across the output channel
	outchan <- u

	unitState, err := unitState(name)
	if err != nil {
		errchan <- fmt.Errorf("Error retrieving Unit state from Registry: %v", err)
		return false
	}

	if unitState.SystemdSubState == desiredState {
		return true
	}
	return false
}

// unitState retrieves a UnitState based on its name.
// FIXME: add this to fleet's API
func unitState(name string) (*schema.UnitState, error) {
	unitStates, err := cAPI.UnitStates()
	if err != nil {
		return nil, err
	}

	for _, state := range unitStates {
		if state.Name == name {
			return state, nil
		}
	}
	return nil, errors.New("Could not find unit state: " + name)
}

// waitForUnitStates polls each of the indicated jobs until each of their
// states is equal to that which the caller indicates
func waitForUnitStates(units []string, desiredState string) (outchan chan *schema.Unit, errchan chan error) {

	var wg sync.WaitGroup
	errchan = make(chan error)
	outchan = make(chan *schema.Unit)

	// check each unit for desired state
	for _, name := range units {
		wg.Add(1)
		go checkUnitState(name, desiredState, outchan, errchan, &wg)
	}

	// wait for all jobs to complete
	go func() {
		wg.Wait()
		close(outchan)
		close(errchan)
	}()

	return outchan, errchan

}

func checkUnitState(name string, desiredState string, outchan chan *schema.Unit, errchan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if assertUnitState(name, desiredState, outchan, errchan) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func assertUnitState(name string, desiredState string, outchan chan *schema.Unit, errchan chan error) bool {
	u, err := cAPI.Unit(name)
	if err != nil {
		errchan <- fmt.Errorf("Error retrieving Job(%s) from Registry: %v", name, err)
		return false
	}

	// send unit across the output channel
	outchan <- u

	if u.DesiredState == u.CurrentState {
		return true
	}
	return false
}

func printUnitSubState(name string, outchan chan *schema.Unit, errchan chan error) error {
	// print output while jobs are transitioning
	defer fmt.Printf("\n")
	for {
		select {
		case u := <-outchan:
			// return on closed channel
			if u == nil {
				return nil
			}
			// ignore units that don't match our unit
			if u.Name != name {
				continue
			}
			// retrieve the unit's state
			unitState, err := unitState(u.Name)
			if err != nil {
				return err
			}
			fmt.Printf("\033[0;33m%v:\033[0m %v                                 \r",
				u.Name, unitState.SystemdSubState)
		// read from error channel
		case err := <-errchan:
			// continue processing if error channel closed
			if err == nil {
				continue
			} else if strings.Contains(err.Error(), "timeout reached") {
				// ignore intermittent timeout errors
				continue
			}
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func printUnitState(name string, outchan chan *schema.Unit, errchan chan error) error {
	// print output while jobs are transitioning
	defer fmt.Printf("\n")
	for {
		select {
		case u := <-outchan:
			// return on closed channel
			if u == nil {
				return nil
			}
			// ignore units that don't match our unit
			if u.Name != name {
				continue
			}
			// otherwise print output
			if u.CurrentState != u.DesiredState {
				fmt.Printf("\033[0;33m%v:\033[0m %v (pending)                       \r",
					u.Name, u.CurrentState)
			} else {
				fmt.Printf("\033[0;33m%v:\033[0m %v                                 \r",
					u.Name, u.CurrentState)
			}
		// read from error channel
		case err := <-errchan:
			// continue processing if error channel closed
			if err == nil {
				continue
			} else if strings.Contains(err.Error(), "timeout reached") {
				// ignore intermittent timeout errors
				continue
			}
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
