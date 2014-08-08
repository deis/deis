package client

// Scale creates or destroys units to match the desired number
func (c *FleetClient) Scale(component string, num int) (err error) {
	for {
		components, err := c.getUnits(component)
		if err != nil {
			return err
		}
		if len(components) == num {
			break
		}
		if len(components) < num {
			err := c.Create(component, false)
			if err != nil {
				return err
			}
			continue
		}
		if len(components) > num {
			err := c.Destroy(component)
			if err != nil {
				return err
			}
			continue
		}
	}
	return
}
