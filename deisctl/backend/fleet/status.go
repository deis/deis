package fleet

import (
	"fmt"
)

// Status prints the systemd status of target unit(s)
func (c *FleetClient) Status(target string) (err error) {
	units, err := c.Units(target)
	if err != nil {
		return
	}
	for _, unit := range units {
		c.printUnitStatus(unit)
		fmt.Println()
	}
	return
}

// printUnitStatus displays the systemd status for a given unit
func (c *FleetClient) printUnitStatus(name string) int {
	u, err := c.Fleet.Unit(name)
	switch {
	case suToGlobal(*u):
		fmt.Fprintf(c.errWriter, "Unable to get status for global unit %s. Check the status on the host using systemctl.\n", name)
		return 1
	case err != nil:
		fmt.Fprintf(c.errWriter, "Error retrieving Unit %s: %v\n", name, err)
		return 1
	case u == nil:
		fmt.Fprintf(c.errWriter, "Unit %s does not exist.\n", name)
		return 1
	case u.CurrentState == "":
		fmt.Fprintf(c.errWriter, "Unit %s does not appear to be running.\n", name)
		return 1
	}
	cmd := fmt.Sprintf("systemctl status -l %s", name)
	return c.runCommand(cmd, u.MachineID)
}
