package client

import (
	"fmt"

	"github.com/coreos/fleet/job"
)

// Destroy units for a given target
func (c *FleetClient) Destroy(targets string) (err error) {
	component, num, err := splitTarget(targets)
	if err != nil {
		return
	}
	if num == 0 {
		num, err = c.lastUnit(component)
		if err != nil {
			return err
		}
	}
	name, err := formatUnitName(component, num)
	if err != nil {
		return err
	}

	desiredState := string(job.JobStateInactive)
	err = c.Fleet.SetUnitTargetState(name, desiredState)
	if err != nil {
		return err
	}
	outchan, errchan := waitForUnitStates([]string{name}, desiredState)
	err = printUnitState(name, outchan, errchan)
	if err != nil {
		return err
	}
	if err = c.Fleet.DestroyUnit(name); err != nil {
		return fmt.Errorf("failed destroying job %s: %v", name, err)
	}
	return
}
