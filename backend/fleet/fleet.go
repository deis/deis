package fleet

import "github.com/coreos/fleet/client"

// FleetClient used to wrap Fleet API calls
type FleetClient struct {
	Fleet client.API
}

// NewClient returns a client used to communicate with Fleet
// using the Registry API
func NewClient() (*FleetClient, error) {
	client, err := getRegistryClient()
	if err != nil {
		return nil, err
	}
	// set global client
	cAPI = client
	return &FleetClient{Fleet: client}, nil
}

// randomMachineID return a random machineID from the Fleet cluster
func randomMachineID(c *FleetClient) (machineID string, err error) {
	machineState, err := c.Fleet.Machines()
	if err != nil {
		return "", err
	}
	var machineIDs []string
	for _, ms := range machineState {
		machineIDs = append(machineIDs, ms.ID)
	}
	return randomValue(machineIDs), nil
}
