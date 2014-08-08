package client

import "fmt"

// Destroy unschedules one unit for a given component type
func (c *FleetClient) Destroy(component string) (err error) {
	num, err := c.lastUnit(component)
	if err != nil {
		return
	}
	if num == 0 {
		return fmt.Errorf("no units to destroy")
	}
	unitName, err := formatUnitName(component, num)
	if err != nil {
		return
	}
	_, err = c.Fleet.Job(unitName)
	if err != nil {
		return
	}
	if err = c.Fleet.DestroyJob(unitName); err != nil {
		return fmt.Errorf("failed destroying job %s: %v", unitName, err)
	}
	fmt.Printf("Destroyed Unit %s\n", unitName)
	return
}
