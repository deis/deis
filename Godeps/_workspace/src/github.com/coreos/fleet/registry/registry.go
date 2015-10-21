// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func NewEtcdRegistry(client etcd.Client, keyPrefix string) *EtcdRegistry {
	return &EtcdRegistry{
		etcd:      client,
		keyPrefix: keyPrefix,
	}
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
