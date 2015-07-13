package domains

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List domains registered with an app.
func List(c *client.Client, appID string) ([]api.Domain, error) {
	u := fmt.Sprintf("/v1/apps/%s/domains/", appID)
	body, status, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return []api.Domain{}, err
	}

	if status != 200 {
		return []api.Domain{}, errors.New(body)
	}

	domains := api.Domains{}
	if err = json.Unmarshal([]byte(body), &domains); err != nil {
		return []api.Domain{}, err
	}

	return domains.Domains, nil
}

// New adds a domain to an app.
func New(c *client.Client, appID string, domain string) (api.Domain, error) {
	u := fmt.Sprintf("/v1/apps/%s/domains/", appID)

	req := api.DomainCreateRequest{Domain: domain}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Domain{}, err
	}

	resBody, status, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return api.Domain{}, err
	}

	if status != 201 {
		return api.Domain{}, errors.New(resBody)
	}

	res := api.Domain{}
	if err = json.Unmarshal([]byte(resBody), &res); err != nil {
		return api.Domain{}, err
	}

	return res, nil
}

// Delete removes a domain from an app.
func Delete(c *client.Client, appID string, domain string) error {
	u := fmt.Sprintf("/v1/apps/%s/domains/%s", appID, domain)
	body, status, err := c.BasicRequest("DELETE", u, nil)

	if err != nil {
		return err
	}

	if status != 204 {
		return errors.New(body)
	}

	return nil
}
