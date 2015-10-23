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

package job

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/coreos/fleet/pkg"
	"github.com/coreos/fleet/unit"
)

func newUnit(t *testing.T, str string) *unit.UnitFile {
	u, err := unit.NewUnitFile(str)
	if err != nil {
		t.Fatalf("Unexpected error creating unit from %q: %v", str, err)
	}
	return u
}

func TestNewJob(t *testing.T) {
	j1 := NewJob("pong.service", *newUnit(t, "Echo"))

	if j1.Name != "pong.service" {
		t.Error("job.Job.Name != 'pong.service'")
	}
}

func TestJobPeers(t *testing.T) {
	testCases := []struct {
		contents string
		peers    []string
	}{
		// no values should be fine
		{``, []string{}},
		// single value should be fine
		{`[X-Fleet]
MachineOf="lol.service"
`, []string{"lol.service"}},
		// multiple values should be fine
		{`[X-Fleet]
MachineOf="foo.service" "bar.service"
`, []string{"foo.service", "bar.service"}},
	}
	for i, tt := range testCases {
		j := NewJob("echo.service", *newUnit(t, tt.contents))
		peers := j.Peers()
		if !reflect.DeepEqual(peers, tt.peers) {
			t.Errorf("case %d: unexpected peers: got %#v, want %#v", i, peers, tt.peers)
		}

	}
}

func TestJobConflicts(t *testing.T) {
	testCases := []struct {
		contents  string
		conflicts []string
	}{
		{``, []string{}},
		{`[Unit]
Description=Testing

[X-Fleet]
Conflicts=*bar*
`, []string{"*bar*"}},
	}
	for i, tt := range testCases {
		j := NewJob("echo.service", *newUnit(t, tt.contents))
		conflicts := j.Conflicts()
		if !reflect.DeepEqual(conflicts, tt.conflicts) {
			t.Errorf("case %d: unexpected conflicts: got %#v, want %#v", i, conflicts, tt.conflicts)
		}

	}
}

func TestParseRequirements(t *testing.T) {
	testCases := []struct {
		contents string
		reqs     map[string][]string
	}{
		// No reqs should be fine
		{``, map[string][]string{}},
		{`[X-Fleet]`, map[string][]string{}},
		// Simple req is fine
		{
			"[X-Fleet]\nFoo=bar",
			map[string][]string{"Foo": []string{"bar"}},
		},
		// Deprecated prefixes should be unmutated
		{
			"[X-Fleet]\nX-Foo=bar",
			map[string][]string{"X-Foo": []string{"bar"}},
		},
		{
			"[X-Fleet]\nX-ConditionFoo=bar",
			map[string][]string{"X-ConditionFoo": []string{"bar"}},
		},
		// Multiple values are fine
		{
			`[X-Fleet]
Foo=asdf
Foo=qwer
`,
			map[string][]string{"Foo": []string{"asdf", "qwer"}},
		},
		// Multiple values with different prefixes
		{
			`[X-Fleet]
Dog=Woof
X-Dog=Yap
X-ConditionDog=Howl
`,
			map[string][]string{
				"Dog":            []string{"Woof"},
				"X-Dog":          []string{"Yap"},
				"X-ConditionDog": []string{"Howl"},
			},
		},
		// Multiple values stacking with the same key
		{
			`[X-Fleet]
Foo=Bar
Foo=Baz
Ping=Pong
Ping=Pang
`,
			map[string][]string{
				"Foo":  []string{"Bar", "Baz"},
				"Ping": []string{"Pong", "Pang"},
			},
		},
	}
	for i, tt := range testCases {
		j := NewJob("foo.service", *newUnit(t, tt.contents))
		reqs := j.requirements()
		if !reflect.DeepEqual(reqs, tt.reqs) {
			t.Errorf("case %d: incorrect requirements: got %#v, want %#v", i, reqs, tt.reqs)
		}
	}
}

