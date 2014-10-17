package fleet

import (
	"fmt"
	"strings"
	"sync"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"
)

// Create schedules a new unit for the given component
func (c *FleetClient) Create(targets []string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	units := make([]*schema.Unit, len(targets))

	for i, target := range targets {
		unitName, unitFile, err := c.createUnitFile(target)
		if err != nil {
			errchan <- err
			return
		}
		units[i] = &schema.Unit{
			Name:    unitName,
			Options: schema.MapUnitFileToSchemaUnitOptions(unitFile),
		}
	}

	for _, unit := range units {
		wg.Add(1)
		go doCreate(c, unit, wg, outchan, errchan)
	}
}

func doCreate(c *FleetClient, unit *schema.Unit, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	defer wg.Done()

	// create unit definition
	if err := c.Fleet.CreateUnit(unit); err != nil {
		// ignore units that already exist
		if err.Error() != "job already exists" {
			errchan <- err
			return
		}
	}

	desiredState := string(job.JobStateLoaded)
	out := fmt.Sprintf("\033[0;33m%v:\033[0m loaded                                 \r", unit.Name)

	// schedule the unit
	if err := c.Fleet.SetUnitTargetState(unit.Name, desiredState); err != nil {
		errchan <- err
		return
	}

	outchan <- out
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
