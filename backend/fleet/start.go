package fleet

import (
	"strings"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/schema"
)

// Start launches target units and blocks until active
func (c *FleetClient) Start(targets []string) error {
	units := make([][]string, len(targets))
	for i, target := range targets {
		u, err := c.Units(target)
		if err != nil {
			return err
		}
		units[i] = u
	}
	desiredState := string(job.JobStateLaunched)
	for _, names := range units {
		for _, name := range names {
			if err := c.Fleet.SetUnitTargetState(name, desiredState); err != nil {
				return err
			}
		}
	}
	var errchan chan error
	var outchan chan *schema.Unit
	for _, names := range units {
		for _, name := range names {
			// wait for systemd to tell us that it's running, not fleet
			// data containers are special snowflakes who just exit
			if strings.Contains(name, "-data.service") {
				outchan, errchan = waitForUnitSubStates(names, "exited")
			} else {
				outchan, errchan = waitForUnitSubStates(names, "running")
			}
			if err := printUnitSubState(name, outchan, errchan); err != nil {
				return err
			}
		}
	}
	return nil
}
