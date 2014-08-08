package deisctl

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/job"
)

// Client interface used to interact with the cluster control plane
type Client interface {
	Create(string) error
	Destroy(string) error
	Start(string) error
	Stop(string) error
	Scale(string, int) error
	List() error
	Status(string) error
}

// FleetClient used to wrap Fleet API calls
type FleetClient struct {
	Fleet client.API
}

// NewClient returns a client used to communicate with Fleet
// using the Registry API
func NewClient() (*FleetClient, error) {
	client, err := getRegistryClient()
	if err != nil {
		return nil, err
	}
	return &FleetClient{Fleet: client}, nil
}

// Create schedules a new unit for the given component
// and blocks until the unit is loaded
func (c *FleetClient) Create(component string) (err error) {
	num, err := c.nextUnit(component)
	if err != nil {
		return
	}
	unitName, err := formatUnitName(component, num)
	if err != nil {
		return
	}
	unit, err := NewUnit(component)
	if err != nil {
		return
	}
	j := job.NewJob(unitName, *unit)
	if err := c.Fleet.CreateJob(j); err != nil {
		return fmt.Errorf("failed creating job %s: %v", unitName, err)
	}
	newState := job.JobStateLoaded
	err = c.Fleet.SetJobTargetState(unitName, newState)
	if err != nil {
		return err
	}
	errchan := waitForJobStates(c.Fleet, []string{unitName}, testJobStateLoaded, 0, os.Stdout)
	for err := range errchan {
		return fmt.Errorf("error waiting for job %s: %v", unitName, err)
	}
	return nil
}

// Destroy unschedules one unit for a given component type
func (c *FleetClient) Destroy(component string) (err error) {
	num, err := c.lastUnit(component)
	if err != nil {
		return
	}
	if num == 0 {
		return fmt.Errorf("no units to destroy")
	}
	unitName, err := formatUnitName(component, num)
	if err != nil {
		return
	}
	_, err = c.Fleet.Job(unitName)
	if err != nil {
		return
	}
	if err = c.Fleet.DestroyJob(unitName); err != nil {
		return fmt.Errorf("failed destroying job %s: %v", unitName, err)
	}
	fmt.Printf("Destroyed Unit %s\n", unitName)
	return
}

// Scale creates or destroys units to match the desired number
func (c *FleetClient) Scale(component string, num int) (err error) {
	for {
		components, err := c.getUnits(component)
		if err != nil {
			return err
		}
		if len(components) == num {
			break
		}
		if len(components) < num {
			c.Create(component)
			continue
		}
		if len(components) > num {
			c.Destroy(component)
			continue
		}
	}
	return
}

// Start launches target units and blocks until active
func (c *FleetClient) Start(target string) (err error) {
	units, err := c.getUnits(target)
	if err != nil {
		return
	}
	newState := job.JobStateLaunched
	for _, unitName := range units {
		err = c.Fleet.SetJobTargetState(unitName, newState)
		if err != nil {
			return err
		}
	}
	errchan := waitForJobStates(c.Fleet, units, testUnitStateActive, 0, os.Stdout)
	for err := range errchan {
		return fmt.Errorf("error waiting for active: %v", err)
	}
	return nil
}

// Stop sets target units to inactive and blocks until complete
func (c *FleetClient) Stop(target string) (err error) {
	units, err := c.getUnits(target)
	if err != nil {
		return
	}
	newState := job.JobStateInactive
	for _, unitName := range units {
		err = c.Fleet.SetJobTargetState(unitName, newState)
		if err != nil {
			return err
		}
	}
	errchan := waitForJobStates(c.Fleet, units, testJobStateInactive, 0, os.Stdout)
	for err := range errchan {
		return fmt.Errorf("error waiting for inactive: %v", err)
	}
	return nil
}

// List prints all Deis-related units to Stdout
func (c *FleetClient) List() (err error) {

	var jobs map[string]job.Job
	var sortable sort.StringSlice

	jobs = make(map[string]job.Job, 0)
	jj, err := c.Fleet.Jobs()
	if err != nil {
		return err
	}
	for _, j := range jj {
		if strings.HasPrefix(j.Name, "deis-") {
			jobs[j.Name] = j
			sortable = append(sortable, j.Name)
		}
	}
	sortable.Sort()
	printList(jobs, sortable)
	return
}

// Status prints the systemd status of target unit(s)
func (c *FleetClient) Status(target string) (err error) {
	units, err := c.getUnits(target)
	if err != nil {
		return
	}
	for _, unit := range units {
		printUnitStatus(c.Fleet, unit)
		fmt.Println()
	}
	return
}

// getUnits returns a list of units filtered by target
func (c *FleetClient) getUnits(target string) (units []string, err error) {
	jobs, err := c.Fleet.Jobs()
	if err != nil {
		return
	}
	for _, j := range jobs {
		if strings.HasPrefix(j.Name, "deis-"+target) {
			units = append(units, j.Name)
		}
	}
	return
}

// nextUnit returns the next unit number for a given component
func (c *FleetClient) nextUnit(component string) (num int, err error) {
	units, err := c.getUnits(component)
	if err != nil {
		return
	}
	num, err = nextUnitNum(units)
	if err != nil {
		return
	}
	return
}

// lastUnit returns the last unit number for a given component
func (c *FleetClient) lastUnit(component string) (num int, err error) {
	units, err := c.getUnits(component)
	if err != nil {
		return
	}
	num, err = lastUnitNum(units)
	if err != nil {
		return
	}
	return
}
