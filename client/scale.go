package client

import "strconv"

// Scale creates or destroys units to match the desired number
func (c *FleetClient) Scale(component string, num int) (err error) {
	for {
		components, err := c.Units(component)
		if err != nil {
			return err
		}
		if len(components) == num {
			break
		}
		if len(components) < num {
			num, err = c.nextUnit(component)
			if err != nil {
				return err
			}
			err := c.Create(component + "@" + strconv.Itoa(num))
			if err != nil {
				return err
			}
			continue
		}
		if len(components) > num {
			num, err = c.lastUnit(component)
			if err != nil {
				return err
			}
			err := c.Destroy(component + "@" + strconv.Itoa(num))
			if err != nil {
				return err
			}
			continue
		}
	}
	return
}
