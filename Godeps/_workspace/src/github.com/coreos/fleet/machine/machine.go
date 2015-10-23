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

package machine

import (
	"github.com/coreos/fleet/log"
	"github.com/coreos/fleet/pkg"
)

type Machine interface {
	State() MachineState
}

// HasMetadata determine if the Metadata of a given MachineState
// matches the indicated values.
func HasMetadata(state *MachineState, metadata map[string]pkg.Set) bool {
	for key, values := range metadata {
		local, ok := state.Metadata[key]
		if !ok {
			log.Debugf("No local values found for Metadata(%s)", key)
			return false
		}

		log.Debugf("Asserting local Metadata(%s) meets requirements", key)

		if values.Contains(local) {
			log.Debugf("Local Metadata(%s) meets requirement", key)
		} else {
			log.Debugf("Local Metadata(%s) does not match requirement", key)
			return false
		}
	}

	return true
}
