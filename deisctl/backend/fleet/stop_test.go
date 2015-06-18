package fleet

import (
	"fmt"
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
)

func TestStop(t *testing.T) {
	testUnits := []*schema.Unit{
		&schema.Unit{
			Name:         "deis-controller.service",
			DesiredState: "launched",
		},
		&schema.Unit{
			Name:         "deis-builder.service",
			DesiredState: "launched",
		},
		&schema.Unit{
			Name:         "deis-publisher.service",
			DesiredState: "launch",
		},
	}

	testFleetClient := stubFleetClient{testUnits: testUnits, unitsMutex: &sync.Mutex{},
		unitStatesMutex: &sync.Mutex{}}

	c := &FleetClient{Fleet: &testFleetClient}

	var errOutput string
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	go logState(outchan, errchan, &errOutput, &logMutex)

	c.Stop([]string{"controller", "builder", "publisher"}, &wg, outchan, errchan)

	wg.Wait()
	close(errchan)
	close(outchan)

	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	expected := []string{"deis-controller.service", "deis-builder.service", "deis-publisher.service"}

	for _, expectedUnit := range expected {
		found := false

		for _, unit := range testFleetClient.testUnitStates {
			if unit.Name == expectedUnit {
				found = true

				if unit.SystemdSubState != "dead" {
					t.Error(fmt.Errorf("Unit %s is %s, expected dead", unit.Name, unit.SystemdSubState))
				}

				break
			}
		}

		if !found {
			t.Error(fmt.Errorf("Expected Unit %s not found in Unit States", expectedUnit))
		}
	}
}
