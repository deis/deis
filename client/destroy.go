package client

import (
	"fmt"
	"regexp"
	"strconv"
)

// Destroy unschedules one unit for a given component type
func (c *FleetClient) Destroy(component string) (err error) {

	// see if we were provided a specific target
	r := regexp.MustCompile(`([a-z-]+)\.([\d]+)`)
	match := r.FindStringSubmatch(component)
	var (
		num int
		unitName string
	)
	if len(match) == 3 {
		num, err = strconv.Atoi(match[2])
		if err != nil {
			return err
		}
		unitName, err = formatUnitName(component, 0)
	} else {
		num, err = c.lastUnit(component)
		if err != nil {
			return err
		}
		unitName, err = formatUnitName(component, num)
		if err != nil {
			return err
		}
	}
	if num == 0 {
		return fmt.Errorf("no units to destroy")
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
