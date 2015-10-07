package fleet

import (
	"bytes"
	"testing"
	"text/tabwriter"

	"github.com/coreos/fleet/machine"
)

func TestListMachines(t *testing.T) {
	testMachines := []machine.MachineState{
		machine.MachineState{
			ID:       "123456",
			PublicIP: "1.1.1.1",
			Metadata: map[string]string{
				"foo":  "bar",
				"ping": "pong",
			},
			Version: "",
		},
		machine.MachineState{
			ID:       "654321",
			PublicIP: "2.2.2.2",
			Metadata: nil,
			Version:  "",
		},
	}

	testWriter := bytes.Buffer{}
	testTabWriter := new(tabwriter.Writer)
	testTabWriter.Init(&testWriter, 0, 8, 1, '\t', 0)

	c := &FleetClient{
		Fleet: &stubFleetClient{
			testMachineStates: testMachines,
		},
		out: testTabWriter,
	}

	err := c.ListMachines()

	if err != nil {
		t.Fatal(err)
	}

	expected := `MACHINE		IP	METADATA
123456...	1.1.1.1	foo=bar,ping=pong
654321...	2.2.2.2	-
`

	actual := testWriter.String()

	if expected != actual {
		t.Errorf("Expected '%s', Got '%s'", expected, actual)
	}
}
