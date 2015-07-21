package fleet

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

type mockStatusCommandRunner struct {
	validUnits []string
}

func (mockStatusCommandRunner) LocalCommand(string) (int, error) {
	return 0, nil
}

func (m mockStatusCommandRunner) RemoteCommand(cmd string, addr string, timeout time.Duration) (int, error) {
	if addr != "1.1.1.1" || timeout != 0 {
		return -1, fmt.Errorf("Got %s %d, which is unexpected", cmd, addr, timeout)
	}

	for _, unit := range m.validUnits {
		if fmt.Sprintf("systemctl status -l %s", unit) == cmd {
			return 0, nil
		}
	}

	return -1, fmt.Errorf("Didn't find command %s to match with units %v", cmd, m.validUnits)
}

func TestStatus(t *testing.T) {
	t.Parallel()

	testMachines := []machine.MachineState{
		machine.MachineState{
			ID:       "test-1",
			PublicIP: "1.1.1.1",
			Metadata: nil,
			Version:  "",
		},
	}

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name:         "deis-router@1.service",
			CurrentState: "loaded",
			MachineID:    "test-1",
		},
		&schema.Unit{
			Name:         "deis-router@2.service",
			CurrentState: "loaded",
			MachineID:    "test-1",
		},
		&schema.Unit{
			Name:      "deis-controller.service",
			MachineID: "test-1",
		},
	}

	testWriter := bytes.Buffer{}

	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits, testMachineStates: testMachines,
		unitsMutex: &sync.Mutex{}}, errWriter: &testWriter, runner: mockStatusCommandRunner{
		validUnits: []string{"deis-router@1.service", "deis-router@2.service"}}}

	err := c.Status("router")

	if err != nil {
		t.Error(err)
	}

	commandErr := testWriter.String()

	if commandErr != "" {
		t.Error(commandErr)
	}

	c.runner = mockStatusCommandRunner{validUnits: []string{}}
	err = c.Status("foo")

	actualErr := err.Error()
	expectedErr := "could not find unit: foo"

	if actualErr != expectedErr {
		t.Errorf("Expected %s, Got %s", expectedErr, actualErr)
	}

	err = c.Status("controller")

	if err != nil {
		t.Error(err)
	}

	expectedErr = "Unit deis-controller.service does not appear to be running.\n"
	commandErr = testWriter.String()

	if commandErr != expectedErr {
		t.Errorf("Expected '%s', Got '%s'", expectedErr, commandErr)
	}
}
