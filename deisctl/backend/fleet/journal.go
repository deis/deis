package fleet

import (
	"fmt"
	"os"
)

// Journal prints the systemd journal of target unit(s)
func (c *FleetClient) Journal(target string) (err error) {
	units, err := c.Units(target)
	if err != nil {
		return
	}
	for _, unit := range units {
		runJournal(unit)
	}
	return
}

// runJournal tails the systemd journal for a given unit
func runJournal(name string) (exit int) {

	u, err := cAPI.Unit(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving Unit %s: %v", name, err)
		return 1
	}
	if u == nil {
		fmt.Fprintf(os.Stderr, "Unit %s does not exist.\n", name)
		return 1
	} else if u.CurrentState == "" {
		fmt.Fprintf(os.Stderr, "Unit %s does not appear to be running.\n", name)
		return 1
	}

	command := fmt.Sprintf("journalctl --unit %s --no-pager -n 40 -f", name)
	return runCommand(command, u.MachineID)
}
