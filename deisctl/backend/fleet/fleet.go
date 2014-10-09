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
