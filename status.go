package deisctl

import (
	"fmt"
	"os"

	"github.com/coreos/fleet/client"
)

// printUnitStatus displays the systemd status for a given unit
func printUnitStatus(cAPI client.API, jobName string) int {
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
	cmd := fmt.Sprintf("systemctl status -l %s", jobName)
	return runCommand(cmd, j.UnitState.MachineState)
}
