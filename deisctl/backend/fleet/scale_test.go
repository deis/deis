package fleet

import (
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/deis/deis/deisctl/config/model"
	"github.com/deis/deis/deisctl/test/mock"

	"github.com/coreos/fleet/schema"
)

func TestScaleUp(t *testing.T) {
	t.Parallel()

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

	testConfigBackend := mock.ConfigBackend{Expected: []*model.ConfigNode{{Key: "/deis/platform/enablePlacementOptions", Value: "true"}}}

	c := &FleetClient{templatePaths: []string{name}, Fleet: &testFleetClient, configBackend: testConfigBackend}

	var errOutput string
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	se := newOutErr()

	c.Scale("router", 3, &wg, se.out, se.ew)

	wg.Wait()

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
			t.Errorf("Expected Unit %s not found in Unit States", expectedUnit)
		}
	}
}

func TestScaleDown(t *testing.T) {
	t.Parallel()

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
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	se := newOutErr()
	c.Scale("router", 1, &wg, se.out, se.ew)

	wg.Wait()

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
			t.Errorf("Expected Unit %s not found in Unit States", expectedUnit)
		}
	}
}

func TestScaleError(t *testing.T) {
	t.Parallel()

	c := &FleetClient{Fleet: &stubFleetClient{}}

	var errOutput string
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	se := newOutErr()
	c.Scale("router", -1, &wg, se.out, se.ew)

	wg.Wait()

	expected := "cannot scale below 0"
	errOutput = strings.TrimSpace(se.ew.String())

	logMutex.Lock()
	if errOutput != expected {
		t.Errorf("Expected '%s', Got '%s'", expected, errOutput)
	}
	logMutex.Unlock()
}
