package fleet

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
)

// Scale creates or destroys units to match the desired number
func (c *FleetClient) Scale(
	component string, requested int, wg *sync.WaitGroup, out, ew io.Writer) {

	if requested < 0 {
		fmt.Fprintln(ew, "cannot scale below 0")
		return
	}
	// check how many currently exist
	components, err := c.Units(component)
	if err != nil {
		// skip checking the first time; we just want a tally
		if !strings.Contains(err.Error(), "could not find unit") {
			fmt.Fprintln(ew, err.Error())
			return
		}
	}

	timesToScale := int(math.Abs(float64(requested - len(components))))
	switch {
	case timesToScale == 0:
		return
	case requested-len(components) > 0:
		c.scaleUp(component, len(components), timesToScale, wg, out, ew)
	default:
		c.scaleDown(component, len(components), timesToScale, wg, out, ew)
	}
}

func (c *FleetClient) scaleUp(component string, numExistingContainers, numTimesToScale int,
	wg *sync.WaitGroup, out, ew io.Writer) {
	for i := 0; i < numTimesToScale; i++ {
		target := component + "@" + strconv.Itoa(numExistingContainers+i+1)
		c.Create([]string{target}, wg, out, ew)
	}
	wg.Wait()
	for i := 0; i < numTimesToScale; i++ {
		target := component + "@" + strconv.Itoa(numExistingContainers+i+1)
		c.Start([]string{target}, wg, out, ew)
	}
}

func (c *FleetClient) scaleDown(component string, numExistingContainers, numTimesToScale int,
	wg *sync.WaitGroup, out, ew io.Writer) {
	for i := 0; i < numTimesToScale; i++ {
		target := component + "@" + strconv.Itoa(numExistingContainers-i)
		c.Destroy([]string{target}, wg, out, ew)
	}
}
