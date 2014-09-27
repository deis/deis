package client

import (
	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

type API interface {
	Machines() ([]machine.MachineState, error)

	Unit(string) (*schema.Unit, error)
	Units() ([]*schema.Unit, error)
	UnitStates() ([]*schema.UnitState, error)

	SetUnitTargetState(name, target string) error
	CreateUnit(*schema.Unit) error
	DestroyUnit(string) error
}
