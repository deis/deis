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
var rootPaths = []string{"/var/lib/deis/units", "units"}

// getUnits returns a list of units filtered by target
func (c *FleetClient) getUnits(target string) (units []string, err error) {
	allUnits, err := c.Fleet.Units()
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
	for _, u := range allUnits {
		match := r.MatchString(u.Name)
		if match {
			units = append(units, u.Name)
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
func NewUnit(component string) (uf *unit.UnitFile, err error) {
	template, err := readTemplate(component)
	if err != nil {
		return
	}
	uf, err = unit.NewUnitFile(string(template))
	if err != nil {
		return
	}
	return
}

// NewDataUnit takes a component type and returns a Fleet unit
// that is hard-scheduled to a machine ID
func NewDataUnit(component string, machineID string) (uf *unit.UnitFile, err error) {
	template, err := readTemplate(component)
	if err != nil {
		return
	}
	// replace CHANGEME with random machineID
	replaced := strings.Replace(string(template), "CHANGEME", machineID, 1)
	uf, err = unit.NewUnitFile(replaced)
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
	}
	return "deis-" + component + "@" + strconv.Itoa(num) + ".service", nil
}

// readTemplate returns the contents of a systemd template for the given component
func readTemplate(component string) (out []byte, err error) {
	templateName := "deis-" + component + ".service"
	var templateFile string

	// first look for unit files in GOPATH
	if os.Getenv("GOPATH") != "" {
		filename := path.Join(os.Getenv("GOPATH"),
			path.Join("src", "github.com", "deis", "deisctl", "units", templateName))
		if _, err := os.Stat(filename); err == nil {
			templateFile = filename
		}
	} else {
		// otherwise look in rootPaths hierarchy
		for _, rootPath := range rootPaths {
			filename := path.Join(rootPath, templateName)
			if _, err := os.Stat(filename); err == nil {
				templateFile = filename
				break
			}
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
