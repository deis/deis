package fleet

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

type mockCommandRunner struct{}

func (mockCommandRunner) LocalCommand(string) (int, error) {
	return 0, nil
}

func (mockCommandRunner) RemoteCommand(cmd string, addr string, timeout time.Duration) (int, error) {
	if cmd != "true" || addr != "1.1.1.1" || timeout != 0 {
		return -1, fmt.Errorf("Got %s %s %d, which is unexpected", cmd, addr, timeout)
	}
	return 0, nil
}

func TestRunCommand(t *testing.T) {
	t.Parallel()

	expected := machine.MachineState{
		ID:       "test-1",
		PublicIP: "1.1.1.1",
		Metadata: nil,
		Version:  "",
	}

	testMachines := []machine.MachineState{
		expected,
	}

	testWriter := bytes.Buffer{}

	c := &FleetClient{Fleet: &stubFleetClient{testMachineStates: testMachines},
		runner: mockCommandRunner{}, errWriter: &testWriter}

	code := c.runCommand("true", "test-1")

	err := testWriter.String()
	if err != "" || code != 0 {
		t.Errorf("Error: %v, Returned %d", err, code)
	}
}

func TestFindUnit(t *testing.T) {
	t.Parallel()

	expectedID := "testing"

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name:         "deis-router@1.service",
			CurrentState: "loaded",
			MachineID:    expectedID,
		},
		&schema.Unit{
			Name: "deis-router@2.service",
		},
	}
	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits, unitsMutex: &sync.Mutex{}}}

	expectedErr := "Error retrieving Unit foo: Unit foo not found"

	_, err := c.findUnit("foo")

	actualErr := err.Error()

	if actualErr != expectedErr {
		t.Errorf("Expected '%s', Got '%s'", expectedErr, actualErr)
	}

	expectedErr = "Unit deis-router@2.service does not appear to be running.\n"

	_, err = c.findUnit("deis-router@2.service")

	actualErr = err.Error()

	if actualErr != expectedErr {
		t.Errorf("Expected '%s', Got '%s'", expectedErr, actualErr)
	}

	result, err := c.findUnit("deis-router@1.service")

	if err != nil {
		t.Error(err)
	}

	if result != expectedID {
		t.Errorf("Expected %s, Got %s", expectedID, result)
	}
}

func TestMachineState(t *testing.T) {
	t.Parallel()

	expected := machine.MachineState{
		ID:       "test-1",
		PublicIP: "1.1.1.1",
		Metadata: nil,
		Version:  "",
	}

	testMachines := []machine.MachineState{
		expected,
		machine.MachineState{
			ID:       "test-2",
			PublicIP: "2.2.2.2",
			Metadata: nil,
			Version:  "",
		},
	}
	c := &FleetClient{Fleet: &stubFleetClient{testMachineStates: testMachines}}

	result, err := c.machineState(expected.ID)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf("Expected %v, Got %v", expected, *result)
	}
}
