package client

import (
	"fmt"
	"strings"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/unit"
)

// Create schedules a new unit for the given component
// and blocks until the unit is loaded
func (c *FleetClient) Create(target string) (err error) {
	var (
		unitName string
		unitPtr  *unit.Unit
	)
	// create unit
	unitName, unitPtr, err = c.createUnit(target)
	if err != nil {
		return err
	}
	// schedule job
	j := job.NewJob(unitName, *unitPtr)
	if err := c.Fleet.CreateJob(j); err != nil {
		return fmt.Errorf("failed creating job %s: %v", unitName, err)
	}
	newState := job.JobStateLoaded
	err = c.Fleet.SetJobTargetState(unitName, newState)
	if err != nil {
		return err
	}
	check := newStateCheck(testJobStateLoaded)
	err = waitForJobStates([]string{unitName}, check)
	if err != nil {
		return err
	}
	return nil
}

func (c *FleetClient) createUnit(target string) (unitName string, unitPtr *unit.Unit, err error) {
	component, num, err := splitTarget(target)
	if err != nil {
		return
	}
	if strings.HasSuffix(component, "-data") {
		unitName, unitPtr, err = c.createDataUnit(component)
	} else {
		unitName, unitPtr, err = c.createServiceUnit(component, num)
	}
	if err != nil {
		return unitName, unitPtr, err
	}
	return
}

// Create normal service unit
func (c *FleetClient) createServiceUnit(component string, num int) (unitName string, unitPtr *unit.Unit, err error) {
	// if number wasn't provided get next unit number
	if num == 0 {
		num, err = c.nextUnit(component)
		if err != nil {
			return "", nil, err
		}
	}
	// build a fleet unit
	unitName, err = formatUnitName(component, num)
	if err != nil {
		return "", nil, err
	}
	unitPtr, err = NewUnit(component)
	if err != nil {
		return
	}
	return unitName, unitPtr, nil
}

// Create data container unit
func (c *FleetClient) createDataUnit(component string) (unitName string, unitPtr *unit.Unit, err error) {
	unitName, err = formatUnitName(component, 0)
	if err != nil {
		return
	}
	machineID, err := randomMachineID(c)
	if err != nil {
		return
	}
	unitPtr, err = NewDataUnit(component, machineID)
	if err != nil {
		return
	}
	return unitName, unitPtr, nil

}
