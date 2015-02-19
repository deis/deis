package fleet

import (
	"fmt"
	"sync"
	"time"

	"github.com/coreos/fleet/schema"
)

// Stop units and wait for their desiredState
func (c *FleetClient) Stop(targets []string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	// expand @* targets
	expandedTargets, err := expandTargets(c, targets)
	if err != nil {
		errchan <- err
		return
	}

	for _, target := range expandedTargets {
		wg.Add(1)
		go doStop(c, target, wg, outchan, errchan)
	}
	return
}

func doStop(c *FleetClient, target string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	defer wg.Done()

	// prepare string representation
	component, num, err := splitTarget(target)
	if err != nil {
		errchan <- err
		return
	}
	name, err := formatUnitName(component, num)
	if err != nil {
		errchan <- err
		return
	}

	requestState := "loaded"
	desiredState := "dead"

	if err := c.Fleet.SetUnitTargetState(name, requestState); err != nil {
		errchan <- err
		return
	}

	// start with the likely subState to avoid sending it across the channel
	lastSubState := "running"

	for {
		// poll for unit states
		states, err := c.Fleet.UnitStates()
		if err != nil {
			errchan <- err
			return
		}

		// FIXME: fleet UnitStates API forces us to iterate for now
		var currentState *schema.UnitState
		for _, s := range states {
			if name == s.Name {
				currentState = s
				break
			}
		}
		if currentState == nil {
			errchan <- fmt.Errorf("could not find unit: %v", name)
			return
		}

		// if subState changed, send it across the output channel
		if lastSubState != currentState.SystemdSubState {
			outchan <- fmt.Sprintf("\033[0;33m%v:\033[0m %v/%v                                 \r",
				name, currentState.SystemdActiveState, currentState.SystemdSubState)
		}

		// break when desired state is reached
		if currentState.SystemdSubState == desiredState {
			return
		}

		lastSubState = currentState.SystemdSubState
		time.Sleep(250 * time.Millisecond)
	}
}
