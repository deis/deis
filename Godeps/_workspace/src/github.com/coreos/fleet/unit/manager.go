package unit

import (
	"github.com/coreos/fleet/pkg"
)

type UnitManager interface {
	Load(string, UnitFile) error
	Unload(string)

	Start(string)
	Stop(string)

	Units() ([]string, error)
	GetUnitStates(pkg.Set) (map[string]*UnitState, error)
	GetUnitState(string) (*UnitState, error)
}
