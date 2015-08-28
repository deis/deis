package perms

import (
	"encoding/json"
	"fmt"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List users that can access an app.
func List(c *client.Client, appID string) ([]string, error) {
	body, err := c.BasicRequest("GET", fmt.Sprintf("/v1/apps/%s/perms/", appID), nil)

	if err != nil {
		return []string{}, err
	}

	var users api.PermsAppResponse
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []string{}, err
	}

	return users.Users, nil
}

// ListAdmins lists administrators.
func ListAdmins(c *client.Client, results int) ([]string, int, error) {
	body, count, err := c.LimitedRequest("/v1/admin/perms/", results)

	if err != nil {
		return []string{}, -1, err
	}

	var users []api.PermsRequest
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []string{}, -1, err
	}

	usersList := []string{}

	for _, user := range users {
		usersList = append(usersList, user.Username)
	}

	return usersList, count, nil
}

// New adds a user to an app.
func New(c *client.Client, appID string, username string) error {
	return doNew(c, fmt.Sprintf("/v1/apps/%s/perms/", appID), username)
}

// NewAdmin makes a user an administrator.
func NewAdmin(c *client.Client, username string) error {
	return doNew(c, "/v1/admin/perms/", username)
}

func doNew(c *client.Client, u string, username string) error {
	req := api.PermsRequest{Username: username}

	reqBody, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", u, reqBody)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes a user from an app.
func Delete(c *client.Client, appID string, username string) error {
	return doDelete(c, fmt.Sprintf("/v1/apps/%s/perms/%s", appID, username))
}

// DeleteAdmin removes administrative privilages from a user.
func DeleteAdmin(c *client.Client, username string) error {
	return doDelete(c, fmt.Sprintf("/v1/admin/perms/%s", username))
}

func doDelete(c *client.Client, u string) error {
	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}
