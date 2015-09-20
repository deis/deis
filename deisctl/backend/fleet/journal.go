package fleet

import (
	"fmt"
	"io"
	"time"
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

// JournalBackground runs Journal in the background
func (c *FleetClient) JournalBackground(targets string, out io.Writer) []chan bool {
	var quitChans []chan bool

	expandedTargets, err := c.expandTargets([]string{targets})
	if err != nil {
		fmt.Fprintln(out, err)
		return []chan bool{}
	}
	for _, target := range expandedTargets {
		quit := c.runJournalBackground(target, out)

		if quit != nil {
			quitChans = append(quitChans, quit)
		}
	}
	return quitChans
}

func (c *FleetClient) runJournalBackground(target string, out io.Writer) chan bool {
	machineID, err := c.findUnit(target)

	// Return if the unit can't be found.
	if err != nil {
		return nil
	}

	ms, err := c.machineState(machineID)
	if err != nil || ms == nil {
		fmt.Fprintf(out, "Error getting machine IP: %v\n", err)
		return nil
	}

	tick := time.Tick(c.journalInterval)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-tick:
				layout := "2006-01-02 15:04:05"
				now := c.timeNow().Add(-c.journalInterval).UTC().Format(layout)
				command := fmt.Sprintf(`journalctl --unit %s --no-pager -o cat --since="%s"`, target, now)
				log, err := c.runRemoteCommandString(command, ms.PublicIP)

				if err != nil {
					fmt.Fprintf(out, "Deisctl Error: %s\n", err.Error())
					continue
				}

				// Journalctl doesn't have logs yet
				if log == "Failed to get data: Cannot assign requested address\n" {
					continue
				}
				fmt.Fprint(out, log)
			case <-quit:
				close(quit)
				return
			}
		}
	}()

	return quit
}
