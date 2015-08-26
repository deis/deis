package builds

import (
	"encoding/json"
	"fmt"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List lists an app's builds.
func List(c *client.Client, appID string, results int) ([]api.Build, int, error) {
	u := fmt.Sprintf("/v1/apps/%s/builds/", appID)
	body, count, err := c.LimitedRequest(u, results)

	if err != nil {
		return []api.Build{}, -1, err
	}

	var builds []api.Build
	if err = json.Unmarshal([]byte(body), &builds); err != nil {
		return []api.Build{}, -1, err
	}

	return builds, count, nil
}

// New creates a build for an app.
func New(c *client.Client, appID string, image string,
	procfile map[string]string) (api.Build, error) {

	u := fmt.Sprintf("/v1/apps/%s/builds/", appID)

	req := api.CreateBuildRequest{Image: image, Procfile: procfile}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Build{}, err
	}

	resBody, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return api.Build{}, err
	}

	build := api.Build{}
	if err = json.Unmarshal([]byte(resBody), &build); err != nil {
		return api.Build{}, err
	}

	return build, nil
}
