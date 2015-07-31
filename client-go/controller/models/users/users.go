package users

import (
	"encoding/json"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List users registered with the controller.
func List(c *client.Client) ([]api.User, error) {
	body, err := c.BasicRequest("GET", "/v1/users/", nil)

	if err != nil {
		return []api.User{}, err
	}

	users := api.Users{}
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []api.User{}, err
	}

	return users.Users, nil
}
