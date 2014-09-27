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

// getMachineIDs returns a list of machineIDs from the Fleet cluster
func getMachineIDs(c *FleetClient) (machineIDs []string, err error) {
	machineState, err := c.Fleet.Machines()
	if err != nil {
		return []string{}, err
	}
	for _, ms := range machineState {
		machineIDs = append(machineIDs, ms.ID)
	}
	return
}

// randomMachineID return a random machineID from the Fleet cluster
func randomMachineID(c *FleetClient) (machineID string, err error) {
	var machineIDs []string
	machineIDs, err = getMachineIDs(c)
	if err != nil {
		return "", err
	}
	return randomValue(machineIDs), nil
}