func TestParseRequirementsInstanceUnit(t *testing.T) {
	contents := `
[X-Fleet]
Foo=%n
Bar=%N
Baz=%p
Qux=%i
Zzz=something
`
	// Ensure the correct values are replaced for a non-instance unit
	j := NewJob("test.service", *newUnit(t, contents))
	reqs := j.requirements()
	for field, want := range map[string]string{
		"Foo": "test.service",
		"Bar": "test",
		"Baz": "test",
		"Qux": "",
		"Zzz": "something",
	} {
		got := reqs[field]
		if len(got) != 1 || got[0] != want {
			t.Errorf("Requirement %q unexpectedly altered for non-instance unit: want %q, got %q", field, want, got)
		}
	}

	// Now ensure that they are substituted appropriately for an instance unit
	j = NewJob("ssh@2.service", *newUnit(t, contents))
	reqs = j.requirements()
	for field, want := range map[string]string{
		"Foo": "ssh@2.service",
		"Bar": "ssh@2",
		"Baz": "ssh",
		"Qux": "2",
		"Zzz": "something",
	} {
		got := reqs[field]
		if len(got) != 1 || got[0] != want {
			t.Errorf("Bad instance unit requirement substitution for %q: want %q, got %q", field, want, got)
		}
	}

}

func TestParseRequirementsMissingSection(t *testing.T) {
	contents := `
[Unit]
Description=Timmy
`
	j := NewJob("foo.service", *newUnit(t, contents))
	reqs := j.requirements()
	if len(reqs) != 0 {
		t.Fatalf("Incorrect number of requirements; got %d, expected 0", len(reqs))
	}
}

func TestJobConditionMachineID(t *testing.T) {
	tests := []struct {
		unit string
		outS string
		outB bool
	}{
		// No value provided
		{
			`[X-Fleet]`,
			"",
			false,
		},
		// Simplest case
		{
			`[X-Fleet]
MachineID=123
`,
			"123",
			true,
		},

		// First value wins
		// TODO(bcwaldon): maybe the last one should win?
		{
			`[X-Fleet]
MachineID="123" "456"
`,
			"123",
			true,
		},
		// Ensure we fall back to the deprecated prefix
		{
			`[X-Fleet]
X-ConditionMachineID="123" "456"
`,
			"123",
			true,
		},
		// Fall back to the deprecated prefix only if the modern syntax is absent
		{
			`[X-Fleet]
X-ConditionMachineID="456"
MachineID="123"
`,
			"123",
			true,
		},
		{
			`[X-Fleet]
MachineID="123"
X-ConditionMachineID="456"
`,
			"123",
			true,
		},

		// Ensure we fall back to the legacy boot ID option
		{
			`[X-Fleet]
X-ConditionMachineBootID=123
`,
			"123",
			true,
		},

		// Fall back to legacy option only if non-boot ID is absent
		{
			`[X-Fleet]
X-ConditionMachineBootID=123
X-ConditionMachineID=456
`,
			"456",
			true,
		},
		{
			`[X-Fleet]
X-ConditionMachineBootID=123
MachineID=456
`,
			"456",
			true,
		},
	}

	for i, tt := range tests {
		j := NewJob("echo.service", *newUnit(t, tt.unit))
		outS, outB := j.RequiredTarget()

		if outS != tt.outS {
			t.Errorf("case %d: Expected target requirement %s, got %s", i, tt.outS, outS)
		}

		if outB != tt.outB {
			t.Errorf("case %d: Expected target requirement ok-val %t, got %t", i, tt.outB, outB)
		}
	}
}

