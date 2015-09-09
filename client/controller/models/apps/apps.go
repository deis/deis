package apps

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List lists apps on a Deis controller.
func List(c *client.Client, results int) ([]api.App, int, error) {
	body, count, err := c.LimitedRequest("/v1/apps/", results)

	if err != nil {
		return []api.App{}, -1, err
	}

	var apps []api.App
	if err = json.Unmarshal([]byte(body), &apps); err != nil {
		return []api.App{}, -1, err
	}

	return apps, count, nil
}

// New creates a new app.
func New(c *client.Client, id string) (api.App, error) {
	body := []byte{}

	var err error
	if id != "" {
		req := api.AppCreateRequest{ID: id}
		body, err = json.Marshal(req)

		if err != nil {
			return api.App{}, err
		}
	}

	resBody, err := c.BasicRequest("POST", "/v1/apps/", body)

	if err != nil {
		return api.App{}, err
	}

	app := api.App{}
	if err = json.Unmarshal([]byte(resBody), &app); err != nil {
		return api.App{}, err
	}

	return app, nil
}

// Get app details from a Deis controller.
func Get(c *client.Client, appID string) (api.App, error) {
	u := fmt.Sprintf("/v1/apps/%s/", appID)

	body, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return api.App{}, err
	}

	app := api.App{}

	if err = json.Unmarshal([]byte(body), &app); err != nil {
		return api.App{}, err
	}

	return app, nil
}

// Logs retrieves logs from an app.
func Logs(c *client.Client, appID string, lines int) (string, error) {
	u := fmt.Sprintf("/v1/apps/%s/logs", appID)

	if lines > 0 {
		u += "?log_lines=" + strconv.Itoa(lines)
	}

	body, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return "", err
	}

	return strings.Trim(body, `"`), nil
}

// Run one time command in an app.
func Run(c *client.Client, appID string, command string) (api.AppRunResponse, error) {
	req := api.AppRunRequest{Command: command}
	body, err := json.Marshal(req)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	u := fmt.Sprintf("/v1/apps/%s/run", appID)

	resBody, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return api.AppRunResponse{}, err
	}

	out := make([]interface{}, 2)

	if err = json.Unmarshal([]byte(resBody), &out); err != nil {
		return api.AppRunResponse{}, err
	}

	return api.AppRunResponse{Output: out[1].(string), ReturnCode: int(out[0].(float64))}, nil
}

// Delete an app.
func Delete(c *client.Client, appID string) error {
	u := fmt.Sprintf("/v1/apps/%s/", appID)

	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}

// Transfer an app to another user.
func Transfer(c *client.Client, appID string, username string) error {
	u := fmt.Sprintf("/v1/apps/%s/", appID)

	req := api.AppUpdateRequest{Owner: username}
	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", u, body)
	return err
}
