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
