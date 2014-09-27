package registry

import (
	"reflect"
	"testing"

	"github.com/coreos/fleet/unit"
)

func TestLegacyPayload(t *testing.T) {
	contents := `
[Service]
ExecStart=/bin/sleep 30000
`[1:]
	legacyPayloadContents := `{"Name":"sleep.service","Unit":{"Contents":{"Service":{"ExecStart":"/bin/sleep 30000"}},"Raw":"[Service]\nExecStart=/bin/sleep 30000\n"}}`
	want, err := unit.NewUnitFile(contents)
	if err != nil {
		t.Fatalf("unexpected error creating unit from %q: %v", contents, err)
	}
	var ljp LegacyJobPayload
	err = unmarshal(legacyPayloadContents, &ljp)
	if err != nil {
		t.Error("Error unmarshaling legacy payload:", err)
	}
	got := ljp.Unit
	if !reflect.DeepEqual(*want, got) {
		t.Errorf("Unit from legacy payload does not match expected!\nwant:\n%s\ngot:\n%s", *want, got)
	}
}
