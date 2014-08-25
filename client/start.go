package client

import "github.com/coreos/fleet/job"

// Start launches target units and blocks until active
func (c *FleetClient) Start(target string) (err error) {
	units, err := c.getUnits(target)
	if err != nil {
		return
	}
	desiredState := string(job.JobStateLaunched)
	for _, name := range units {
		err = c.Fleet.SetUnitTargetState(name, desiredState)
		if err != nil {
			return err
		}
		outchan, errchan := waitForUnitStates(units, desiredState)
		err = printUnitState(name, outchan, errchan)
		if err != nil {
			return err
		}
	}
	return nil
}
