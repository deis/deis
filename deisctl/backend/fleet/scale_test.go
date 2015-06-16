package fleet

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
)

func TestScaleUp(t *testing.T) {
	name, err := ioutil.TempDir("", "deisctl-fleetctl")

	if err != nil {
		t.Fatal(err)
	}

	ioutil.WriteFile(path.Join(name, "deis-router.service"), []byte("[Unit]"), 777)

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name:         "deis-router@1.service",
			DesiredState: "launched",
		},
	}

	testFleetClient := stubFleetClient{testUnits: testUnits, testUnitStates: []*schema.UnitState{}, unitsMutex: &sync.Mutex{}, unitStatesMutex: &sync.Mutex{}}

	c := &FleetClient{templatePaths: []string{name}, Fleet: &testFleetClient}

	var errOutput string
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	go logState(outchan, errchan, &errOutput, &logMutex)

	c.Scale("router", 3, &wg, outchan, errchan)

	wg.Wait()
	close(errchan)
	close(outchan)

	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	expectedUnits := []string{"deis-router@1.service", "deis-router@2.service",
		"deis-router@3.service"}

	for _, expectedUnit := range expectedUnits {
		found := false

		for _, unit := range testFleetClient.testUnits {
			if unit.Name == expectedUnit {
				found = true
				break
			}
		}

		if !found {
			t.Error(fmt.Errorf("Expected Unit %s not found in Unit States", expectedUnit))
		}
	}
}

func TestScaleDown(t *testing.T) {
	testUnits := []*schema.Unit{
		&schema.Unit{
			Name:         "deis-router@1.service",
			DesiredState: "launched",
		},
		&schema.Unit{
			Name:         "deis-router@2.service",
			DesiredState: "launched",
		},
		&schema.Unit{
			Name:         "deis-router@3.service",
			DesiredState: "launched",
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

	c.Scale("router", 1, &wg, outchan, errchan)

	wg.Wait()
	close(errchan)
	close(outchan)

	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	expectedUnits := []string{"deis-router@1.service"}

	for _, expectedUnit := range expectedUnits {
		found := false

		for _, unit := range testFleetClient.testUnits {
			if unit.Name == expectedUnit {
				found = true
				break
			}
		}

		if !found {
			t.Error(fmt.Errorf("Expected Unit %s not found in Unit States", expectedUnit))
		}
	}
}

func TestScaleError(t *testing.T) {
	c := &FleetClient{Fleet: &stubFleetClient{}}

	var errOutput string
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	go logState(outchan, errchan, &errOutput, &logMutex)

	c.Scale("router", -1, &wg, outchan, errchan)

	wg.Wait()
	close(errchan)
	close(outchan)

	expected := "cannot scale below 0\n"

	logMutex.Lock()
	if errOutput != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, errOutput))
	}
	logMutex.Unlock()
}
