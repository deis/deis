package unit

type UnitManager interface {
	Load(string, Unit) error
	Unload(string)

	Start(string)
	Stop(string)

	Units() ([]string, error)
	GetUnitState(string) (*UnitState, error)

	MarshalJSON() ([]byte, error)
}
