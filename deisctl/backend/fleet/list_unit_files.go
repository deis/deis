package fleet

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

const (
	defaultListUnitFilesFields = "unit,hash,dstate,state,tmachine"
)

var (
	listUnitFilesFields = map[string]unitToField{
		"unit": func(c *FleetClient, u schema.Unit, full bool) string {
			return u.Name
		},
		"global": func(c *FleetClient, u schema.Unit, full bool) string {
			return strconv.FormatBool(suToGlobal(u))
		},
		"dstate": func(c *FleetClient, u schema.Unit, full bool) string {
			if u.DesiredState == "" {
				return "-"
			}
			return u.DesiredState
		},
		"tmachine": func(c *FleetClient, u schema.Unit, full bool) string {
			if suToGlobal(u) || u.MachineID == "" {
				return "-"
			}
			ms := c.cachedMachineState(u.MachineID)
			if ms == nil {
				ms = &machine.MachineState{ID: u.MachineID}
			}

			return machineFullLegend(*ms, full)
		},
		"state": func(c *FleetClient, u schema.Unit, full bool) string {
			if suToGlobal(u) || u.CurrentState == "" {
				return "-"
			}
			return u.CurrentState
		},
		"hash": func(c *FleetClient, u schema.Unit, full bool) string {
			uf := schema.MapSchemaUnitOptionsToUnitFile(u.Options)
			if !full {
				return uf.Hash().Short()
			}
			return uf.Hash().String()
		},
		"desc": func(c *FleetClient, u schema.Unit, full bool) string {
			uf := schema.MapSchemaUnitOptionsToUnitFile(u.Options)
			d := uf.Description()
			if d == "" {
				return "-"
			}
			return d
		},
	}
)

type unitToField func(c *FleetClient, u schema.Unit, full bool) string

// ListUnitFiles prints all Deis-related unit files to Stdout
func (c *FleetClient) ListUnitFiles() (err error) {
	var sortable sort.StringSlice
	units := make(map[string]*schema.Unit, 0)

	us, err := c.Fleet.Units()
	if err != nil {
		return err
	}

	for _, u := range us {
		if strings.HasPrefix(u.Name, "deis-") {
			units[u.Name] = u
			sortable = append(sortable, u.Name)
		}
	}
	sortable.Sort()
	c.printUnitFiles(units, sortable)
	return
}

// printUnitFiles writes unit files to stdout using a tabwriter
func (c *FleetClient) printUnitFiles(units map[string]*schema.Unit, sortable sort.StringSlice) {
	cols := strings.Split(defaultListUnitFilesFields, ",")
	fmt.Fprintln(c.out, strings.ToUpper(strings.Join(cols, "\t")))
	for _, name := range sortable {
		var f []string
		u := units[name]
		for _, col := range cols {
			f = append(f, listUnitFilesFields[col](c, *u, false))
		}
		fmt.Fprintln(c.out, strings.Join(f, "\t"))
	}
	c.out.Flush()
}

// suToGlobal returns whether or not a schema.Unit refers to a global unit
func suToGlobal(su schema.Unit) bool {
	u := job.Unit{
		Unit: *schema.MapSchemaUnitOptionsToUnitFile(su.Options),
	}
	return u.IsGlobal()
}
