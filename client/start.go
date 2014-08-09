package client

import (
	"fmt"
	"os"
	"regexp"

	"github.com/coreos/fleet/job"
)

// Start launches target units and blocks until active
func (c *FleetClient) Start(target string, data bool) (err error) {

	// see if we were provided a specific target
	r := regexp.MustCompile(`([a-z-]+)\.([\d]+)`)
	match := r.FindStringSubmatch(target)
	var component string
	if len(match) == 3 {
		component = match[1]
	} else {
		component = target
	}
	units, err := c.getUnits(component)
	if err != nil {
		return
	}
	newState := job.JobStateLaunched
	for _, unitName := range units {
		err = c.Fleet.SetJobTargetState(unitName, newState)
		if err != nil {
			return err
		}
	}
	errchan := waitForJobStates(units, testUnitStateActive, 0, os.Stdout)
	for err := range errchan {
		return fmt.Errorf("error waiting for active: %v", err)
	}
	return nil
}
