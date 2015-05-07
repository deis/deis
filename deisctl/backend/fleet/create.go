package fleet

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/schema"
	"github.com/coreos/fleet/unit"

	"github.com/deis/deis/pkg/prettyprint"
)

// Create schedules unit files for the given components.
func (c *FleetClient) Create(
	targets []string, wg *sync.WaitGroup, out, ew io.Writer) {

	units := make([]*schema.Unit, len(targets))

	for i, target := range targets {
		unitName, unitFile, err := c.createUnitFile(target)
		if err != nil {
			fmt.Fprintf(ew, "Error creating: %s\n", err)
			return
		}
		units[i] = &schema.Unit{
			Name:    unitName,
			Options: schema.MapUnitFileToSchemaUnitOptions(unitFile),
		}
	}

	for _, unit := range units {
		wg.Add(1)
		go doCreate(c, unit, wg, out, ew)
	}
}

func doCreate(c *FleetClient, unit *schema.Unit, wg *sync.WaitGroup, out, ew io.Writer) {
	defer wg.Done()

	// create unit definition
	if err := c.Fleet.CreateUnit(unit); err != nil {
		// ignore units that already exist
		if err.Error() != "job already exists" {
			fmt.Fprintln(ew, err.Error())
			return
		}
	}

	desiredState := string(job.JobStateLoaded)
	tpl := prettyprint.Colorize("{{.Yellow}}%v:{{.Default}} loaded")
	msg := fmt.Sprintf(tpl, unit.Name)

	// schedule the unit
	if err := c.Fleet.SetUnitTargetState(unit.Name, desiredState); err != nil {
		fmt.Fprintln(ew, err)
		return
	}

	// loop until the unit actually exists in unit states
outerLoop:
	for {
		time.Sleep(250 * time.Millisecond)
		unitStates, err := c.Fleet.UnitStates()
		if err != nil {
			fmt.Fprintln(ew, err)
		}
		for _, us := range unitStates {
			if strings.HasPrefix(us.Name, unit.Name) {
				break outerLoop
			}
		}
	}

	fmt.Fprintln(out, msg)
}

func (c *FleetClient) createUnitFile(target string) (unitName string, uf *unit.UnitFile, err error) {
	component, num, err := splitTarget(target)
	if err != nil {
		return
	}
	unitName, uf, err = c.createServiceUnit(component, num)
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
	decorateStr, err := c.configBackend.GetWithDefault("/deis/platform/enablePlacementOptions", "false")
	if err != nil {
		return "", nil, err
	}
	decorate, err := strconv.ParseBool(decorateStr)
	if err != nil {
		return "", nil, err
	}
	uf, err = NewUnit(component, c.templatePaths, decorate)
	if err != nil {
		return
	}
	return name, uf, nil
}
