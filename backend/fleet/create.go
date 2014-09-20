package fleet

import (
	"fmt"
	"strings"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

// Create schedules a new unit for the given component
// and blocks until the unit is loaded
func (c *FleetClient) Create(targets []string) error {
	units := make([]*schema.Unit, len(targets))
	for i, target := range targets {
		unitName, unitFile, err := c.createUnitFile(target)
		if err != nil {
			return err
		}
		units[i] = &schema.Unit{
			Name:    unitName,
			Options: schema.MapUnitFileToSchemaUnitOptions(unitFile),
		}
	}
	for _, unit := range units {
		// schedule unit
		if err := c.Fleet.CreateUnit(unit); err != nil {
			// ignore units that already exist
			if err.Error() != "job already exists" {
				return fmt.Errorf("failed creating job %s: %v", unit.Name, err)
			}
		}
		desiredState := string(job.JobStateLoaded)
		if err := c.Fleet.SetUnitTargetState(unit.Name, desiredState); err != nil {
			return err
		}
	}
	for _, unit := range units {
		desiredState := string(job.JobStateLoaded)
		outchan, errchan := waitForUnitStates([]string{unit.Name}, desiredState)
		if err := printUnitState(unit.Name, outchan, errchan); err != nil {
			return err
		}
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
