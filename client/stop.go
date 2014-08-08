package client

import (
	"fmt"
	"os"

	"github.com/coreos/fleet/job"
)

// Stop sets target units to inactive and blocks until complete
func (c *FleetClient) Stop(target string) (err error) {
	units, err := c.getUnits(target)
	if err != nil {
		return
	}
	newState := job.JobStateInactive
	for _, unitName := range units {
		err = c.Fleet.SetJobTargetState(unitName, newState)
		if err != nil {
			return err
		}
	}
	errchan := waitForJobStates(units, testJobStateInactive, 0, os.Stdout)
	for err := range errchan {
		return fmt.Errorf("error waiting for inactive: %v", err)
	}
	return nil
}
