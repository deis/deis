package fleet

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
	"github.com/deis/deis/deisctl/utils"
)

// initialize tabwriter on stdout
func init() {
	out = new(tabwriter.Writer)
	out.Init(os.Stdout, 0, 8, 1, '\t', 0)
}

const (
	//defaultListUnitFields  = "unit,state,load,active,sub,machine"
	defaultListUnitsFields = "unit,machine,load,active,sub"
)

type usToField func(us *schema.UnitState, full bool) string

var (
	out             *tabwriter.Writer
	listUnitsFields = map[string]usToField{
		"unit": func(us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.Name
		},
		"load": func(us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.SystemdLoadState
		},
		"active": func(us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.SystemdActiveState
		},
		"sub": func(us *schema.UnitState, full bool) string {
			if us == nil {
				return "-"
			}
			return us.SystemdSubState
		},
		"machine": func(us *schema.UnitState, full bool) string {
			if us == nil || us.MachineID == "" {
				return "-"
			}
			ms := cachedMachineState(us.MachineID)
			if ms == nil {
				ms = &machine.MachineState{ID: us.MachineID}
			}
			return machineFullLegend(*ms, full)
		},
		"hash": func(us *schema.UnitState, full bool) string {
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
	var sortable sort.StringSlice
	states := make(map[string]*schema.UnitState, 0)

	unitStates, err := cAPI.UnitStates()
	if err != nil {
		return err
	}

	for _, us := range unitStates {
		if strings.HasPrefix(us.Name, "deis-") {
			states[us.Name] = us
			sortable = append(sortable, us.Name)
		}
	}
	sortable.Sort()
	printUnits(states, sortable)
	return
}

func (c *FleetClient) GetLocaljobs() sort.StringSlice {
	var sortable sort.StringSlice
	unitStates, err := c.Fleet.UnitStates()
	if err != nil {
		return sortable
	}
	for _, us := range unitStates {
		if strings.HasPrefix(us.Name, "deis-") && us.MachineID == utils.GetMachineID("/") {
			sortable = append(sortable, us.Name)
		}
	}
	return sortable
}

// printUnits writes units to stdout using a tabwriter
func printUnits(states map[string]*schema.UnitState, sortable sort.StringSlice) {
	cols := strings.Split(defaultListUnitsFields, ",")
	for _, s := range cols {
		if _, ok := listUnitsFields[s]; !ok {
			fmt.Fprintf(os.Stderr, "Invalid key in output format: %q\n", s)
		}
	}
	fmt.Fprintln(out, strings.ToUpper(strings.Join(cols, "\t")))
	for _, name := range sortable {
		var f []string
		us := states[name]
		for _, c := range cols {
			f = append(f, listUnitsFields[c](us, false))
		}
		fmt.Fprintln(out, strings.Join(f, "\t"))
	}
	out.Flush()
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
