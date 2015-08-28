package users

import (
	"encoding/json"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List users registered with the controller.
func List(c *client.Client, results int) ([]api.User, int, error) {
	body, count, err := c.LimitedRequest("/v1/users/", results)

	if err != nil {
		return []api.User{}, -1, err
	}

	var users []api.User
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []api.User{}, -1, err
	}

	return users, count, nil
}
