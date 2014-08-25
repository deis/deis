package client

import "github.com/coreos/fleet/job"

// Start launches target units and blocks until active
func (c *FleetClient) Start(target string) (err error) {
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
	err = waitForJobStates(units, unitStateActive)
	if err != nil {
		return err
	}
	return nil
}
