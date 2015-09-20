package fleet

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/coreos/fleet/schema"
	"github.com/deis/deis/pkg/prettyprint"
)

// Start units and wait for their desiredState
func (c *FleetClient) Start(targets []string, wg *sync.WaitGroup, out, ew io.Writer) {
	// expand @* targets
	expandedTargets, err := c.expandTargets(targets)
	if err != nil {
		io.WriteString(ew, err.Error())
		return
	}

	for _, target := range expandedTargets {
		wg.Add(1)
		go doStart(c, target, wg, out, ew)
	}
	return
}

func doStart(c *FleetClient, target string, wg *sync.WaitGroup, out, ew io.Writer) {
	defer wg.Done()

	// prepare string representation
	component, num, err := splitTarget(target)
	if err != nil {
		io.WriteString(ew, err.Error())
		return
	}
	name, err := formatUnitName(component, num)
	if err != nil {
		io.WriteString(ew, err.Error())
		return
	}

	requestState := "launched"
	desiredState := "running"

	if err := c.Fleet.SetUnitTargetState(name, requestState); err != nil {
		io.WriteString(ew, err.Error())
		return
	}

	// start with the likely subState to avoid sending it across the channel
	lastSubState := "dead"

	for {
		// poll for unit states
		states, err := c.Fleet.UnitStates()
		if err != nil {
			io.WriteString(ew, err.Error())
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
			fmt.Fprintf(ew, "Could not find unit: %v\n", name)
			return
		}

		// if subState changed, send it across the output channel
		if lastSubState != currentState.SystemdSubState {
			fmt.Fprintf(out, stateFmt+"\n", name, currentState.SystemdActiveState, currentState.SystemdSubState)
		}

		// break when desired state is reached
		if currentState.SystemdSubState == desiredState {
			return
		}

		lastSubState = currentState.SystemdSubState

		if lastSubState == "failed" {
			o := prettyprint.Colorize("{{.Red}}The service '%s' failed while starting.{{.Default}}\n")
			fmt.Fprintf(ew, o, target)
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
}
