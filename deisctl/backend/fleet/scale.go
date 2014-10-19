package fleet

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
)

// Scale creates or destroys units to match the desired number
func (c *FleetClient) Scale(component string, requested int) error {

	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	if requested < 0 {
		return errors.New("cannot scale below 0")
	}
	// check how many currently exist
	components, err := c.Units(component)
	if err != nil {
		// skip checking the first time; we just want a tally
		if !strings.Contains(err.Error(), "could not find unit") {
			return err
		}
	}

	timesToScale := int(math.Abs(float64(requested - len(components))))
	if timesToScale == 0 {
		return nil
	}
	if requested-len(components) > 0 {
		return scaleUp(c, component, len(components), timesToScale, &wg, outchan, errchan)
	}
	return scaleDown(c, component, len(components), timesToScale, &wg, outchan, errchan)
}

func scaleUp(c *FleetClient, component string, numExistingContainers, numTimesToScale int,
	wg *sync.WaitGroup, outchan chan string, errchan chan error) error {
	for i := 0; i < numTimesToScale; i++ {
		target := component + "@" + strconv.Itoa(numExistingContainers+i+1)
		c.Create([]string{target}, wg, outchan, errchan)
	}
	return nil
}

func scaleDown(c *FleetClient, component string, numExistingContainers, numTimesToScale int,
	wg *sync.WaitGroup, outchan chan string, errchan chan error) error {
	for i := 0; i < numTimesToScale; i++ {
		target := component + "@" + strconv.Itoa(numExistingContainers-i)
		c.Destroy([]string{target}, wg, outchan, errchan)
	}
	return nil
}
