package fleet

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
)

// Scale creates or destroys units to match the desired number
func (c *FleetClient) Scale(
	component string, requested int, wg *sync.WaitGroup, outchan chan string, errchan chan error) {

	if requested < 0 {
		errchan <- errors.New("cannot scale below 0")
	}
	// check how many currently exist
	components, err := c.Units(component)
	if err != nil {
		// skip checking the first time; we just want a tally
		if !strings.Contains(err.Error(), "could not find unit") {
			errchan <- err
			return
		}
	}

	timesToScale := int(math.Abs(float64(requested - len(components))))
	if timesToScale == 0 {
		return
	}
	if requested-len(components) > 0 {
		scaleUp(c, component, len(components), timesToScale, wg, outchan, errchan)
	} else {
		scaleDown(c, component, len(components), timesToScale, wg, outchan, errchan)
	}
}

func scaleUp(c *FleetClient, component string, numExistingContainers, numTimesToScale int,
	wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	for i := 0; i < numTimesToScale; i++ {
		target := component + "@" + strconv.Itoa(numExistingContainers+i+1)
		c.Create([]string{target}, wg, outchan, errchan)
	}
}

func scaleDown(c *FleetClient, component string, numExistingContainers, numTimesToScale int,
	wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	for i := 0; i < numTimesToScale; i++ {
		target := component + "@" + strconv.Itoa(numExistingContainers-i)
		c.Destroy([]string{target}, wg, outchan, errchan)
	}
}
