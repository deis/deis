package domains

import (
	"encoding/json"
	"fmt"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List domains registered with an app.
func List(c *client.Client, appID string, results int) ([]api.Domain, int, error) {
	u := fmt.Sprintf("/v1/apps/%s/domains/", appID)
	body, count, err := c.LimitedRequest(u, results)

	if err != nil {
		return []api.Domain{}, -1, err
	}

	var domains []api.Domain
	if err = json.Unmarshal([]byte(body), &domains); err != nil {
		return []api.Domain{}, -1, err
	}

	return domains, count, nil
}

// New adds a domain to an app.
func New(c *client.Client, appID string, domain string) (api.Domain, error) {
	u := fmt.Sprintf("/v1/apps/%s/domains/", appID)

	req := api.DomainCreateRequest{Domain: domain}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Domain{}, err
	}

	resBody, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return api.Domain{}, err
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
	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}
