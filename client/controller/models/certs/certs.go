package certs

import (
	"encoding/json"
	"fmt"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List certs registered with the controller.
func List(c *client.Client, results int) ([]api.Cert, int, error) {
	body, count, err := c.LimitedRequest("/v1/certs/", results)

	if err != nil {
		return []api.Cert{}, -1, err
	}

	var res []api.Cert
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return []api.Cert{}, -1, err
	}

	return res, count, nil
}

// New creates a new cert.
func New(c *client.Client, cert string, key string, commonName string) (api.Cert, error) {
	req := api.CertCreateRequest{Certificate: cert, Key: key, Name: commonName}
	reqBody, err := json.Marshal(req)

	if err != nil {
		return api.Cert{}, err
	}

	resBody, err := c.BasicRequest("POST", "/v1/certs/", reqBody)

	if err != nil {
		return api.Cert{}, err
	}

	resCert := api.Cert{}
	if err = json.Unmarshal([]byte(resBody), &resCert); err != nil {
		return api.Cert{}, err
	}

	return resCert, nil
}

// Delete removes a cert.
func Delete(c *client.Client, commonName string) error {
	u := fmt.Sprintf("/v1/certs/%s", commonName)

	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}
