package unit

import (
	"github.com/coreos/fleet/pkg"
)

type UnitManager interface {
	Load(string, UnitFile) error
	Unload(string)

	TriggerStart(string)
	TriggerStop(string)

	Units() ([]string, error)
	GetUnitStates(pkg.Set) (map[string]*UnitState, error)
	GetUnitState(string) (*UnitState, error)
}
