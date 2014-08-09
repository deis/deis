package client

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/unit"
)

// Create schedules a new unit for the given component
// and blocks until the unit is loaded
func (c *FleetClient) Create(component string, data bool) (err error) {
	var (
		unitName string
		unitPtr  *unit.Unit
	)
	// create unit
	if data == true {
		unitName, unitPtr, err = c.createDataUnit(component)
	} else {
		unitName, unitPtr, err = c.createServiceUnit(component)
	}
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
	errchan := waitForJobStates([]string{unitName}, testJobStateLoaded, 0, os.Stdout)
	for err := range errchan {
		return fmt.Errorf("error waiting for job %s: %v", unitName, err)
	}
	return nil
}

// Create normal service unit
func (c *FleetClient) createServiceUnit(target string) (unitName string, unitPtr *unit.Unit, err error) {

	// see if we were provided a specific target
	r := regexp.MustCompile(`([a-z-]+)\.([\d]+)`)
	match := r.FindStringSubmatch(target)
	var (
		num int
		component string
	)
	if len(match) == 3 {
		component = match[1]
		num, err = strconv.Atoi(match[2])
		if err != nil {
			return "", nil, err
		}
		unitName, err = formatUnitName(component, num)
	} else {
		component = target
		num, err := c.nextUnit(component)
		if err != nil {
			return "", nil, err
		}
		unitName, err = formatUnitName(component, num)
		if err != nil {
			return "", nil, err
		}
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
