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
	}

	for i, tt := range testCases {
		ms := &MachineState{Metadata: tt.metadata}
		got := HasMetadata(ms, tt.match)
		if got != tt.want {
			t.Errorf("case %d: HasMetadata returned %t, expected %t", i, got, tt.want)
		}
	}
}
