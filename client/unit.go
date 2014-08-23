package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/coreos/fleet/unit"
)

// path hierarchy for finding systemd service templates
var rootPaths = []string{"/var/lib/deis/units", "units", "../units"}

// getUnits returns a list of units filtered by target
func (c *FleetClient) getUnits(target string) (units []string, err error) {
	jobs, err := c.Fleet.Jobs()
	if err != nil {
		return
	}
	var r *regexp.Regexp
	if strings.HasSuffix(target, "-data") {
		r = regexp.MustCompile(`deis\-(` + target + `)\.service`)
	} else if strings.Contains(target, ".") {
		r = regexp.MustCompile(`deis\-(` + target + `)\.service`)
	} else {
		r = regexp.MustCompile(`deis\-(` + target + `)@([\d]+)\.service`)
	}
	for _, j := range jobs {
		match := r.MatchString(j.Name)
		if match {
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

// NewUnit takes a component type and returns a Fleet unit
// that includes the relevant systemd service template
func NewUnit(component string) (u *unit.Unit, err error) {
	template, err := readTemplate(component)
	if err != nil {
		return
	}
	u, err = unit.NewUnit(string(template))
	if err != nil {
		return
	}
	return
}

// NewDataUnit takes a component type and returns a Fleet unit
// that is hard-scheduled to a machine ID
func NewDataUnit(component string, machineID string) (u *unit.Unit, err error) {
	template, err := readTemplate(component)
	if err != nil {
		return
	}
	// replace CHANGEME with random machineID
	replaced := strings.Replace(string(template), "CHANGEME", machineID, 1)
	u, err = unit.NewUnit(replaced)
	if err != nil {
		return
	}
	return
}

// formatUnitName returns a properly formatted systemd service name
// using the given component type and number
func formatUnitName(component string, num int) (unitName string, err error) {
	if num == 0 {
		return "deis-" + component + ".service", nil
	} else {
		return "deis-" + component + "@" + strconv.Itoa(num) + ".service", nil
	}
}

// readTemplate returns the contents of a systemd template for the given component
func readTemplate(component string) (out []byte, err error) {
	templateName := "deis-" + component + ".service"
	var templateFile string
	for _, rootPath := range rootPaths {
		filename := path.Join(rootPath, templateName)
		if _, err := os.Stat(filename); err == nil {
			templateFile = filename
			break
		}
	}
	if templateFile == "" {
		return nil, fmt.Errorf("Could not find unit template for %v", component)
	}
	out, err = ioutil.ReadFile(templateFile)
	if err != nil {
		return
	}
	return
}
