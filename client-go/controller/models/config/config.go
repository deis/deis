package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// Set a config variable.
func Set(c *client.Client, app string, configValues map[string]string) error {
	config := api.ConfigSet{Values: configValues}

	body, err := json.Marshal(config)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("/v1/apps/%s/config", app)

	resBody, status, err := c.BasicRequest("POST", url, body)

	if err != nil {
		return err
	}

	if status != 201 {
		return errors.New(resBody)
	}

	return nil
}
