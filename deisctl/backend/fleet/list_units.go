package fleet

import (
	"fmt"
	"strings"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
	"github.com/deis/deis/deisctl/units"
)

const (
	//defaultListUnitFields  = "unit,state,load,active,sub,machine"
	defaultListUnitsFields = "unit,machine,load,active,sub"
)

type usToField func(c *FleetClient, us *schema.UnitState, full bool) string

var (
	listUnitsFields = map[string]usToField{
		"unit": func(c *FleetClient, us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.Name
		},
		"load": func(c *FleetClient, us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.SystemdLoadState
		},
		"active": func(c *FleetClient, us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.SystemdActiveState
		},
		"sub": func(c *FleetClient, us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.SystemdSubState
		},
		"machine": func(c *FleetClient, us *schema.UnitState, full bool) string {
			if us == nil || us.MachineID == "" {
				return "-"
			}
			ms := c.cachedMachineState(us.MachineID)
			if ms == nil {
				ms = &machine.MachineState{ID: us.MachineID}
			}
			return machineFullLegend(*ms, full)
		},
		"hash": func(c *FleetClient, us *schema.UnitState, full bool) string {
			if us == nil || us.Hash == "" {
				return "-"
			}
			if !full {
				return us.Hash[:7]
			}
			return us.Hash
		},
	}
)

// ListUnits prints all Deis-related units to Stdout
func (c *FleetClient) ListUnits() (err error) {
	var states []*schema.UnitState

	unitStates, err := c.Fleet.UnitStates()
	if err != nil {
		return err
	}

	for _, us := range unitStates {
		for _, prefix := range units.Names {
			if strings.HasPrefix(us.Name, prefix) {
				states = append(states, us)
				break
			}
		}
	}
	c.printUnits(states)
	return
}

// printUnits writes units to stdout using a tabwriter
func (c *FleetClient) printUnits(states []*schema.UnitState) {
	cols := strings.Split(defaultListUnitsFields, ",")
	fmt.Fprintln(c.out, strings.ToUpper(strings.Join(cols, "\t")))
	for _, us := range states {
		var f []string
		for _, col := range cols {
			f = append(f, listUnitsFields[col](c, us, false))
		}
		fmt.Fprintln(c.out, strings.Join(f, "\t"))
	}
	c.out.Flush()
}

func machineIDLegend(ms machine.MachineState, full bool) string {
	legend := ms.ID
	if !full {
		legend = fmt.Sprintf("%s...", ms.ShortID())
	}
	return legend
}

func machineFullLegend(ms machine.MachineState, full bool) string {
	legend := machineIDLegend(ms, full)
	if len(ms.PublicIP) > 0 {
		legend = fmt.Sprintf("%s/%s", legend, ms.PublicIP)
	}
	return legend
}
