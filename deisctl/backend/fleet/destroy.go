package fleet

import (
	"fmt"
	"strings"

	"github.com/coreos/fleet/job"
)

// Destroy units for a given target
func (c *FleetClient) Destroy(targets []string) error {
	for _, target := range targets {
		// check if the unit exists
		units, err := c.Units(target)
		if err != nil {
			return err
		}
		component, num, err := splitTarget(target)
		if err != nil {
			return err
		}
		// if no number is specified, destroy ALL THE UNITS!
		if num == 0 {
			num = len(units)
		}
		if strings.HasSuffix(component, "-data") {
			err = c.destroyDataUnit(component)
		} else {
			err = c.destroyServiceUnit(component, num)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *FleetClient) destroyServiceUnit(component string, num int) (err error) {
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
	return err
}

func (c *FleetClient) destroyDataUnit(component string) (err error) {
	name, err := formatUnitName(component, 0)
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
	return err

}
