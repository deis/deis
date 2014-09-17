package client

import (
	"strings"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/schema"
)

// Start launches target units and blocks until active
func (c *FleetClient) Start(target string) (err error) {
	units, err := c.Units(target)
	if err != nil {
		return
	}
	desiredState := string(job.JobStateLaunched)
	for _, name := range units {
		err = c.Fleet.SetUnitTargetState(name, desiredState)
		if err != nil {
			return err
		}
		// wait for systemd to tell us that it's running, not fleet
		var errchan chan error
		var outchan chan *schema.Unit
		// data containers are special snowflakes who just exit
		if strings.Contains(name, "-data.service") {
			outchan, errchan = waitForUnitSubStates(units, "exited")
		} else {
			outchan, errchan = waitForUnitSubStates(units, "running")
		}
		err = printUnitSubState(name, outchan, errchan)
		if err != nil {
			return err
		}
	}
	return nil
}
