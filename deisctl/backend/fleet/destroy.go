package fleet

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/deis/deis/pkg/prettyprint"
)

// Destroy units for a given target
func (c *FleetClient) Destroy(targets []string, wg *sync.WaitGroup, out, ew io.Writer) {
	// expand @* targets
	expandedTargets, err := c.expandTargets(targets)
	if err != nil {
		fmt.Fprintln(ew, err.Error())
		return
	}

	for _, target := range expandedTargets {
		wg.Add(1)
		go doDestroy(c, target, wg, out, ew)
	}
	return
}

func doDestroy(c *FleetClient, target string, wg *sync.WaitGroup, out, ew io.Writer) {
	defer wg.Done()

	// prepare string representation
	component, num, err := splitTarget(target)
	if err != nil {
		fmt.Fprintln(ew, err.Error())
		return
	}
	name, err := formatUnitName(component, num)
	if err != nil {
		fmt.Fprintln(ew, err.Error())
		return
	}
	tpl := prettyprint.Colorize("{{.Yellow}}%v:{{.Default}} destroyed")
	destroyed := fmt.Sprintf(tpl, name)

	// tell fleet to destroy the unit
	c.Fleet.DestroyUnit(name)

	// loop until the unit is actually gone from unit states
outerLoop:
	for {
		time.Sleep(250 * time.Millisecond)
		unitStates, err := c.Fleet.UnitStates()
		if err != nil {
			fmt.Fprintln(ew, err.Error())
		}
		for _, us := range unitStates {
			if strings.HasPrefix(us.Name, name) {
				continue outerLoop
			}
		}
		fmt.Fprintln(out, destroyed)
		return
	}
}
