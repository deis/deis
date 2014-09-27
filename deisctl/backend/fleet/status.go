package fleet

import (
	"fmt"
	"os"
)

// Status prints the systemd status of target unit(s)
func (c *FleetClient) Status(target string) (err error) {
	units, err := c.Units(target)
	if err != nil {
		return
	}
	for _, unit := range units {
		printUnitStatus(unit)
		fmt.Println()
	}
	return
}

// printUnitStatus displays the systemd status for a given unit
func printUnitStatus(name string) int {
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
	cmd := fmt.Sprintf("systemctl status -l %s", name)
	return runCommand(cmd, u.MachineID)
}
