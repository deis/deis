package job

import (
	"testing"

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

func TestJobWithPeers(t *testing.T) {
	j := NewJob("echo.service", *newUnit(t, ``))
	peers := j.Peers()

	if len(peers) != 0 {
		t.Fatalf("Unexpected number of peers %d, expected 0", len(peers))
	}
}

func TestJobWithoutPeers(t *testing.T) {
	contents := `[X-Fleet]
X-ConditionMachineOf="foo.service" "bar.service"
`
	j := NewJob("echo.service", *newUnit(t, contents))
	peers := j.Peers()

	if len(peers) != 2 {
		t.Fatalf("Unexpected number of peers %d, expected 2", len(peers))
	}

	if peers[0] != "foo.service" {
		t.Errorf("Expected first peer to be foo.service, got %s", peers[0])
	}

	if peers[1] != "bar.service" {
		t.Errorf("Expected second peer to be bar.service, got %s", peers[1])
	}
}

func TestJobConflicts(t *testing.T) {
	contents := `[Unit]
Description=Testing

[X-Fleet]
X-Conflicts=*bar*
`
	j := NewJob("echo.service", *newUnit(t, contents))
	conflicts := j.Conflicts()

	if len(conflicts) != 1 {
		t.Errorf("Expected 1 conflict, received %v", conflicts)
	}

	if conflicts[0] != "*bar*" {
		t.Errorf("Expected first conflict to be '*bar*', received %s", conflicts[1])
	}
}

func TestJobConflictsNotProvided(t *testing.T) {
	j := NewJob("echo.socket", *newUnit(t, ""))
	conflicts := j.Conflicts()

	if len(conflicts) > 0 {
		t.Fatalf("Expected no conflicts, received %v", conflicts)
	}
}

func TestParseRequirements(t *testing.T) {
	contents := `
[X-Fleet]
X-Foo=Bar
Ping=Pong
X-Key=Value
`
	j := NewJob("foo.service", *newUnit(t, contents))
	reqs := j.Requirements()
	if len(reqs) != 2 {
		t.Fatalf("Incorrect number of requirements; got %d, expected 2", len(reqs))
	}

	if len(reqs["Foo"]) != 1 || reqs["Foo"][0] != "Bar" {
		t.Fatalf("Incorrect value %q of requirement 'Foo'", reqs["Foo"])
	}

	if len(reqs["Key"]) != 1 || reqs["Key"][0] != "Value" {
		t.Fatalf("Incorrect value %q of requirement 'Key'", reqs["Key"])
	}
}

func TestParseRequirementsMultipleValuesForKeyStack(t *testing.T) {
	contents := `
[X-Fleet]
X-Foo=Bar
X-Foo=Baz
X-Ping=Pong
X-Ping=Pang
`
	j := NewJob("foo.service", *newUnit(t, contents))
	reqs := j.Requirements()
	if len(reqs) != 2 {
		t.Fatalf("Incorrect number of requirements; got %d, expected 2: %v", len(reqs), reqs)
	}

	if len(reqs["Foo"]) != 2 || reqs["Foo"][0] != "Bar" || reqs["Foo"][1] != "Baz" {
		t.Fatalf("Incorrect value %v of requirement 'Foo'", reqs["Foo"])
	}

	if len(reqs["Ping"]) != 2 || reqs["Ping"][0] != "Pong" || reqs["Ping"][1] != "Pang" {
		t.Fatalf("Incorrect value %v of requirement 'Ping'", reqs["Ping"])
	}
}

func TestParseRequirementsInstanceUnit(t *testing.T) {
	contents := `
[X-Fleet]
X-Foo=%n
X-Bar=%N
X-Baz=%p
X-Qux=%i
X-Zzz=something
`
	// Ensure the correct values are replaced for a non-instance unit
	j := NewJob("test.service", *newUnit(t, contents))
	reqs := j.Requirements()
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
	reqs = j.Requirements()
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
	reqs := j.Requirements()
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
		// Simplest case
		{
			`[X-Fleet]
X-ConditionMachineID=123
`,
			"123",
			true,
		},

		// First value wins
		// TODO(bcwaldon): maybe the last one should win?
		{
			`[X-Fleet]
X-ConditionMachineID="123" "456"
`,
			"123",
			true,
		},

		// No value provided
		{
			`[X-Fleet]`,
			"",
			false,
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
	}

	for _, tt := range tests {
		j := NewJob("echo.service", *newUnit(t, tt.unit))
		outS, outB := j.RequiredTarget()

		if outS != tt.outS {
			t.Errorf("Expected target requirement %s, got %s", tt.outS, outS)
		}

		if outB != tt.outB {
			t.Errorf("Expected target requirement ok-val %t, got %t", tt.outB, outB)
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
		{"[X-Fleet]\nX-ConditionMachineOf=bar", false},
		{"Global=true", false},
		// specified in wrong section
		{"[Service]\nGlobal=true", false},
		// bad values
		{"[X-Fleet]\nX-ConditionMachineOf=bar\nGlobal=false", false},
		{"[X-Fleet]\nX-ConditionMachineOf=bar\nGlobal=what", false},
		// correct specifications
		{"[X-Fleet]\nX-ConditionMachineOf=foo\nGlobal=true", true},
		{"[X-Fleet]\nX-ConditionMachineOf=foo\nGlobal=True", true},
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
