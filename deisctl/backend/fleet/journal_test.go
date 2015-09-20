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
		return -1, fmt.Errorf("Got %s %s %d, which is unexpected", cmd, addr, timeout)
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

var ITERATIONS int
var done chan bool

func testBackgroundRunner(cmd string, IP string) (string, error) {

	if IP != "1.1.1.1" {
		if done != nil {
			close(done)
		}
		return "", fmt.Errorf("Bad IP: %s", IP)
	}

	if cmd != `journalctl --unit deis-router@1.service --no-pager -o cat --since="2006-01-02 15:04:00"` {
		if done != nil {
			close(done)
		}
		return "", fmt.Errorf("Bad Command: %s", cmd)
	}

	defer func() { ITERATIONS++ }()

	switch ITERATIONS {
	case 0:
		return "Failed to get data: Cannot assign requested address\n", nil
	case 1:
		return "test\n", nil
	case 2:
		return "foo\n", nil
	case 3:
		defer func() {
			if done != nil {
				close(done)
			}
		}()
		return "bar\n", nil
	default:
		return "", nil
	}
}

func timeStub() time.Time {
	timeStub, _ := time.Parse("2006-01-02 15:04:05 UTC", "2006-01-02 15:04:01 UTC")
	return timeStub
}

func TestJournalBackground(t *testing.T) {
	t.Parallel()

	done = make(chan bool)

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
		unitsMutex: &sync.Mutex{}}, errWriter: &testWriter, journalInterval: time.Nanosecond,
		runRemoteCommandString: testBackgroundRunner, timeNow: timeStub}

	expected := "test\nfoo\nbar\n"
	quitChans := c.JournalBackground("router@1", &testWriter)
	if done != nil {
		<-done
	}

	if len(quitChans) != 1 {
		t.Errorf("Expected 1 quit channel, Got %d", len(quitChans))
	}

	for _, quit := range quitChans {
		quit <- true
		<-quit
	}

	actual := testWriter.String()

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}
