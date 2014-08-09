package client

import (
	"fmt"
	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/deis/deisctl/utils"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// initialize tabwriter on stdout
func init() {
	out = new(tabwriter.Writer)
	out.Init(os.Stdout, 0, 8, 1, '\t', 0)
}

const (
	defaultListUnitFields = "unit,state,load,active,sub,machine"
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
		"dstate": func(j *job.Job, full bool) string {
			return string(j.TargetState)
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
			if us == nil || us.MachineID == "" {
				return "-"
			}
			ms := cachedMachineState(us.MachineID)
			if ms == nil {
				ms = &machine.MachineState{ID: us.MachineID}
			}
			return machineFullLegend(*ms, full)
		},
		"tmachine": func(j *job.Job, full bool) string {
			if j.TargetMachineID == "" {
				return "-"
			}
			ms := cachedMachineState(j.TargetMachineID)
			if ms == nil {
				ms = &machine.MachineState{ID: j.TargetMachineID}
			}
			return machineFullLegend(*ms, full)
		},
		"hash": func(j *job.Job, full bool) string {
			us := j.UnitState
			if us == nil || us.UnitHash == "" {
				return "-"
			}
			if !full {
				return us.UnitHash[:7]
			}
			return us.UnitHash
		},
	}
)

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

func (c *FleetClient) GetLocaljobs() sort.StringSlice {
	var sortable sort.StringSlice
	jj, err := c.Fleet.Jobs()
	if err != nil {
		return sortable
	}
	for _, j := range jj {
		if strings.HasPrefix(j.Name, "deis-") && j.UnitState.MachineID == utils.GetMachineID("/") {
			sortable = append(sortable, j.Name)
		}
	}
	return sortable
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