func TestJobRequiredMetadata(t *testing.T) {
	testCases := []struct {
		unit string
		out  map[string]pkg.Set
	}{
		// no metadata
		{
			`[X-Fleet]`,
			map[string]pkg.Set{},
		},
		// simplest case - one key/value
		{
			`[X-Fleet]
MachineMetadata=foo=bar`,
			map[string]pkg.Set{
				"foo": pkg.NewUnsafeSet("bar"),
			},
		},
		// multiple different values for a key in one line
		{
			`[X-Fleet]
MachineMetadata="foo=bar" "foo=baz"`,
			map[string]pkg.Set{
				"foo": pkg.NewUnsafeSet("bar", "baz"),
			},
		},
		// multiple different values for a key in different lines
		{
			`[X-Fleet]
MachineMetadata=foo=bar
MachineMetadata=foo=baz
MachineMetadata=foo=asdf`,
			map[string]pkg.Set{
				"foo": pkg.NewUnsafeSet("bar", "baz", "asdf"),
			},
		},
		// multiple different key-value pairs in a single line
		{
			`[X-Fleet]
MachineMetadata="foo=bar" "duck=quack"`,
			map[string]pkg.Set{
				"foo":  pkg.NewUnsafeSet("bar"),
				"duck": pkg.NewUnsafeSet("quack"),
			},
		},
		// multiple different key-value pairs in different lines
		{
			`[X-Fleet]
MachineMetadata=foo=bar
MachineMetadata=dog=woof
MachineMetadata=cat=miaow`,
			map[string]pkg.Set{
				"foo": pkg.NewUnsafeSet("bar"),
				"dog": pkg.NewUnsafeSet("woof"),
				"cat": pkg.NewUnsafeSet("miaow"),
			},
		},
		// support deprecated prefixed syntax
		{
			`[X-Fleet]
X-ConditionMachineMetadata=foo=bar`,
			map[string]pkg.Set{
				"foo": pkg.NewUnsafeSet("bar"),
			},
		},
		// support deprecated prefixed syntax mixed with modern syntax
		{
			`[X-Fleet]
MachineMetadata=foo=bar
X-ConditionMachineMetadata=foo=asdf`,
			map[string]pkg.Set{
				"foo": pkg.NewUnsafeSet("bar", "asdf"),
			},
		},
		// bad fields just get ignored
		{
			`[X-Fleet]
MachineMetadata=foo=`,
			map[string]pkg.Set{},
		},
		{
			`[X-Fleet]
MachineMetadata==asdf`,
			map[string]pkg.Set{},
		},
		{
			`[X-Fleet]
MachineMetadata=foo=asdf=WHAT`,
			map[string]pkg.Set{},
		},
		// mix everything up
		{
			`[X-Fleet]
MachineMetadata=ignored=
MachineMetadata=oh=yeah
MachineMetadata=whynot=zoidberg
X-ConditionMachineMetadata=oh=no
X-ConditionMachineMetadata="one=abc" "two=def"`,
			map[string]pkg.Set{
				"oh":     pkg.NewUnsafeSet("yeah", "no"),
				"whynot": pkg.NewUnsafeSet("zoidberg"),
				"one":    pkg.NewUnsafeSet("abc"),
				"two":    pkg.NewUnsafeSet("def"),
			},
		},
	}
	for i, tt := range testCases {
		j := NewJob("echo.service", *newUnit(t, tt.unit))
		md := j.RequiredTargetMetadata()
		if !reflect.DeepEqual(md, tt.out) {
			t.Errorf("case %d: metadata differs", i)
			t.Logf("got: %#v", md)
			t.Logf("want: %#v", tt.out)
		}
	}
}

