package releases

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List lists an app's releases.
func List(c *client.Client, appID string) ([]api.Release, error) {
	u := fmt.Sprintf("/v1/apps/%s/releases/", appID)

	body, status, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return []api.Release{}, err
	}

	if status != 200 {
		return []api.Release{}, errors.New(body)
	}

	releases := api.Releases{}
	if err = json.Unmarshal([]byte(body), &releases); err != nil {
		return []api.Release{}, err
	}

	return releases.Releases, nil
}

// Get a release of an app.
func Get(c *client.Client, appID string, version int) (api.Release, error) {
	u := fmt.Sprintf("/v1/apps/%s/releases/v%d/", appID, version)

	body, status, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return api.Release{}, err
	}

	if status != 200 {
		return api.Release{}, errors.New(body)
	}

	release := api.Release{}
	if err = json.Unmarshal([]byte(body), &release); err != nil {
		return api.Release{}, err
	}

	return release, nil
}

// Rollback rolls back an app to a previous release.
func Rollback(c *client.Client, appID string, version int) (int, error) {
	u := fmt.Sprintf("/v1/apps/%s/releases/rollback/", appID)

	req := api.ReleaseRollback{Version: version}

	var err error
	var reqBody []byte
	if version != -1 {
		reqBody, err = json.Marshal(req)

		if err != nil {
			return -1, err
		}
	}

	body, status, err := c.BasicRequest("POST", u, reqBody)

	if err != nil {
		return -1, err
	}

	if status != 201 {
		return -1, errors.New(body)
	}

	response := api.ReleaseRollback{}

	if err = json.Unmarshal([]byte(body), &response); err != nil {
		return -1, err
	}

	return response.Version, nil
}
