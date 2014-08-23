package client

import "github.com/coreos/fleet/job"

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
	check := newStateCheck(testJobStateInactive)
	err = waitForJobStates(units, check)
	if err != nil {
		return err
	}
	return nil
}
