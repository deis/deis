package client

import "github.com/coreos/fleet/job"

// Stop sets target units to inactive and blocks until complete
func (c *FleetClient) Stop(target string) (err error) {
	units, err := c.getUnits(target)
	if err != nil {
		return
	}
	desiredState := string(job.JobStateLoaded)
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
