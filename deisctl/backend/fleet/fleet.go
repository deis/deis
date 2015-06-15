package fleet

import (
	"github.com/coreos/fleet/client"
	"github.com/coreos/fleet/machine"
)

// FleetClient used to wrap Fleet API calls
type FleetClient struct {
	Fleet client.API

	// used to cache MachineStates
	machineStates map[string]*machine.MachineState
}

// NewClient returns a client used to communicate with Fleet
// using the Registry API
func NewClient() (*FleetClient, error) {
	client, err := getRegistryClient()
	if err != nil {
		return nil, err
	}
	return &FleetClient{Fleet: client}, nil
}
