package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List lists an app's config.
func List(c *client.Client, app string) (api.Config, error) {
	u := fmt.Sprintf("/v1/apps/%s/config/", app)

	body, status, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return api.Config{}, err
	}

	if status != 200 {
		return api.Config{}, errors.New(body)
	}

	config := api.Config{}
	if err = json.Unmarshal([]byte(body), &config); err != nil {
		return api.Config{}, err
	}

	return config, nil
}

// Set sets an app's config variables.
func Set(c *client.Client, app string, config api.Config) (api.Config, error) {
	body, err := json.Marshal(config)

	if err != nil {
		return api.Config{}, err
	}

	u := fmt.Sprintf("/v1/apps/%s/config/", app)

	resBody, status, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return api.Config{}, err
	}

	if status != 201 {
		return api.Config{}, errors.New(resBody)
	}

	newConfig := api.Config{}
	if err = json.Unmarshal([]byte(resBody), &newConfig); err != nil {
		return api.Config{}, err
	}

	return newConfig, nil
}
