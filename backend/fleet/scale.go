package fleet

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

// Scale creates or destroys units to match the desired number
func (c *FleetClient) Scale(component string, requested int) error {
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
	if requested-len(components) > 0 {
		return scaleUp(c, component, len(components), timesToScale)
	} else {
		return scaleDown(c, component, len(components), timesToScale)
	}
}

func scaleUp(c *FleetClient, component string, numExistingContainers, numTimesToScale int) error {
	for i := 0; i < numTimesToScale; i++ {
		if err := c.Create([]string{component + "@" + strconv.Itoa(numExistingContainers+i+1)}); err != nil {
			return err
		}
	}
	return nil
}

func scaleDown(c *FleetClient, component string, numExistingContainers, numTimesToScale int) error {
	for i := 0; i < numTimesToScale; i++ {
		if err := c.Destroy([]string{component + "@" + strconv.Itoa(numExistingContainers-i)}); err != nil {
			return err
		}
	}
	return nil
}
