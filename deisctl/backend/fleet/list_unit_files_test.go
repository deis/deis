package fleet

import (
	"bytes"
	"sync"
	"testing"
	"text/tabwriter"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

func TestListUnitFiles(t *testing.T) {
	t.Parallel()

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name:         "deis-controller.service",
			DesiredState: "launched",
			CurrentState: "launched",
			MachineID:    "123456",
			Options: []*schema.UnitOption{
				&schema.UnitOption{
					Section: "Unit",
					Name:    "Description",
					Value:   "deis-controller",
				},
			},
		},
		&schema.Unit{
			Name:         "deis-router@1.service",
			DesiredState: "launched",
			CurrentState: "launched",
			MachineID:    "654321",
			Options: []*schema.UnitOption{
				&schema.UnitOption{
					Section: "Unit",
					Name:    "Description",
					Value:   "deis-router@1",
				},
			},
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

	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits,
		testMachineStates: testMachines, unitsMutex: &sync.Mutex{}},
		out: testTabWriter}

	err := c.ListUnitFiles()

	if err != nil {
		t.Fatal(err)
	}

	expected := `UNIT			HASH	DSTATE		STATE		TMACHINE
deis-controller.service	61030e3	launched	launched	123456.../1.1.1.1
deis-router@1.service	1182ecf	launched	launched	654321.../2.2.2.2
`

	actual := testWriter.String()

	if expected != actual {
		t.Errorf("Expected '%s', Got '%s'", expected, actual)
	}
}
