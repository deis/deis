package client

import (
	"fmt"
	"os"
)

// Journal prints the systemd journal of target unit(s)
func (c *FleetClient) Journal(target string) (err error) {
	units, err := c.getUnits(target)
	if err != nil {
		return
	}
	for _, unit := range units {
		runJournal(unit)
	}
	return
}

// runJournal tails the systemd journal for a given unit
func runJournal(jobName string) (exit int) {

	j, err := cAPI.Job(jobName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving Job %s: %v", jobName, err)
		return 1
	}
	if j == nil {
		fmt.Fprintf(os.Stderr, "Job %s does not exist.\n", jobName)
		return 1
	} else if j.UnitState == nil {
		fmt.Fprintf(os.Stderr, "Job %s does not appear to be running.\n", jobName)
		return 1
	}

	command := fmt.Sprintf("journalctl --unit %s --no-pager -n 40 -f", jobName)
	return runCommand(command, j.UnitState.MachineID)
}
