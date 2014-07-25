package job

import (
	"strconv"
	"strings"

	log "github.com/coreos/fleet/third_party/github.com/golang/glog"

	"github.com/coreos/fleet/resource"
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
	// Memory required in MB
	fleetXMemoryReservation = "MemoryReservation"
	// Cores required in hundreds, ie 100=1core, 50=0.5core, 200=2cores
	fleetXCoresReservation = "CoresReservation"
	// Disk required in MB
	fleetXDiskReservation = "DiskReservation"
	// Machine metadata key for the deprecated `require` flag
	fleetFlagMachineMetadata = "MachineMetadata"
)

func ParseJobState(s string) *JobState {
	js := JobState(s)
	if js != JobStateInactive && js != JobStateLoaded && js != JobStateLaunched {
		return nil
	}
	return &js
}

type Job struct {
	Name      string
	State     *JobState
	Unit      unit.Unit
	UnitHash  unit.Hash
	UnitState *unit.UnitState
}

// NewJob creates a new Job based on the given name and Unit.
// The returned Job has a populated UnitHash and empty JobState and
// UnitState. nil is returned on failure.
func NewJob(name string, unit unit.Unit) *Job {
	return &Job{
		Name:      name,
		State:     nil,
		Unit:      unit,
		UnitHash:  unit.Hash(),
		UnitState: nil,
	}
}

// IsOneshot identifies whether the Job is intended to execute once, and, on
// completion, not be migrated around the cluster. This is determined by
// whether the unit file associated with the Job is a Service of type "oneshot".
func (j *Job) IsOneshot() bool {
	s, ok := j.Unit.Contents["Service"]
	if !ok {
		return false
	}
	t, ok := s["Type"]
	if !ok || len(t) == 0 {
		return false
	}
	// If multiple Types are defined, systemd uses the last
	return t[len(t)-1] == "oneshot"
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

func (j *Job) resourceFromKey(resKey string) int {
	valStr, ok := j.Requirements()[resKey]
	if ok && len(valStr) > 0 {
		val, err := strconv.Atoi(valStr[0])
		if err != nil {
			log.Errorf("failed to parse resource requirement %s from %s: %v", resKey, j.Name, err)
			return 0
		}
		return val
	}
	return 0
}

func (j *Job) Resources() resource.ResourceTuple {
	return resource.ResourceTuple{
		Cores:  j.resourceFromKey(fleetXCoresReservation),
		Memory: j.resourceFromKey(fleetXMemoryReservation),
		Disk:   j.resourceFromKey(fleetXDiskReservation),
	}
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
