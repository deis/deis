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
	"path"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/log"
	"github.com/coreos/fleet/unit"
)

const (
	unitPrefix = "/unit/"
)

func (r *EtcdRegistry) storeOrGetUnitFile(u unit.UnitFile) (err error) {
	um := unitModel{
		Raw: u.String(),
	}

	json, err := marshal(um)
	if err != nil {
		return err
	}

	req := etcd.Create{
		Key:   r.hashedUnitPath(u.Hash()),
		Value: json,
	}
	_, err = r.etcd.Do(&req)
	// unit is already stored
	if err != nil && etcd.IsNodeExist(err) {
		// TODO(jonboulle): verify more here?
		err = nil
	}
	return
}

// getUnitByHash retrieves from the Registry the Unit associated with the given Hash
func (r *EtcdRegistry) getUnitByHash(hash unit.Hash) *unit.UnitFile {
	req := etcd.Get{
		Key:       r.hashedUnitPath(hash),
		Recursive: true,
	}
	resp, err := r.etcd.Do(&req)
	if err != nil {
		if etcd.IsKeyNotFound(err) {
			err = nil
		}
		return nil
	}
	var um unitModel
	if err := unmarshal(resp.Node.Value, &um); err != nil {
		log.Errorf("error unmarshaling Unit(%s): %v", hash, err)
		return nil
	}

	u, err := unit.NewUnitFile(um.Raw)
	if err != nil {
		log.Errorf("error parsing Unit(%s): %v", hash, err)
		return nil
	}

	return u
}

func (r *EtcdRegistry) hashedUnitPath(hash unit.Hash) string {
	return path.Join(r.keyPrefix, unitPrefix, hash.String())
}

type unitModel struct {
	Raw string
}
