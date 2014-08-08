package deisctl

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
)

const (
	defaultListUnitFields = "unit,state,load,active,sub,desc,machine"
)

type jobToField func(j *job.Job, full bool) string

var (
	out             *tabwriter.Writer
	listUnitsFields = map[string]jobToField{
		"unit": func(j *job.Job, full bool) string {
			return j.Name
		},
		"state": func(j *job.Job, full bool) string {
			js := j.State
			if js != nil {
				return string(*js)
			}
			return "-"
		},
		"load": func(j *job.Job, full bool) string {
			us := j.UnitState
			if us == nil {
				return "-"
			}
			return us.LoadState
		},
		"active": func(j *job.Job, full bool) string {
			us := j.UnitState
			if us == nil {
				return "-"
			}
			return us.ActiveState
		},
		"sub": func(j *job.Job, full bool) string {
			us := j.UnitState
			if us == nil {
				return "-"
			}
			return us.SubState
		},
		"desc": func(j *job.Job, full bool) string {
			d := j.Unit.Description()
			if d == "" {
				return "-"
			}
			return d
		},
		"machine": func(j *job.Job, full bool) string {
			us := j.UnitState
			if us == nil || us.MachineState == nil {
				return "-"
			}
			return machineFullLegend(*us.MachineState, full)
		},
		"hash": func(j *job.Job, full bool) string {
			if !full {
				return j.UnitHash.Short()
			}
			return j.UnitHash.String()
		},
	}
)

// initialize tabwriter on stdout
func init() {
	out = new(tabwriter.Writer)
	out.Init(os.Stdout, 0, 8, 1, '\t', 0)
}

// printList writes units to stdout using a tabwriter
func printList(jobs map[string]job.Job, sortable sort.StringSlice) {
	cols := strings.Split(defaultListUnitFields, ",")
	for _, s := range cols {
		if _, ok := listUnitsFields[s]; !ok {
			fmt.Fprintf(os.Stderr, "Invalid key in output format: %q\n", s)
		}
	}
	fmt.Fprintln(out, strings.ToUpper(strings.Join(cols, "\t")))
	for _, name := range sortable {
		var f []string
		j := jobs[name]
		for _, c := range cols {
			f = append(f, listUnitsFields[c](&j, false))
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
