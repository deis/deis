package schema

import (
	gsunit "github.com/coreos/fleet/Godeps/_workspace/src/github.com/coreos/go-systemd/unit"

	"github.com/coreos/fleet/job"
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/unit"
)

func MapUnitFileToSchemaUnitOptions(uf *unit.UnitFile) []*UnitOption {
	sopts := make([]*UnitOption, len(uf.Options))
	for i, opt := range uf.Options {
		sopts[i] = &UnitOption{
			Section: opt.Section,
			Name:    opt.Name,
			Value:   opt.Value,
		}
	}
	return sopts
}

func MapSchemaUnitOptionsToUnitFile(sopts []*UnitOption) *unit.UnitFile {
	opts := make([]*gsunit.UnitOption, len(sopts))
	for i, sopt := range sopts {
		opts[i] = &gsunit.UnitOption{
			Section: sopt.Section,
			Name:    sopt.Name,
			Value:   sopt.Value,
		}
	}

	return unit.NewUnitFromOptions(opts)
}

func MapSchemaUnitToUnit(entity *Unit) *job.Unit {
	uf := MapSchemaUnitOptionsToUnitFile(entity.Options)
	j := job.Unit{
		Name: entity.Name,
		Unit: *uf,
	}
	return &j
}

func MapUnitToSchemaUnit(u *job.Unit, su *job.ScheduledUnit) *Unit {
	s := Unit{
		Name:         u.Name,
		Options:      MapUnitFileToSchemaUnitOptions(&(u.Unit)),
		DesiredState: string(u.TargetState),
	}

	if su != nil {
		s.MachineID = su.TargetMachineID
		if su.State != nil {
			s.CurrentState = string(*su.State)
		}
	}

	return &s
}

func MapSchemaUnitsToUnits(entities []*Unit) []job.Unit {
	units := make([]job.Unit, len(entities))
	for i, _ := range entities {
		entity := entities[i]
		units[i] = *(MapSchemaUnitToUnit(entity))
	}

	return units
}

func MapMachineStateToSchema(ms *machine.MachineState) *Machine {
	sm := Machine{
		Id:        ms.ID,
		PrimaryIP: ms.PublicIP,
	}

	sm.Metadata = make(map[string]string, len(ms.Metadata))
	for k, v := range ms.Metadata {
		sm.Metadata[k] = v
	}

	return &sm
}

func MapSchemaToMachineStates(entities []*Machine) []machine.MachineState {
	machines := make([]machine.MachineState, len(entities))
	for i, _ := range entities {
		me := entities[i]

		ms := machine.MachineState{
			ID:       me.Id,
			PublicIP: me.PrimaryIP,
		}

		ms.Metadata = make(map[string]string, len(me.Metadata))
		for k, v := range me.Metadata {
			ms.Metadata[k] = v
		}

		machines[i] = ms
	}

	return machines
}

func MapUnitStatesToSchemaUnitStates(entities []*unit.UnitState) []*UnitState {
	sus := make([]*UnitState, len(entities))
	for i, e := range entities {
		sus[i] = MapUnitStateToSchemaUnitState(e)
	}

	return sus
}

func MapUnitStateToSchemaUnitState(entity *unit.UnitState) *UnitState {
	us := UnitState{
		Name:               entity.UnitName,
		Hash:               entity.UnitHash,
		MachineID:          entity.MachineID,
		SystemdLoadState:   entity.LoadState,
		SystemdActiveState: entity.ActiveState,
		SystemdSubState:    entity.SubState,
	}

	return &us
}

func MapSchemaUnitStatesToUnitStates(entities []*UnitState) []*unit.UnitState {
	us := make([]*unit.UnitState, len(entities))
	for i, e := range entities {
		us[i] = &unit.UnitState{
			UnitName:    e.Name,
			UnitHash:    e.Hash,
			MachineID:   e.MachineID,
			LoadState:   e.SystemdLoadState,
			ActiveState: e.SystemdActiveState,
			SubState:    e.SystemdSubState,
		}
	}

	return us
}

func MapSchemaUnitToScheduledUnit(entity *Unit) *job.ScheduledUnit {
	cs := job.JobState(entity.CurrentState)
	return &job.ScheduledUnit{
		Name:            entity.Name,
		State:           &cs,
		TargetMachineID: entity.MachineID,
	}
}

func MapSchemaUnitsToScheduledUnits(entities []*Unit) []job.ScheduledUnit {
	su := make([]job.ScheduledUnit, len(entities))
	for i, e := range entities {
		su[i] = *MapSchemaUnitToScheduledUnit(e)
	}

	return su
}
