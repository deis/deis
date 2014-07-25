package unit

import (
	"reflect"
	"testing"
)

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

[Fleet]
X-ConditionMachineMetadata=foo=bar
X-ConditionMachineMetadata=baz=qux
`

	expected := map[string]map[string][]string{
		"Unit": {
			"Description": {"Foo"},
		},
		"Service": {
			"ExecStart": {`echo "ping";`},
			"ExecStop":  {`echo "pong"`, "echo post"},
		},
		"Fleet": {
			"X-ConditionMachineMetadata": {"foo=bar", "baz=qux"},
		},
	}

	unitFile, err := NewUnit(contents)
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
	unitFile, err := NewUnit(contents)
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
	unitFile, err := NewUnit(contents)
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
	deserialized, err := NewUnit(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}
	section := deserialized.Contents["Unit"]
	if val, ok := section["Description"]; !ok || val[0] != "Foo" {
		t.Errorf("Failed to persist data through serialize/deserialize: %v", val)
	}

	serialized := deserialized.String()
	deserialized, err = NewUnit(serialized)
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

	unitFile, err := NewUnit(contents)
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

	unitFile, err := NewUnit(contents)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", contents, err)
	}
	if unitFile.Description() != "" {
		t.Fatalf("Unit.Description is incorrect")
	}
}

func TestNewSystemdUnitFileFromLegacyContents(t *testing.T) {
	legacy := map[string]map[string]string{
		"Unit": {
			"Description": "foobar",
		},
		"Service": {
			"Type":      "oneshot",
			"ExecStart": "/usr/bin/echo bar",
		},
	}

	expected := map[string]map[string][]string{
		"Unit": {
			"Description": {"foobar"},
		},
		"Service": {
			"Type":      {"oneshot"},
			"ExecStart": {"/usr/bin/echo bar"},
		},
	}

	u, err := NewUnitFromLegacyContents(legacy)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", legacy, err)
	}
	actual := u.Contents

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Map func did not produce expected output.\nActual=%v\nExpected=%v", actual, expected)
	}
}

func TestDeserializeLine(t *testing.T) {
	deserializeLineExamples := map[string][]string{
		`key=foo=bar`:             {`foo=bar`},
		`key="foo=bar"`:           {`foo=bar`},
		`key="foo=bar" "baz=qux"`: {`foo=bar`, `baz=qux`},
		`key="foo=bar baz"`:       {`foo=bar baz`},
		`key="foo=bar" baz`:       {`"foo=bar" baz`},
		`key=baz "foo=bar"`:       {`baz "foo=bar"`},
		`key="foo=bar baz=qux"`:   {`foo=bar baz=qux`},
	}

	for q, w := range deserializeLineExamples {
		k, g, err := deserializeUnitLine(q)
		if err != nil {
			t.Fatalf("Unexpected error testing %q: %v", q, err)
		}
		if k != "key" {
			t.Fatalf("Unexpected key, got %q, want %q", k, "key")
		}
		if !reflect.DeepEqual(g, w) {
			t.Errorf("Unexpected line parse for %q:\ngot %q\nwant %q", q, g, w)
		}
	}

	// Any non-empty line without an '=' is bad
	badLines := []string{
		`<<<<<<<<<<<<<<<<<<<<<<<<`,
		`asdjfkl;`,
		`>>>>>>>>>>>>>>>>>>>>>>>>`,
		`!@#$%^&&*`,
	}
	for _, l := range badLines {
		_, _, err := deserializeUnitLine(l)
		if err == nil {
			t.Fatalf("Did not get expected error deserializing %q", l)
		}
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
		if _, err := NewUnit(tt); err == nil {
			t.Fatalf("Did not get expected error creating Unit from %q", tt)
		}
	}
}
