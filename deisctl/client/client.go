package client

import (
	"errors"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/cmd"
)

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

type Client struct {
	Backend backend.Backend
}

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

func (c *Client) Config() error {
	return cmd.Config()
}

func (c *Client) Install(targets []string) error {
	return cmd.Install(c.Backend, targets)
}

func (c *Client) Journal(targets []string) error {
	return cmd.Journal(c.Backend, targets)
}

func (c *Client) List() error {
	return cmd.ListUnits(c.Backend)
}

func (c *Client) RefreshUnits() error {
	return cmd.RefreshUnits()
}

func (c *Client) Restart(targets []string) error {
	return cmd.Restart(c.Backend, targets)
}

func (c *Client) Scale(targets []string) error {
	return cmd.Scale(c.Backend, targets)
}

func (c *Client) Start(targets []string) error {
	return cmd.Start(c.Backend, targets)
}

func (c *Client) Status(targets []string) error {
	return cmd.Status(c.Backend, targets)
}

func (c *Client) Stop(targets []string) error {
	return cmd.Stop(c.Backend, targets)
}

func (c *Client) Uninstall(targets []string) error {
	return cmd.Uninstall(c.Backend, targets)
}

func (c *Client) Update() error {
	return cmd.Update()
}
