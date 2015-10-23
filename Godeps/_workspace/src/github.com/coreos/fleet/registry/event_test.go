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
	"reflect"
	"testing"

	"github.com/coreos/fleet/etcd"
	"github.com/coreos/fleet/pkg"
)

func TestFilterEtcdEvents(t *testing.T) {
	tests := []struct {
		in string
		ev pkg.Event
		ok bool
	}{
		{
			in: "",
			ok: false,
		},
		{
			in: "/",
			ok: false,
		},
		{
			in: "/fleet",
			ok: false,
		},
		{
			in: "/fleet/job",
			ok: false,
		},
		{
			in: "/fleet/job/foo/object",
			ok: false,
		},
		{
			in: "/fleet/machine/asdf",
			ok: false,
		},
		{
			in: "/fleet/state/asdf",
			ok: false,
		},
		{
			in: "/fleet/job/foobarbaz/target-state",
			ev: JobTargetStateChangeEvent,
			ok: true,
		},
		{
			in: "/fleet/job/asdf/target",
			ev: JobTargetChangeEvent,
			ok: true,
		},
	}

	for i, tt := range tests {
		for _, action := range []string{"set", "update", "create", "delete"} {
			prefix := "/fleet"

			res := &etcd.Result{
				Node: &etcd.Node{
					Key: tt.in,
				},
				Action: action,
			}
			ev, ok := parse(res, prefix)
			if ok != tt.ok {
				t.Errorf("case %d: expected ok=%t, got %t", i, tt.ok, ok)
				continue
			}

			if !reflect.DeepEqual(tt.ev, ev) {
				t.Errorf("case %d: received incorrect event\nexpected %#v\ngot %#v", i, tt.ev, ev)
				t.Logf("action: %v", action)
			}
		}
	}
}
