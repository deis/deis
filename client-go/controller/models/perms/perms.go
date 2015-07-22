package perms

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List users that can access an app.
func List(c *client.Client, appID string) ([]string, error) {
	body, err := doList(c, fmt.Sprintf("/v1/apps/%s/perms/", appID))

	if err != nil {
		return []string{}, err
	}

	users := api.PermsAppResponse{}
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []string{}, err
	}

	return users.Users, nil
}

// ListAdmins lists administrators.
func ListAdmins(c *client.Client) ([]string, error) {
	body, err := doList(c, "/v1/admin/perms/")

	if err != nil {
		return []string{}, err
	}

	users := api.PermsAdminResponse{}
	if err = json.Unmarshal([]byte(body), &users); err != nil {
		return []string{}, err
	}

	usersList := []string{}

	for _, user := range users.Users {
		usersList = append(usersList, user.Username)
	}

	return usersList, nil
}

func doList(c *client.Client, u string) (string, error) {
	body, status, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return "", err
	}

	if status != 200 {
		return "", errors.New(body)
	}

	return body, nil
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

	body, status, err := c.BasicRequest("POST", u, reqBody)

	if err != nil {
		return err
	}

	if status != 201 {
		return errors.New(body)
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
	body, status, err := c.BasicRequest("DELETE", u, nil)

	if err != nil {
		return err
	}

	if status != 204 {
		return errors.New(body)
	}

	return nil
}
