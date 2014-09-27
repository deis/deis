package fleet

import (
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/schema"
)

// Stop sets target units to inactive and blocks until complete
func (c *FleetClient) Stop(targets []string) error {
	units := make([][]string, len(targets))
	for i, target := range targets {
		u, err := c.Units(target)
		if err != nil {
			return err
		}
		units[i] = u
	}
	desiredState := string(job.JobStateLoaded)
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
			outchan, errchan = waitForUnitSubStates(names, "dead")
			if err := printUnitSubState(name, outchan, errchan); err != nil {
				return err
			}
		}
	}
	return nil

}
