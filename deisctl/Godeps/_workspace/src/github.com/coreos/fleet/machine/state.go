package machine

import "github.com/coreos/fleet/resource"

const (
	shortIDLen = 8
)

// MachineState represents a point-in-time snapshot of the
// state of the local host.
type MachineState struct {
	ID             string
	PublicIP       string
	Metadata       map[string]string
	Version        string
	TotalResources resource.ResourceTuple
}

func (s MachineState) ShortID() string {
	if len(s.ID) <= shortIDLen {
		return s.ID
	}
	return s.ID[0:shortIDLen]
}

func (s MachineState) MatchID(ID string) bool {
	return s.ID == ID || s.ShortID() == ID
}

func stackState(top, bottom MachineState) MachineState {
	state := MachineState(bottom)

	if top.PublicIP != "" {
		state.PublicIP = top.PublicIP
	}

	if top.ID != "" {
		state.ID = top.ID
	}

	if top.TotalResources.Cores > 0 {
		state.TotalResources.Cores = top.TotalResources.Cores
	}

	if top.TotalResources.Memory > 0 {
		state.TotalResources.Memory = top.TotalResources.Memory
	}

	if top.TotalResources.Disk > 0 {
		state.TotalResources.Disk = top.TotalResources.Disk
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
