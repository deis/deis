package deisctl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"

	"github.com/coreos/fleet/unit"
)

// path hierarchy for finding systemd service templates
var rootPaths = []string{"/var/lib/deis/units", "units", "../units"}

// NewUnit takes a component type and returns a Fleet unit
// that includes the relevant systemd service template
func NewUnit(component string, templateName string, c *FleetClient) (u *unit.Unit, err error) {
	template, err := readTemplate(component, templateName)
	if err != nil {
		return
	}
	template = processTemplate(template, c)
	u, err = unit.NewUnit(string(template))
	if err != nil {
		return
	}
	return
}

// formatUnitName returns a properly formatted systemd service name
// using the given component type and number
func formatUnitName(component string, num int) (unitName string, err error) {
	unitName = "deis-" + component + "." + strconv.Itoa(num) + ".service"
	return
}

// readTemplate returns the contents of a systemd template for the given template
func readTemplate(component string, templateName string) (out []byte, err error) {
	var templateFile string
	for _, rootPath := range rootPaths {
		filename := path.Join(rootPath, component, templateName)
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

// returns the same template with its placeholders replaced with real values
func processTemplate(template []byte, c *FleetClient) (out []byte){
	if bytes.Index(template, []byte("CHANGEME")) != -1 {
		machines, _ := c.Fleet.Machines()
		randomMachine := machines[rand.Intn(len(machines))]
		template = bytes.Replace(template, []byte("CHANGEME"), []byte(randomMachine.ID), -1)
	}
	return template
}
