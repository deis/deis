package fleet

import (
	"bytes"
	"sync"
	"testing"
	"text/tabwriter"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

func TestListUnits(t *testing.T) {
	t.Parallel()

	testUnitStates := []*schema.UnitState{
		&schema.UnitState{
			Name:               "deis-controller.service",
			MachineID:          "123456",
			SystemdLoadState:   "loaded",
			SystemdActiveState: "active",
			SystemdSubState:    "running",
			Hash:               "abcd",
		},
		&schema.UnitState{
			Name:               "deis-router@1.service",
			MachineID:          "654321",
			SystemdLoadState:   "loaded",
			SystemdActiveState: "active",
			SystemdSubState:    "running",
			Hash:               "dcba",
		},
	}

	testMachines := []machine.MachineState{
		machine.MachineState{
			ID:       "123456",
			PublicIP: "1.1.1.1",
			Metadata: nil,
			Version:  "",
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

	c := &FleetClient{Fleet: &stubFleetClient{testUnitStates: testUnitStates,
		testMachineStates: testMachines, unitStatesMutex: &sync.Mutex{}},
		out: testTabWriter}

	err := c.ListUnits()

	if err != nil {
		t.Fatal(err)
	}

	expected := `UNIT			MACHINE			LOAD	ACTIVE	SUB
deis-controller.service	123456.../1.1.1.1	loaded	active	running
deis-router@1.service	654321.../2.2.2.2	loaded	active	running
`

	actual := testWriter.String()

	if expected != actual {
		t.Errorf("Expected '%s', Got '%s'", expected, actual)
	}
}
