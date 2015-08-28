package fleet

import (
	"fmt"
)

// Journal prints the systemd journal of target unit(s)
func (c *FleetClient) Journal(target string) (err error) {
	units, err := c.Units(target)
	if err != nil {
		return
	}
	for _, unit := range units {
		c.runJournal(unit)
	}
	return
}

// runJournal tails the systemd journal for a given unit
func (c *FleetClient) runJournal(name string) (exit int) {
	u, err := c.Fleet.Unit(name)
	if suToGlobal(*u) {
		fmt.Fprintf(c.errWriter, "Unable to get journal for global unit %s. Check on a host directly using journalctl.\n", name)
		return 1
	}

	machineID, err := c.findUnit(name)

	if err != nil {
		return 1
	}

	command := fmt.Sprintf("journalctl --unit %s --no-pager -n 40 -f", name)
	return c.runCommand(command, machineID)
}
