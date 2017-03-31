package fleet

import (
	"fmt"
	"github.com/deis/deis/deisctl/units"
	"strings"
)

//checking if all units are active
func (c *FleetClient) isPlatformRunning() (running bool) {
	unitStates, _ := c.Fleet.UnitStates() //@todo: probs wanna check error

	for _, us := range unitStates {
		for _, prefix := range units.Names {
			if strings.HasPrefix(us.Name, prefix) {
				if us.SystemdActiveState != "active" {
					return false
				}
			}
		}
	}

	return true
}

func (c *FleetClient) Doctor() (err error) {
	fmt.Println("Calculating some things ...")
	// @todo: compare deisctl version against deis version

	// confirm all units are active
	running := c.isPlatformRunning()
	if running != true {
		fmt.Println("\tAll units are not active")
	} else {
		fmt.Println("\tThumbs up! All units are active.")
	}

	// @todo: check to see if all image tags match

	return
}
