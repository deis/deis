package keys

import (
	"encoding/json"
	"fmt"

	"github.com/deis/deis/client-go/controller/api"
	"github.com/deis/deis/client-go/controller/client"
)

// List keys on a controller.
func List(c *client.Client) ([]api.Key, error) {
	body, err := c.BasicRequest("GET", "/v1/keys/", nil)

	if err != nil {
		return []api.Key{}, err
	}

	keys := api.Keys{}
	if err = json.Unmarshal([]byte(body), &keys); err != nil {
		return []api.Key{}, err
	}

	return keys.Keys, nil
}

// New creates a new key.
func New(c *client.Client, id string, pubKey string) (api.Key, error) {
	req := api.KeyCreateRequest{ID: id, Public: pubKey}
	body, err := json.Marshal(req)

	resBody, err := c.BasicRequest("POST", "/v1/keys/", body)

	if err != nil {
		return api.Key{}, err
	}

	key := api.Key{}
	if err = json.Unmarshal([]byte(resBody), &key); err != nil {
		return api.Key{}, err
	}

	return key, nil
}

// Delete a key.
func Delete(c *client.Client, keyID string) error {
	u := fmt.Sprintf("/v1/keys/%s", keyID)

	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}
