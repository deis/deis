package client

import (
	"fmt"
	"os"

	"github.com/coreos/fleet/job"
)

// Start launches target units and blocks until active
func (c *FleetClient) Start(target string, data bool) (err error) {
	units, err := c.getUnits(target)
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
