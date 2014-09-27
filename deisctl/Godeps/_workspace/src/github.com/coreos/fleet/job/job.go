package job

import (
	"fmt"
	"strings"

	"github.com/coreos/fleet/unit"
)

type JobState string

const (
	JobStateInactive = JobState("inactive")
	JobStateLoaded   = JobState("loaded")
	JobStateLaunched = JobState("launched")
)

// fleet-specific unit file requirement keys.
// "X-" prefix only appears in unit file and is dropped in code before the value is used.
const (
	// Require the unit be scheduled to a specific machine identified by given ID.
	fleetXConditionMachineID = "ConditionMachineID"
	// Legacy form of FleetXConditionMachineID.
	fleetXConditionMachineBootID = "ConditionMachineBootID"
	// Limit eligible machines to the one that hosts a specific unit.
	fleetXConditionMachineOf = "ConditionMachineOf"
	// Prevent a unit from being collocated with other units using glob-matching on the other unit names.
	fleetXConflicts = "Conflicts"
	// Machine metadata key in the unit file, without the X- prefix
	fleetXConditionMachineMetadata = "ConditionMachineMetadata"
	// Machine metadata key for the deprecated `require` flag
	fleetFlagMachineMetadata = "MachineMetadata"
)

func ParseJobState(s string) (JobState, error) {
	js := JobState(s)

	var err error
	if js != JobStateInactive && js != JobStateLoaded && js != JobStateLaunched {
		err = fmt.Errorf("invalid value %q for JobState", s)
		js = JobStateInactive
	}

	return js, err
}

// Job is a legacy construct encapsulating a scheduled unit in fleet
type Job struct {
	Name            string
	State           *JobState
	TargetState     JobState
	TargetMachineID string
	Unit            unit.UnitFile
}

// ScheduledUnit represents a Unit known by fleet and encapsulates its current scheduling state. This does not include Global units.
type ScheduledUnit struct {
	Name            string
	State           *JobState
	TargetMachineID string
}

// Unit represents a Unit that has been submitted to fleet
// (list-unit-files)
type Unit struct {
	Name        string
	Unit        unit.UnitFile
	TargetState JobState
}

// IsGlobal returns whether a Unit is considered a global unit
func (u *Unit) IsGlobal() bool {
	for key, values := range u.Unit.Contents["X-Fleet"] {
		if key == "Global" {
			// Last value found wins
			return strings.ToLower(values[len(values)-1]) == "true"
		}
	}
	return false
}

// NewJob creates a new Job based on the given name and Unit.
// The returned Job has a populated UnitHash and empty JobState.
// nil is returned on failure.
func NewJob(name string, unit unit.UnitFile) *Job {
	return &Job{
		Name:            name,
		State:           nil,
		TargetState:     JobStateInactive,
		TargetMachineID: "",
		Unit:            unit,
	}
}

// Requirements returns all relevant options from the [X-Fleet] section of a unit file.
// Relevant options are identified with a `X-` prefix in the unit.
// This prefix is stripped from relevant options before being returned.
// Furthermore, specifier substitution (using unitPrintf) is performed on all requirements.
func (j *Job) Requirements() map[string][]string {
	uni := unit.NewUnitNameInfo(j.Name)
	requirements := make(map[string][]string)
	for key, values := range j.Unit.Contents["X-Fleet"] {
		if !strings.HasPrefix(key, "X-") {
			continue
		}

		// Strip off leading X-
		key = key[2:]

		if _, ok := requirements[key]; !ok {
			requirements[key] = make([]string, 0)
		}

		if uni != nil {
			for i, v := range values {
				values[i] = unitPrintf(v, *uni)
			}
		}

		requirements[key] = values
	}

	return requirements
}

// Conflicts returns a list of Job names that cannot be scheduled to the same
// machine as this Job.
func (j *Job) Conflicts() []string {
	conflicts, ok := j.Requirements()[fleetXConflicts]
	if ok {
		return conflicts
	}
	return make([]string, 0)
}

// Peers returns a list of Job names that must be scheduled to the same
// machine as this Job.
func (j *Job) Peers() []string {
	peers, ok := j.Requirements()[fleetXConditionMachineOf]
	if !ok {
		return []string{}
	}
	return peers
}

// RequiredTarget determines whether or not this Job must be scheduled to
// a specific machine. If such a requirement exists, the first value returned
// represents the ID of such a machine, while the second value will be a bool
// true. If no requirement exists, an empty string along with a bool false
// will be returned.
func (j *Job) RequiredTarget() (string, bool) {
	requirements := j.Requirements()

	machIDs, ok := requirements[fleetXConditionMachineID]
	if ok && len(machIDs) != 0 {
		return machIDs[0], true
	}

	// Fall back to the legacy option if it exists. This is unlikely
	// to actually work as the user intends, but it's better to
	// prevent a job from starting that has a legacy requirement
	// than to ignore the requirement and let it start.
	bootIDs, ok := requirements[fleetXConditionMachineBootID]
	if ok && len(bootIDs) != 0 {
		return bootIDs[0], true
	}

	return "", false
}

// RequiredTargetMetadata return all machine-related metadata from a Job's requirements
func (j *Job) RequiredTargetMetadata() map[string][]string {
	metadata := make(map[string][]string)
	for key, values := range j.Requirements() {
		// Deprecated syntax added to the metadata via the old `--require` flag.
		if strings.HasPrefix(key, fleetFlagMachineMetadata) {
			if len(values) == 0 {
				continue
			}

			metadata[key[15:]] = values
		} else if key == fleetXConditionMachineMetadata {
			for _, valuePair := range values {
				s := strings.Split(valuePair, "=")

				if len(s) != 2 {
					continue
				}

				if len(s[0]) == 0 || len(s[1]) == 0 {
					continue
				}

				var mValues []string
				if mv, ok := metadata[s[0]]; ok {
					mValues = mv
				}

				metadata[s[0]] = append(mValues, s[1])
			}
		}
	}

	return metadata
}

func (j *Job) Scheduled() bool {
	return len(j.TargetMachineID) > 0
}

// unitPrintf is analogous to systemd's `unit_name_printf`. It will take the
// given string and replace the following specifiers with the values from the
// provided UnitNameInfo:
// 	%n: the full name of the unit               (foo@bar.waldo)
// 	%N: the name of the unit without the suffix (foo@bar)
// 	%p: the prefix                              (foo)
// 	%i: the instance                            (bar)
func unitPrintf(s string, nu unit.UnitNameInfo) (out string) {
	out = strings.Replace(s, "%n", nu.FullName, -1)
	out = strings.Replace(out, "%N", nu.Name, -1)
	out = strings.Replace(out, "%p", nu.Prefix, -1)
	out = strings.Replace(out, "%i", nu.Instance, -1)
	return
}
