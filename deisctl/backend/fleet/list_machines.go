package fleet

import (
	"fmt"
	"sort"
	"strings"

	"github.com/coreos/fleet/machine"
)

const (
	defaultListMachinesFields = "machine,ip,metadata"
)

var (
	listMachinesFields = map[string]machineToField{
		"machine": func(ms *machine.MachineState, full bool) string {
			return machineIDLegend(*ms, full)
		},
		"ip": func(ms *machine.MachineState, full bool) string {
			if len(ms.PublicIP) == 0 {
				return "-"
			}
			return ms.PublicIP
		},
		"metadata": func(ms *machine.MachineState, full bool) string {
			if len(ms.Metadata) == 0 {
				return "-"
			}
			return formatMetadata(ms.Metadata)
		},
	}
)

type machineToField func(ms *machine.MachineState, full bool) string

// ListMachines prints all nodes to Stdout
func (c *FleetClient) ListMachines() (err error) {
	machines, err := c.Fleet.Machines()
	if err != nil {
		return err
	}

	c.printMachines(machines)
	return
}

// printUnits writes units to stdout using a tabwriter
func (c *FleetClient) printMachines(states []machine.MachineState) {
	cols := strings.Split(defaultListMachinesFields, ",")
	fmt.Fprintln(c.out, strings.ToUpper(strings.Join(cols, "\t")))
	for _, ms := range states {
		var f []string
		for _, c := range cols {
			f = append(f, listMachinesFields[c](&ms, false))
		}
		fmt.Fprintln(c.out, strings.Join(f, "\t"))
	}
	c.out.Flush()
}

func formatMetadata(metadata map[string]string) string {
	pairs := make([]string, len(metadata))
	idx := 0
	var sorted sort.StringSlice
	for k := range metadata {
		sorted = append(sorted, k)
	}
	sorted.Sort()
	for _, key := range sorted {
		value := metadata[key]
		pairs[idx] = fmt.Sprintf("%s=%s", key, value)
		idx++
	}
	return strings.Join(pairs, ",")
}
