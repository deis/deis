package machine

const (
	shortIDLen = 8
)

// MachineState represents a point-in-time snapshot of the
// state of the local host.
type MachineState struct {
	ID       string
	PublicIP string
	Metadata map[string]string
	Version  string
}

func (ms MachineState) ShortID() string {
	if len(ms.ID) <= shortIDLen {
		return ms.ID
	}
	return ms.ID[0:shortIDLen]
}

func (ms MachineState) MatchID(ID string) bool {
	return ms.ID == ID || ms.ShortID() == ID
}

// stackState is used to merge two MachineStates. Values configured on the top
// MachineState always take precedence over those on the bottom.
func stackState(top, bottom MachineState) MachineState {
	state := MachineState(bottom)

	if top.PublicIP != "" {
		state.PublicIP = top.PublicIP
	}

	if top.ID != "" {
		state.ID = top.ID
	}

	//FIXME: This will *always* overwrite the bottom's metadata,
	// but the only use-case we have today does not ever have
	// metadata on the bottom.
	if len(top.Metadata) > 0 {
		state.Metadata = top.Metadata
	}

	if top.Version != "" {
		state.Version = top.Version
	}

	return state
}
