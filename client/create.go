package client

import (
	"fmt"
	"strings"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

// Create schedules a new unit for the given component
// and blocks until the unit is loaded
func (c *FleetClient) Create(target string) (err error) {
	var (
		unitName string
		unitFile *unit.UnitFile
	)
	// create unit
	unitName, unitFile, err = c.createUnitFile(target)
	if err != nil {
		return err
	}
	//
	u := &schema.Unit{
		Name:    unitName,
		Options: schema.MapUnitFileToSchemaUnitOptions(unitFile),
	}
	// schedule unit
	if err := c.Fleet.CreateUnit(u); err != nil {
		// ignore units that already exist
		if err.Error() != "job already exists" {
			return fmt.Errorf("failed creating job %s: %v", unitName, err)
		}
	}
	desiredState := string(job.JobStateLoaded)
	err = c.Fleet.SetUnitTargetState(unitName, desiredState)
	if err != nil {
		return err
	}
	outchan, errchan := waitForUnitStates([]string{unitName}, desiredState)
	err = printUnitState(unitName, outchan, errchan)
	if err != nil {
		return err
	}
	return nil
}

func (c *FleetClient) createUnitFile(target string) (unitName string, uf *unit.UnitFile, err error) {
	component, num, err := splitTarget(target)
	if err != nil {
		return
	}
	if strings.HasSuffix(component, "-data") {
		unitName, uf, err = c.createDataUnit(component)
	} else {
		unitName, uf, err = c.createServiceUnit(component, num)
	}
	if err != nil {
		return unitName, uf, err
	}
	return
}

// Create normal service unit
func (c *FleetClient) createServiceUnit(component string, num int) (name string, uf *unit.UnitFile, err error) {
	// if number wasn't provided get next unit number
	if num == 0 {
		num, err = c.nextUnit(component)
		if err != nil {
			return "", nil, err
		}
	}
	// build a fleet unit
	name, err = formatUnitName(component, num)
	if err != nil {
		return "", nil, err
	}
	uf, err = NewUnit(component)
	if err != nil {
		return
	}
	return name, uf, nil
}

// Create data container unit
func (c *FleetClient) createDataUnit(component string) (name string, uf *unit.UnitFile, err error) {
	name, err = formatUnitName(component, 0)
	if err != nil {
		return
	}
	machineID, err := randomMachineID(c)
	if err != nil {
		return
	}
	uf, err = NewDataUnit(component, machineID)
	if err != nil {
		return
	}
	return name, uf, nil

}
