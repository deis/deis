package ps

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List an app's processes.
func List(c *client.Client, appID string) ([]api.Process, error) {
	u := fmt.Sprintf("/v1/apps/%s/containers/", appID)
	body, status, err := c.BasicRequest("GET", u, nil)

	if err != nil {
		return []api.Process{}, err
	}

	if status != 200 {
		return []api.Process{}, errors.New(body)
	}

	procs := api.Processes{}
	if err = json.Unmarshal([]byte(body), &procs); err != nil {
		return []api.Process{}, err
	}

	return procs.Processes, nil
}

// Scale an app's processes.
func Scale(c *client.Client, appID string, targets map[string]int) error {
	u := fmt.Sprintf("/v1/apps/%s/scale/", appID)

	body, err := json.Marshal(targets)

	if err != nil {
		return err
	}

	resBody, status, err := c.BasicRequest("POST", u, body)

	if err != nil {
		return err
	}

	if status != 204 {
		return errors.New(resBody)
	}

	return nil
}

// Restart an app's processes.
func Restart(c *client.Client, appID string, procType string, num int) ([]api.Process, error) {
	u := fmt.Sprintf("/v1/apps/%s/containers/", appID)

	if procType == "" {
		u += "restart/"
	} else {
		if num == -1 {
			u += procType + "/restart/"
		} else {
			u += procType + "/" + strconv.Itoa(num) + "/restart/"
		}
	}

	body, status, err := c.BasicRequest("POST", u, nil)

	if err != nil {
		return []api.Process{}, err
	}

	if status != 200 {
		return []api.Process{}, errors.New(body)
	}

	procs := []api.Process{}
	if err = json.Unmarshal([]byte(body), &procs); err != nil {
		return []api.Process{}, err
	}

	return procs, nil
}

// ByType organizes processes of an app by process type.
func ByType(processes []api.Process) map[string][]api.Process {
	psMap := make(map[string][]api.Process)

	for _, ps := range processes {
		psMap[ps.Type] = append(psMap[ps.Type], ps)
	}

	return psMap
}
