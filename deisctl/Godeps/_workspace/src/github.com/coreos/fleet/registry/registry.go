package registry

import (
	"encoding/json"
	"fmt"

	"github.com/coreos/fleet/etcd"
)

const DefaultKeyPrefix = "/_coreos.com/fleet/"

// EtcdRegistry fulfils the Registry interface and uses etcd as a backend
type EtcdRegistry struct {
	etcd      etcd.Client
	keyPrefix string
}

// New creates a new EtcdRegistry with the given parameters
func New(client etcd.Client, keyPrefix string) (registry Registry) {
	return &EtcdRegistry{client, keyPrefix}
}

func marshal(obj interface{}) (string, error) {
	encoded, err := json.Marshal(obj)
	if err == nil {
		return string(encoded), nil
	}
	return "", fmt.Errorf("unable to JSON-serialize object: %s", err)
}

func unmarshal(val string, obj interface{}) error {
	err := json.Unmarshal([]byte(val), &obj)
	if err == nil {
		return nil
	}
	return fmt.Errorf("unable to JSON-deserialize object: %s", err)
}

func isKeyNotFound(err error) bool {
	e, ok := err.(etcd.Error)
	return ok && e.ErrorCode == etcd.ErrorKeyNotFound
}

func isNodeExist(err error) bool {
	e, ok := err.(etcd.Error)
	return ok && e.ErrorCode == etcd.ErrorNodeExist
}
