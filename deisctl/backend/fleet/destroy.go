package fleet

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Destroy units for a given target
func (c *FleetClient) Destroy(targets []string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	// expand @* targets
	expandedTargets, err := expandTargets(c, targets)
	if err != nil {
		errchan <- err
		return
	}

	for _, target := range expandedTargets {
		wg.Add(1)
		go doDestroy(c, target, wg, outchan, errchan)
	}
	return
}

func doDestroy(c *FleetClient, target string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	defer wg.Done()

	// prepare string representation
	component, num, err := splitTarget(target)
	if err != nil {
		errchan <- err
		return
	}
	name, err := formatUnitName(component, num)
	if err != nil {
		errchan <- err
		return
	}
	destroyed := fmt.Sprintf("\033[0;33m%v:\033[0m destroyed                                 \r", name)

	// tell fleet to destroy the unit
	c.Fleet.DestroyUnit(name)

	// loop until the unit is actually gone from unit states
outerLoop:
	for {
		time.Sleep(250 * time.Millisecond)
		unitStates, err := cAPI.UnitStates()
		if err != nil {
			errchan <- err
		}
		for _, us := range unitStates {
			if strings.HasPrefix(us.Name, name) {
				continue outerLoop
			}
		}
		outchan <- destroyed
		return
	}
}
