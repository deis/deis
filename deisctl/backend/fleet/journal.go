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
		runJournal(c, unit)
	}
	return
}

// runJournal tails the systemd journal for a given unit
func runJournal(c *FleetClient, name string) (exit int) {
	machineID, err := findUnit(c, name)

	if err != nil {
		return 1
	}

	command := fmt.Sprintf("journalctl --unit %s --no-pager -n 40 -f", name)
	return runCommand(c, command, machineID)
}
