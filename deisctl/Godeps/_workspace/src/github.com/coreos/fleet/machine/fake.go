package machine

type FakeMachine struct {
	MachineState MachineState
}

func (fm *FakeMachine) State() MachineState {
	return fm.MachineState
}
