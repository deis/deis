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

const (
	shortIDLen = 8
)

// MachineState represents a point-in-time snapshot of the
// state of the local host.
type MachineState struct {
	ID       string
	PublicIP string
	Metadata map[string]string
	Version  string
}

func (ms MachineState) ShortID() string {
	if len(ms.ID) <= shortIDLen {
		return ms.ID
	}
	return ms.ID[0:shortIDLen]
}

func (ms MachineState) MatchID(ID string) bool {
	return ms.ID == ID || ms.ShortID() == ID
}

// stackState is used to merge two MachineStates. Values configured on the top
// MachineState always take precedence over those on the bottom.
func stackState(top, bottom MachineState) MachineState {
	state := MachineState(bottom)

	if top.PublicIP != "" {
		state.PublicIP = top.PublicIP
	}

	if top.ID != "" {
		state.ID = top.ID
	}

	//FIXME: This will *always* overwrite the bottom's metadata,
	// but the only use-case we have today does not ever have
	// metadata on the bottom.
	if len(top.Metadata) > 0 {
		state.Metadata = top.Metadata
	}

	if top.Version != "" {
		state.Version = top.Version
	}

	return state
}
