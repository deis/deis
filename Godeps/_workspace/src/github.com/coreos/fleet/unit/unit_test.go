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

package unit

import (
	"reflect"
	"testing"
)

func TestUnitHash(t *testing.T) {
	u, err := NewUnitFile("[Service]\nExecStart=/bin/sleep 100\n")
	if err != nil {
		t.Fatalf("Unexpected error encountered creating unit: %v", err)
	}

	gotHash := u.Hash()
	gotHashString := gotHash.String()
	expectHashString := "1c6fb6f3684bafb0c173d8b8b957ceff031180c1"
	if gotHashString != expectHashString {
		t.Fatalf("Unit Hash (%s) does not match expected (%s)", gotHashString, expectHashString)
	}

	expectHashShort := "1c6fb6f"
	gotHashShort := gotHash.Short()
	if gotHashShort != expectHashShort {
		t.Fatalf("Unit Hash short (%s) does not match expected (%s)", gotHashShort, expectHashShort)
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
	want := &UnitState{
		LoadState:   "ls",
		ActiveState: "as",
		SubState:    "ss",
		MachineID:   "id",
	}

	got := NewUnitState("ls", "as", "ss", "id")
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
			t.Errorf("NewUnitNameInfo(%s).IsInstance returned %t, want %t", tt.name, i, tt.isinst)
		}
	}

	bad := []string{"foo", "bar@baz"}
	for _, tt := range bad {
		if NewUnitNameInfo(tt) != nil {
			t.Errorf("NewUnitNameInfo returned non-nil InstanceUnit unexpectedly")
		}
	}

}

func TestDeserialize(t *testing.T) {
	contents := `
This=Ignored
[Unit]
;ignore this guy
Description = Foo

[Service]
ExecStart=echo "ping";
ExecStop=echo "pong"
# ignore me, too
ExecStop=echo post

[X-Fleet]
MachineMetadata=foo=bar
MachineMetadata=baz=qux
`

	expected := map[string]map[string][]string{
		"Unit": {
			"Description": {"Foo"},
		},
		"Service": {
			"ExecStart": {`echo "ping";`},
			"ExecStop":  {`echo "pong"`, "echo post"},
		},
		"X-Fleet": {
			"MachineMetadata": {"foo=bar", "baz=qux"},
		},
	}

	unitFile, err := NewUnitFile(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}

	if !reflect.DeepEqual(expected, unitFile.Contents) {
		t.Fatalf("Map func did not produce expected output.\nActual=%v\nExpected=%v", unitFile.Contents, expected)
	}
}

func TestDeserializedUnitGarbage(t *testing.T) {
	contents := `
>>>>>>>>>>>>>
[Service]
ExecStart=jim
# As long as a line has an equals sign, systemd is happy, so we should pass it through
<<<<<<<<<<<=bar
`
	expected := map[string]map[string][]string{
		"Service": {
			"ExecStart":   {"jim"},
			"<<<<<<<<<<<": {"bar"},
		},
	}
	unitFile, err := NewUnitFile(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}

	if !reflect.DeepEqual(expected, unitFile.Contents) {
		t.Fatalf("Map func did not produce expected output.\nActual=%v\nExpected=%v", unitFile.Contents, expected)
	}
}

func TestDeserializeEscapedMultilines(t *testing.T) {
	contents := `
[Service]
ExecStart=echo \
  "pi\
  ng"
ExecStop=\
echo "po\
ng"
# comments within continuation should not be ignored
ExecStopPre=echo\
#pang
ExecStopPost=echo\
#peng\
pung
`
	expected := map[string]map[string][]string{
		"Service": {
			"ExecStart":    {`echo    "pi   ng"`},
			"ExecStop":     {`echo "po ng"`},
			"ExecStopPre":  {`echo #pang`},
			"ExecStopPost": {`echo #peng pung`},
		},
	}
	unitFile, err := NewUnitFile(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}

	if !reflect.DeepEqual(expected, unitFile.Contents) {
		t.Fatalf("Map func did not produce expected output.\nActual=%v\nExpected=%v", unitFile.Contents, expected)
	}
}

func TestSerializeDeserialize(t *testing.T) {
	contents := `
[Unit]
Description = Foo
`
	deserialized, err := NewUnitFile(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}
	section := deserialized.Contents["Unit"]
	if val, ok := section["Description"]; !ok || val[0] != "Foo" {
		t.Errorf("Failed to persist data through serialize/deserialize: %v", val)
	}

	serialized := deserialized.String()
	deserialized, err = NewUnitFile(serialized)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", serialized, err)
	}

	section = deserialized.Contents["Unit"]
	if val, ok := section["Description"]; !ok || val[0] != "Foo" {
		t.Errorf("Failed to persist data through serialize/deserialize: %v", val)
	}
}

func TestDescription(t *testing.T) {
	contents := `
[Unit]
Description = Foo

[Service]
ExecStart=echo "ping";
ExecStop=echo "pong";
`

	unitFile, err := NewUnitFile(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}
	if unitFile.Description() != "Foo" {
		t.Fatalf("Unit.Description is incorrect")
	}
}

func TestDescriptionNotDefined(t *testing.T) {
	contents := `
[Unit]

[Service]
ExecStart=echo "ping";
ExecStop=echo "pong";
`

	unitFile, err := NewUnitFile(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}
	if unitFile.Description() != "" {
		t.Fatalf("Unit.Description is incorrect")
	}
}

func TestBadUnitsFail(t *testing.T) {
	bad := []string{
		`
[Unit]

[Service]
<<<<<<<<<<<<<<<<
`,
		`
[Unit]
nonsense upon stilts
`,
	}
	for _, tt := range bad {
		if _, err := NewUnitFile(tt); err == nil {
			t.Fatalf("Did not get expected error creating Unit from %q", tt)
		}
	}
}

func TestParseMultivalueLine(t *testing.T) {
	tests := []struct {
		in  string
		out []string
	}{
		{`"bar" "ping" "pong"`, []string{`bar`, `ping`, `pong`}},
		{`"bar"`, []string{`bar`}},
		{``, []string{""}},
		{`""`, []string{``}},
		{`"bar`, []string{`"bar`}},
		{`bar"`, []string{`bar"`}},
		{`foo\"bar`, []string{`foo\"bar`}},

		{
			`"bar" "`,
			[]string{`bar`, ``},
			//TODO(bcwaldon): should be something like this:
			// []string{`bar`},
		},

		{
			`"foo\"bar"`,
			[]string{`foo\bar`},
			//TODO(bcwaldon): should be something like this:
			// []string{`foo\"bar`},
		},
	}
	for i, tt := range tests {
		out := parseMultivalueLine(tt.in)
		if !reflect.DeepEqual(tt.out, out) {
			t.Errorf("case %d:, epected %v, got %v", i, tt.out, out)
		}
	}
}
