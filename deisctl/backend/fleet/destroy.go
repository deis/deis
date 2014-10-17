package fleet

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Destroy units for a given target
func (c *FleetClient) Destroy(targets []string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	for _, target := range targets {
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

	// bail early if unit doesn't exist
	_, err = c.Units(name)
	if err != nil {
		if strings.Contains(err.Error(), "could not find unit") {
			outchan <- destroyed
		}
		return
	}

	// otherwise destroy it
	if err = c.Fleet.DestroyUnit(name); err != nil {
		// ignore already destroyed units
		if !strings.Contains(err.Error(), "could not find unit") {
			errchan <- err
			return
		}
	}

	// loop until actually destroyed
	for {
		_, err = c.Units(name)
		if err != nil {
			if strings.Contains(err.Error(), "could not find unit") {
				outchan <- destroyed
				return
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
}
