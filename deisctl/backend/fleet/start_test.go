package fleet

import (
	"strings"
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
)

var startTestUnits = []*schema.Unit{
	&schema.Unit{
		Name:         "deis-controller.service",
		DesiredState: "loaded",
	},
	&schema.Unit{
		Name:         "deis-builder.service",
		DesiredState: "loaded",
	},
	&schema.Unit{
		Name:         "deis-publisher.service",
		DesiredState: "loaded",
	},
}

func TestStart(t *testing.T) {
	t.Parallel()

	testFleetClient := stubFleetClient{testUnits: startTestUnits,
		unitsMutex: &sync.Mutex{}, unitStatesMutex: &sync.Mutex{}}

	c := &FleetClient{Fleet: &testFleetClient}

	var errOutput string
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	se := newOutErr()
	c.Start([]string{"controller", "builder", "publisher"}, &wg, se.out, se.ew)

	wg.Wait()

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

				if unit.SystemdSubState != "running" {
					t.Errorf("Unit %s is %s, expected running", unit.Name, unit.SystemdSubState)
				}

				break
			}
		}

		if !found {
			t.Errorf("Expected Unit %s not found in Unit States", expectedUnit)
		}
	}
}

func TestStartFail(t *testing.T) {
	fc := &failingFleetClient{stubFleetClient{
		testUnits:       startTestUnits,
		unitStatesMutex: &sync.Mutex{},
		unitsMutex:      &sync.Mutex{},
	}}
	var wg sync.WaitGroup
	c := &FleetClient{Fleet: fc}

	var b syncBuffer
	c.Start([]string{"deis-builder.service"}, &wg, &b, &b)
	wg.Wait()

	if !strings.Contains(b.String(), "failed while starting") {
		t.Errorf("Expected failure during start. Got '%s'", b.String())
	}

}
