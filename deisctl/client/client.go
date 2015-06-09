package client

import (
	"errors"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/cmd"
)

// DeisCtlClient manages Deis components, configuration, and related tasks.
type DeisCtlClient interface {
	Config(argv []string) error
	Install(argv []string) error
	Journal(argv []string) error
	List(argv []string) error
	RefreshUnits(argv []string) error
	Restart(argv []string) error
	Scale(argv []string) error
	SSH(argv []string) error
	Start(argv []string) error
	Status(argv []string) error
	Stop(argv []string) error
	Uninstall(argv []string) error
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
func (c *Client) Config(argv []string) error {
	return cmd.Config(argv)
}

// Install loads the definitions of components from local unit files.
// After Install, the components will be available to Start.
func (c *Client) Install(argv []string) error {
	return cmd.Install(argv, c.Backend)
}

// Journal prints log output for the specified components.
func (c *Client) Journal(argv []string) error {
	return cmd.Journal(argv, c.Backend)
}

// List prints a summary of installed components.
func (c *Client) List(argv []string) error {
	return cmd.ListUnits(argv, c.Backend)
}

// RefreshUnits overwrites local unit files with those requested.
func (c *Client) RefreshUnits(argv []string) error {
	return cmd.RefreshUnits(argv)
}

// Restart stops and then starts components.
func (c *Client) Restart(argv []string) error {
	return cmd.Restart(argv, c.Backend)
}

// Scale grows or shrinks the number of running components.
func (c *Client) Scale(argv []string) error {
	return cmd.Scale(argv, c.Backend)
}

// SSH opens an interactive shell with a machine in the cluster.
func (c *Client) SSH(argv []string) error {
	return cmd.SSH(argv, c.Backend)
}

// Start activates the specified components.
func (c *Client) Start(argv []string) error {
	return cmd.Start(argv, c.Backend)
}

// Status prints the current status of components.
func (c *Client) Status(argv []string) error {
	return cmd.Status(argv, c.Backend)
}

// Stop deactivates the specified components.
func (c *Client) Stop(argv []string) error {
	return cmd.Stop(argv, c.Backend)
}

// Uninstall unloads the definitions of the specified components.
// After Uninstall, the components will be unavailable until Install is called.
func (c *Client) Uninstall(argv []string) error {
	return cmd.Uninstall(argv, c.Backend)
}
