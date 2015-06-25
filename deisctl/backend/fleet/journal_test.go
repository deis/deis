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

type mockJournalCommandRunner struct {
	validUnits []string
}

func (mockJournalCommandRunner) LocalCommand(string) (int, error) {
	return 0, nil
}

func (m mockJournalCommandRunner) RemoteCommand(cmd string, addr string, timeout time.Duration) (int, error) {
	if addr != "1.1.1.1" || timeout != 0 {
		return -1, fmt.Errorf("Got %s %d, which is unexpected", cmd, addr, timeout)
	}

	for _, unit := range m.validUnits {
		if fmt.Sprintf("journalctl --unit %s --no-pager -n 40 -f", unit) == cmd {
			return 0, nil
		}
	}

	return -1, fmt.Errorf("Didn't find command %s to match with units %v", cmd, m.validUnits)
}

func TestJournal(t *testing.T) {
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
	}

	testWriter := bytes.Buffer{}

	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits, testMachineStates: testMachines,
		unitsMutex: &sync.Mutex{}}, errWriter: &testWriter, runner: mockJournalCommandRunner{
		validUnits: []string{"deis-router@1.service", "deis-router@2.service"}}}

	err := c.Journal("router")

	if err != nil {
		t.Error(err)
	}

	commandErr := testWriter.String()

	if commandErr != "" {
		t.Error(commandErr)
	}
}
