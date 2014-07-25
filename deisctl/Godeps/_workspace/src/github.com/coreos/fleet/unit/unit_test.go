package unit

import (
	"reflect"
	"testing"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/resource"
)

const (
	// $ echo -n "foo" | sha1sum
	// 0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33 -
	testData      = "foo"
	testShaString = "0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"
	testShaShort  = "0beec7b"
)

func TestUnitHash(t *testing.T) {
	u, err := NewUnit(testData)
	if err != nil {
		t.Fatalf("Unexpected error encountered creating unit: %v", err)
	}
	h := u.Hash()
	if h.String() != testShaString {
		t.Fatalf("Unit Hash (%s) does not match expected (%s)", h.String(), testShaString)
	}

	if h.Short() != testShaShort {
		t.Fatalf("Unit Hash short (%s) does not match expected (%s)", h.Short(), testShaShort)
	}

	eh := &Hash{}
	if !eh.Empty() {
		t.Fatalf("Empty hash check failed: %v", eh.Empty())
	}
}

func TestRecognizedUnitTypes(t *testing.T) {
	tts := []struct {
		name string
		ok   bool
	}{
		{"foo.service", true},
		{"foo.socket", true},
		{"foo.path", true},
		{"foo.timer", true},
		{"foo.mount", true},
		{"foo.automount", true},
		{"foo.device", true},
		{"foo.swap", false},
		{"foo.target", false},
		{"foo.snapshot", false},
		{"foo.network", false},
		{"foo.netdev", false},
		{"foo.link", false},
		{"foo.unknown", false},
	}

	for _, tt := range tts {
		ok := RecognizedUnitType(tt.name)
		if ok != tt.ok {
			t.Errorf("Case failed: name=%s expect=%t result=%t", tt.name, tt.ok, ok)
		}
	}
}

func TestDefaultUnitType(t *testing.T) {
	tts := []struct {
		name string
		out  string
	}{
		{"foo", "foo.service"},
		{"foo.service", "foo.service.service"},
		{"foo.link", "foo.link.service"},
	}

	for _, tt := range tts {
		out := DefaultUnitType(tt.name)
		if out != tt.out {
			t.Errorf("Case failed: name=%s expect=%s result=%s", tt.name, tt.out, out)
		}
	}
}

func TestNewUnitState(t *testing.T) {
	ms := &machine.MachineState{"id", "ip", nil, "version", resource.ResourceTuple{}}
	want := &UnitState{
		LoadState:    "ls",
		ActiveState:  "as",
		SubState:     "ss",
		MachineState: ms,
	}

	got := NewUnitState("ls", "as", "ss", ms)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NewUnitState did not create a correct UnitState: got %s, want %s", got, want)
	}

}

func TestNamedUnit(t *testing.T) {
	tts := []struct {
		fname  string
		name   string
		pref   string
		tmpl   string
		inst   string
		isinst bool
	}{
		{"foo.service", "foo", "foo", "", "", false},
		{"foo@.service", "foo@", "foo", "foo@.service", "", false},
		{"foo@bar.service", "foo@bar", "foo", "foo@.service", "bar", true},
		{"foo@bar@baz.service", "foo@bar@baz", "foo", "foo@.service", "bar@baz", true},
		{"foo.1@.service", "foo.1@", "foo.1", "foo.1@.service", "", false},
		{"foo.1@2.service", "foo.1@2", "foo.1", "foo.1@.service", "2", true},
		{"ssh@.socket", "ssh@", "ssh", "ssh@.socket", "", false},
		{"ssh@1.socket", "ssh@1", "ssh", "ssh@.socket", "1", true},
	}
	for _, tt := range tts {
		u := NewUnitNameInfo(tt.fname)
		if u == nil {
			t.Errorf("NewUnitNameInfo(%s) returned nil InstanceUnit!", tt.name)
			continue
		}
		if u.FullName != tt.fname {
			t.Errorf("NewUnitNameInfo(%s) returned bad fullname: got %s, want %s", tt.name, u.FullName, tt.fname)
		}
		if u.Name != tt.name {
			t.Errorf("NewUnitNameInfo(%s) returned bad name: got %s, want %s", tt.name, u.Name, tt.name)
		}
		if u.Template != tt.tmpl {
			t.Errorf("NewUnitNameInfo(%s) returned bad template name: got %s, want %s", tt.name, u.Template, tt.tmpl)
		}
		if u.Prefix != tt.pref {
			t.Errorf("NewUnitNameInfo(%s) returned bad prefix name: got %s, want %s", tt.name, u.Prefix, tt.pref)
		}
		if u.Instance != tt.inst {
			t.Errorf("NewUnitNameInfo(%s) returned bad instance name: got %s, want %s", tt.name, u.Instance, tt.inst)
		}
		i := u.IsInstance()
		if i != tt.isinst {
			t.Errorf("NewUnitNameInfo(%s).IsInstance returned %t, want %s", tt.name, i, tt.isinst)
		}
	}

	bad := []string{"foo", "bar@baz"}
	for _, tt := range bad {
		if NewUnitNameInfo(tt) != nil {
			t.Errorf("NewUnitNameInfo returned non-nil InstanceUnit unexpectedly")
		}
	}

}
