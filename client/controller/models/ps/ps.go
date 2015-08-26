package ps

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
)

// List an app's processes.
func List(c *client.Client, appID string, results int) ([]api.Process, int, error) {
	u := fmt.Sprintf("/v1/apps/%s/containers/", appID)
	body, count, err := c.LimitedRequest(u, results)

	if err != nil {
		return []api.Process{}, -1, err
	}

	var procs []api.Process
	if err = json.Unmarshal([]byte(body), &procs); err != nil {
		return []api.Process{}, -1, err
	}

	return procs, count, nil
}

// Scale an app's processes.
func Scale(c *client.Client, appID string, targets map[string]int) error {
	u := fmt.Sprintf("/v1/apps/%s/scale/", appID)

	body, err := json.Marshal(targets)

	if err != nil {
		return err
	}

	_, err = c.BasicRequest("POST", u, body)
	return err
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

	body, err := c.BasicRequest("POST", u, nil)

	if err != nil {
		return []api.Process{}, err
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
