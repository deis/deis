package fleet

import (
	"fmt"
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
)

func TestDestroy(t *testing.T) {
	t.Parallel()

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name: "deis-registry.service",
		},
		&schema.Unit{
			Name: "deis-builder.service",
		},
		&schema.Unit{
			Name: "deis-router@1.service",
		},
	}

	testFleetClient := stubFleetClient{testUnits: testUnits, testUnitStates: []*schema.UnitState{}, unitsMutex: &sync.Mutex{}, unitStatesMutex: &sync.Mutex{}}

	c := &FleetClient{Fleet: &testFleetClient}

	var errOutput string
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	go logState(outchan, errchan, &errOutput, &logMutex)

	c.Destroy([]string{"controller", "registry", "router@1"}, &wg, outchan, errchan)

	wg.Wait()
	close(errchan)
	close(outchan)

	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	if len(testFleetClient.testUnits) != 1 || testFleetClient.testUnits[0].Name != "deis-builder.service" {
		t.Error(fmt.Errorf("Got %d Units (want 1), first unit %s (want builder)", len(testFleetClient.testUnits), testFleetClient.testUnits[0].Name))
	}
}
