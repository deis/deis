package client

import (
	"errors"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/cmd"
)

// DeisCtlClient manages Deis components, configuration, and related tasks.
type DeisCtlClient interface {
	Config() error
	Install(targets []string) error
	Journal(targets []string) error
	List() error
	RefreshUnits() error
	Restart(targets []string) error
	Scale(targets []string) error
	Start(targets []string) error
	Status(targets []string) error
	Stop(targets []string) error
	Uninstall(targets []string) error
	Update() error
}

// Client uses a backend to implement the DeisCtlClient interface.
type Client struct {
	Backend backend.Backend
}

// NewClient returns a Client using the requested backend.
// The only backend currently supported is "fleet".
func NewClient(requestedBackend string) (*Client, error) {
	var backend backend.Backend

	if requestedBackend == "" {
		requestedBackend = "fleet"
	}

	switch requestedBackend {
	case "fleet":
		b, err := fleet.NewClient()
		if err != nil {
			return nil, err
		}
		backend = b
	default:
		return nil, errors.New("invalid backend")
	}
	return &Client{Backend: backend}, nil
}

// Config gets or sets a configuration value from the cluster.
//
// A configuration value is stored and retrieved from a key/value store (in this case, etcd)
// at /deis/<component>/<config>. Configuration values are typically used for component-level
// configuration, such as enabling TLS for the routers.
func (c *Client) Config() error {
	return cmd.Config()
}

// Install loads components' definitions from local unit files.
func (c *Client) Install(targets []string) error {
	return cmd.Install(c.Backend, targets)
}

// Journal prints log output for the specified components.
func (c *Client) Journal(targets []string) error {
	return cmd.Journal(c.Backend, targets)
}

// List prints a summary of installed components.
func (c *Client) List() error {
	return cmd.ListUnits(c.Backend)
}

// RefreshUnits overwrites local unit files with those requested.
func (c *Client) RefreshUnits() error {
	return cmd.RefreshUnits()
}

// Restart stops and then starts components.
func (c *Client) Restart(targets []string) error {
	return cmd.Restart(c.Backend, targets)
}

// Scale grows or shrinks the number of running components.
func (c *Client) Scale(targets []string) error {
	return cmd.Scale(c.Backend, targets)
}

// Start activates the specified components.
func (c *Client) Start(targets []string) error {
	return cmd.Start(c.Backend, targets)
}

// Status prints the current state of components.
func (c *Client) Status(targets []string) error {
	return cmd.Status(c.Backend, targets)
}

// Stop deactivates the specified components.
func (c *Client) Stop(targets []string) error {
	return cmd.Stop(c.Backend, targets)
}

// Uninstall unloads components' definitions.
func (c *Client) Uninstall(targets []string) error {
	return cmd.Uninstall(c.Backend, targets)
}

// Update changes the platform version on a cluster host.
func (c *Client) Update() error {
	return cmd.Update()
}
