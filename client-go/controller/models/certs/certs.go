package certs

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List certs registered with the controller.
func List(c *client.Client) ([]api.Cert, error) {
	body, status, err := c.BasicRequest("GET", "/v1/certs/", nil)

	if err != nil {
		return []api.Cert{}, err
	}

	if status != 200 {
		return []api.Cert{}, errors.New(body)
	}

	res := api.Certs{}
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return []api.Cert{}, err
	}

	return res.Certs, nil
}

// New creates a new cert.
func New(c *client.Client, cert string, key string, commonName string) (api.Cert, error) {
	req := api.CertCreateRequest{Certificate: cert, Key: key, Name: commonName}
	reqBody, err := json.Marshal(req)

	if err != nil {
		return api.Cert{}, err
	}

	resBody, status, err := c.BasicRequest("POST", "/v1/certs/", reqBody)

	if err != nil {
		return api.Cert{}, err
	}

	if status != 201 {
		return api.Cert{}, errors.New(resBody)
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

	resBody, status, err := c.BasicRequest("DELETE", u, nil)

	if err != nil {
		return err
	}

	if status != 204 {
		return errors.New(resBody)
	}

	return nil
}
