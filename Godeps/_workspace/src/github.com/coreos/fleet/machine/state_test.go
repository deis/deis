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

import "testing"

func TestStackState(t *testing.T) {
	top := MachineState{
		ID:       "c31e44e1-f858-436e-933e-59c642517860",
		PublicIP: "1.2.3.4",
		Metadata: map[string]string{"ping": "pong"},
		Version:  "1",
	}
	bottom := MachineState{
		ID:       "595989bb-cbb7-49ce-8726-722d6e157b4e",
		PublicIP: "5.6.7.8",
		Metadata: map[string]string{"foo": "bar"},
		Version:  "",
	}
	stacked := stackState(top, bottom)

	if stacked.ID != "c31e44e1-f858-436e-933e-59c642517860" {
		t.Errorf("Unexpected ID value %s", stacked.ID)
	}

	if stacked.PublicIP != "1.2.3.4" {
		t.Errorf("Unexpected PublicIp value %s", stacked.PublicIP)
	}

	if len(stacked.Metadata) != 1 || stacked.Metadata["ping"] != "pong" {
		t.Errorf("Unexpected Metadata %v", stacked.Metadata)
	}

	if stacked.Version != "1" {
		t.Errorf("Unexpected Version value %s", stacked.Version)
	}
}

func TestStackStateEmptyTop(t *testing.T) {
	top := MachineState{}
	bottom := MachineState{
		ID:       "595989bb-cbb7-49ce-8726-722d6e157b4e",
		PublicIP: "5.6.7.8",
		Metadata: map[string]string{"foo": "bar"},
	}
	stacked := stackState(top, bottom)

	if stacked.ID != "595989bb-cbb7-49ce-8726-722d6e157b4e" {
		t.Errorf("Unexpected ID value %s", stacked.ID)
	}

	if stacked.PublicIP != "5.6.7.8" {
		t.Errorf("Unexpected PublicIP value %s", stacked.PublicIP)
	}

	if len(stacked.Metadata) != 1 || stacked.Metadata["foo"] != "bar" {
		t.Errorf("Unexpected Metadata %v", stacked.Metadata)
	}

	if stacked.Version != "" {
		t.Errorf("Unexpected Version value %s", stacked.Version)
	}
}

var shortIDTests = []struct {
	m MachineState
	s string
	l string
}{
	{
		m: MachineState{},
		s: "",
		l: "",
	},
	{
		m: MachineState{
			"595989bb-cbb7-49ce-8726-722d6e157b4e",
			"5.6.7.8",
			map[string]string{"foo": "bar"},
			"",
		},
		s: "595989bb",
		l: "595989bb-cbb7-49ce-8726-722d6e157b4e",
	},
	{
		m: MachineState{
			ID:       "5959",
			PublicIP: "5.6.7.8",
			Metadata: map[string]string{"foo": "bar"},
		},
		s: "5959",
		l: "5959",
	},
}

func TestStateShortID(t *testing.T) {
	for i, tt := range shortIDTests {
		if g := tt.m.ShortID(); g != tt.s {
			t.Errorf("#%d: got %q, want %q", i, g, tt.s)
		}
	}
}

func TestStateMatchID(t *testing.T) {
	for i, tt := range shortIDTests {
		if tt.s != "" {
			if ok := tt.m.MatchID(""); ok {
				t.Errorf("#%d: expected %v", i, false)
			}
		}

		if ok := tt.m.MatchID("foobar"); ok {
			t.Errorf("#%d: expected %v", i, false)
		}

		if ok := tt.m.MatchID(tt.l); !ok {
			t.Errorf("#%d: expected %v", i, true)
		}

		if ok := tt.m.MatchID(tt.s); !ok {
			t.Errorf("#%d: expected %v", i, true)
		}
	}
}
