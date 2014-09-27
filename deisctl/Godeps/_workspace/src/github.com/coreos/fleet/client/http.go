package client

import (
	"net/http"
	"net/url"
	"path"

	"github.com/coreos/fleet/Godeps/_workspace/src/code.google.com/p/google-api-go-client/googleapi"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

func NewHTTPClient(c *http.Client, endpoint string) (API, error) {
	svc, err := schema.New(c)
	if err != nil {
		return nil, err
	}

	ep, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// append a slash so the schema.Service knows this is the root path
	ep.Path = path.Join(ep.Path, "v1-alpha") + "/"
	svc.BasePath = ep.String()

	return &HTTPClient{svc: svc}, nil
}

type HTTPClient struct {
	svc *schema.Service

	//NOTE(bcwaldon): This is only necessary until the API interface
	// is fully implemented by HTTPClient
	API
}

func (c *HTTPClient) Machines() ([]machine.MachineState, error) {
	machines := make([]machine.MachineState, 0)
	call := c.svc.Machines.List()
	for call != nil {
		page, err := call.Do()
		if err != nil {
			return nil, err
		}

		machines = append(machines, schema.MapSchemaToMachineStates(page.Machines)...)

		if len(page.NextPageToken) > 0 {
			call = c.svc.Machines.List()
			call.NextPageToken(page.NextPageToken)
		} else {
			call = nil
		}
	}
	return machines, nil
}

func (c *HTTPClient) Units() ([]*schema.Unit, error) {
	var units []*schema.Unit
	call := c.svc.Units.List()
	for call != nil {
		page, err := call.Do()
		if err != nil {
			return nil, err
		}

		units = append(units, page.Units...)

		if len(page.NextPageToken) > 0 {
			call = c.svc.Units.List()
			call.NextPageToken(page.NextPageToken)
		} else {
			call = nil
		}
	}
	return units, nil
}

func (c *HTTPClient) Unit(name string) (*schema.Unit, error) {
	u, err := c.svc.Units.Get(name).Do()
	if err != nil && !is404(err) {
		return nil, err
	}
	return u, nil
}

func (c *HTTPClient) UnitStates() ([]*schema.UnitState, error) {
	var states []*schema.UnitState
	call := c.svc.UnitState.List()
	for call != nil {
		page, err := call.Do()
		if err != nil {
			return nil, err
		}

		states = append(states, page.States...)

		if len(page.NextPageToken) > 0 {
			call = c.svc.UnitState.List()
			call.NextPageToken(page.NextPageToken)
		} else {
			call = nil
		}
	}
	return states, nil
}

func (c *HTTPClient) DestroyUnit(name string) error {
	return c.svc.Units.Delete(name).Do()
}

func (c *HTTPClient) CreateUnit(u *schema.Unit) error {
	return c.svc.Units.Set(u.Name, u).Do()
}

func (c *HTTPClient) SetUnitTargetState(name, target string) error {
	u := schema.Unit{
		Name:         name,
		DesiredState: target,
	}
	return c.svc.Units.Set(name, &u).Do()
}

func is404(err error) bool {
	googerr, ok := err.(*googleapi.Error)
	return ok && googerr.Code == http.StatusNotFound
}
