package client

import (
	"fmt"
)

// Destroy units for a given target
func (c *FleetClient) Destroy(target string) (err error) {
	component, num, err := splitTarget(target)
	if err != nil {
		return
	}
	if num == 0 {
		num, err = c.lastUnit(component)
		if err != nil {
			return err
		}
	}
	unitName, err := formatUnitName(component, num)
	if err != nil {
		return err
	}
	_, err = c.Fleet.Unit(unitName)
	if err != nil {
		return
	}
	if err = c.Fleet.DestroyUnit(unitName); err != nil {
		return fmt.Errorf("failed destroying job %s: %v", unitName, err)
	}
	fmt.Printf("Destroyed Unit %s\n", unitName)
	return
}
