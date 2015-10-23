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
	"testing"

	"github.com/coreos/fleet/pkg"
)

func TestHasMetadata(t *testing.T) {
	testCases := []struct {
		metadata map[string]string
		match    map[string]pkg.Set
		want     bool
	}{
		{
			map[string]string{
				"region": "us-east-1",
			},
			map[string]pkg.Set{
				"region": pkg.NewUnsafeSet("us-east-1"),
			},
			true,
		},
		{
			map[string]string{
				"groups": "ping",
			},
			map[string]pkg.Set{
				"groups": pkg.NewUnsafeSet("ping", "pong"),
			},
			true,
		},
		{
			map[string]string{
				"groups": "ping",
			},
			map[string]pkg.Set{
				"groups": pkg.NewUnsafeSet("pong"),
			},
			false,
		},
		{
			map[string]string{
				"region": "us-east-1",
				"groups": "ping",
			},
			map[string]pkg.Set{
				"region": pkg.NewUnsafeSet("us-east-1"),
				"groups": pkg.NewUnsafeSet("pong"),
			},
			false,
		},
		{
			map[string]string{},
			map[string]pkg.Set{
				"groups": pkg.NewUnsafeSet("pong"),
			},
			false,
		},
	}

	for i, tt := range testCases {
		ms := &MachineState{Metadata: tt.metadata}
		got := HasMetadata(ms, tt.match)
		if got != tt.want {
			t.Errorf("case %d: HasMetadata returned %t, expected %t", i, got, tt.want)
		}
	}
}
