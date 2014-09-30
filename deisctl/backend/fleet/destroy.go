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
		_, err := c.Units(target)
		if err != nil {
			return err
		}
		component, num, err := splitTarget(target)
		if err != nil {
			return err
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

func (c *FleetClient) destroyServiceUnit(component string, num int) error {
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
	if err := printUnitState(name, outchan, errchan); err != nil {
		return err
	}
	if err = c.Fleet.DestroyUnit(name); err != nil {
		return fmt.Errorf("failed destroying job %s: %v", name, err)
	}
	return nil
}

func (c *FleetClient) destroyDataUnit(component string) error {
	name, err := formatUnitName(component, 0)
	desiredState := string(job.JobStateInactive)
	if err != nil {
		return err
	}
	if err := c.Fleet.SetUnitTargetState(name, desiredState); err != nil {
		return err
	}
	outchan, errchan := waitForUnitStates([]string{name}, desiredState)
	if err := printUnitState(name, outchan, errchan); err != nil {
		return err
	}
	if err := c.Fleet.DestroyUnit(name); err != nil {
		return fmt.Errorf("failed destroying job %s: %v", name, err)
	}
	return nil
}
