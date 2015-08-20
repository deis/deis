package fleet

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/coreos/fleet/schema"
	"github.com/deis/deis/pkg/prettyprint"
)

var stateFmt = prettyprint.Colorize("{{.Yellow}}%v:{{.Default}} %v/%v")

// Stop units and wait for their desiredState
func (c *FleetClient) Stop(targets []string, wg *sync.WaitGroup, out, ew io.Writer) {
	// expand @* targets
	expandedTargets, err := c.expandTargets(targets)
	if err != nil {
		fmt.Fprintln(ew, err.Error())
		return
	}

	for _, target := range expandedTargets {
		wg.Add(1)
		go doStop(c, target, wg, out, ew)
	}
	return
}

func doStop(c *FleetClient, target string, wg *sync.WaitGroup, out, ew io.Writer) {
	defer wg.Done()

	// prepare string representation
	component, num, err := splitTarget(target)
	if err != nil {
		fmt.Fprintln(ew, err.Error())
		return
	}
	name, err := formatUnitName(component, num)
	if err != nil {
		fmt.Fprintln(ew, err.Error())
		return
	}

	requestState := "loaded"
	desiredState := "dead"

	if err := c.Fleet.SetUnitTargetState(name, requestState); err != nil {
		fmt.Fprintln(ew, err.Error())
		return
	}

	// start with the likely subState to avoid sending it across the channel
	lastSubState := "running"

	for {
		// poll for unit states
		states, err := c.Fleet.UnitStates()
		if err != nil {
			fmt.Fprintln(ew, err.Error())
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
			l := prettyprint.Overwritef(stateFmt, name, currentState.SystemdActiveState, currentState.SystemdSubState)
			fmt.Fprintf(out, l)
		}

		// break when desired state is reached
		if currentState.SystemdSubState == desiredState {
			fmt.Fprintln(out)
			return
		}

		lastSubState = currentState.SystemdSubState

		if lastSubState == "failed" {
			o := prettyprint.Colorize("{{.Red}}The service '%s' failed while stopping.{{.Default}}\n")
			fmt.Fprintf(ew, o, target)
			return
		}
		fmt.Println(lastSubState)

		time.Sleep(250 * time.Millisecond)
	}
}
