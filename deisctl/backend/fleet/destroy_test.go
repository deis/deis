package fleet

import (
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
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	oe := newOutErr()
	c.Destroy([]string{"controller", "registry", "router@1"}, &wg, oe.out, oe.ew)

	wg.Wait()

	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	if len(testFleetClient.testUnits) != 1 || testFleetClient.testUnits[0].Name != "deis-builder.service" {
		t.Errorf("Got %d Units (want 1), first unit %s (want builder)", len(testFleetClient.testUnits), testFleetClient.testUnits[0].Name)
	}
}