func TestInstanceUnitPrintf(t *testing.T) {
	u := unit.NewUnitNameInfo("foo@bar.waldo")
	if u == nil {
		t.Fatal("NewNamedUnit returned nil - aborting")
	}
	for _, tt := range []struct {
		in   string
		want string
	}{
		{"%n", "foo@bar.waldo"},
		{"%N", "foo@bar"},
		{"%p", "foo"},
		{"%i", "bar"},
	} {
		got := unitPrintf(tt.in, *u)
		if got != tt.want {
			t.Errorf("Replacement of %q failed: got %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseJobState(t *testing.T) {
	tests := []struct {
		in  string
		out JobState
		err bool
	}{
		{"inactive", JobStateInactive, false},
		{"loaded", JobStateLoaded, false},
		{"launched", JobStateLaunched, false},
		{"active", JobStateInactive, true},
	}

	for i, tt := range tests {
		out, err := ParseJobState(tt.in)
		if (err != nil) != tt.err {
			t.Errorf("case %d: expected error=%t, got %v", i, tt.err, err)
		}
		if out != tt.out {
			t.Errorf("case %d: expected JobState=%v, got %v", i, tt.out, out)
		}
	}
}

func TestJobScheduled(t *testing.T) {
	j1 := NewJob("pong.service", *newUnit(t, "Echo"))

	if j1.Scheduled() {
		t.Error("Job should not be scheduled yet")
	}

	j1.TargetMachineID = "XXX"

	if !j1.Scheduled() {
		t.Error("Job should now be scheduled")
	}
}

func TestUnitIsGlobal(t *testing.T) {
	for i, tt := range []struct {
		contents string
		want     bool
	}{
		// empty unit file
		{"", false},
		// no relevant sections
		{"foobarbaz", false},
		{"[Service]\nExecStart=/bin/true", false},
		{"[X-Fleet]\nMachineOf=bar", false},
		{"Global=true", false},
		// specified in wrong section
		{"[Service]\nGlobal=true", false},
		// bad values
		{"[X-Fleet]\nMachineOf=bar\nGlobal=false", false},
		{"[X-Fleet]\nMachineOf=bar\nGlobal=what", false},
		{"[X-Fleet]\nX-Global=true", false},
		{"[X-Fleet]\nX-ConditionGlobal=true", false},
		// correct specifications
		{"[X-Fleet]\nMachineOf=foo\nGlobal=true", true},
		{"[X-Fleet]\nMachineOf=foo\nGlobal=True", true},
		// multiple parameters - last wins
		{"[X-Fleet]\nGlobal=true\nGlobal=false", false},
		{"[X-Fleet]\nGlobal=false\nGlobal=true", true},
	} {
		u := Unit{
			Unit: *newUnit(t, tt.contents),
		}
		got := u.IsGlobal()
		if got != tt.want {
			t.Errorf("case %d: IsGlobal returned %t, want %t", i, got, tt.want)
		}
	}
}

func TestValidateRequirements(t *testing.T) {
	tests := []string{
		"MachineID=asdf",
		"X-ConditionMachineID=123456",
		"X-ConditionMachineBootID=woofwoof",
		"X-ConditionMachineOf=asdf",
		"MachineOf=joe.service",
		"X-Conflicts=bar.service",
		"Conflicts=foo",
		"X-ConditionMachineMetadata=up=down",
		"MachineMetadata=true=false",
		"Global=true",
	}
	for i, req := range tests {
		contents := fmt.Sprintf("[X-Fleet]\n%s", req)
		j := NewJob("echo.service", *newUnit(t, contents))
		if err := j.ValidateRequirements(); err != nil {
			t.Errorf("case %d: unexpected non-nil error for req %q: %v", i, req, err)
		}
	}
}

func TestBadValidateRequirements(t *testing.T) {
	tests := []string{
		"X-ConditionConflicts=asdf",
		"X-Peers=one",
		"Machineof=something",
		"X-Global=true",
		"global=true",
		"X-ConditionMachineId=one",
		"MachineId=true",
		"X-MachineMetadata=none",
		"X-ConditionMetadata=foo=foo",
	}
	for i, req := range tests {
		contents := fmt.Sprintf("[X-Fleet]\n%s", req)
		j := NewJob("echo.service", *newUnit(t, contents))
		if err := j.ValidateRequirements(); err == nil {
			t.Errorf("case %d: unexpected nil error for requirement: %q", i, req)
		}
	}
}
