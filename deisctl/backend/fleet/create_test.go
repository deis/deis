package fleet

import (
	"io/ioutil"
	"path"
	"sync"
	"testing"

	"github.com/deis/deis/deisctl/config/model"
	"github.com/deis/deis/deisctl/test/mock"

	"github.com/coreos/fleet/schema"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	unitFiles := []string{"deis-controller.service", "deis-builder.service",
		"deis-router.service"}

	unitContents := []byte("[Unit]")

	name, err := ioutil.TempDir("", "deisctl-fleetctl")

	if err != nil {
		t.Fatal(err)
	}

	for _, unit := range unitFiles {
		ioutil.WriteFile(path.Join(name, unit), unitContents, 777)
	}

	testFleetClient := stubFleetClient{testUnits: []*schema.Unit{}, unitsMutex: &sync.Mutex{},
		unitStatesMutex: &sync.Mutex{}}

	testConfigBackend := mock.ConfigBackend{Expected: []*model.ConfigNode{{Key: "/deis/platform/enablePlacementOptions", Value: "true"}}}

	c := &FleetClient{templatePaths: []string{name}, Fleet: &testFleetClient, configBackend: testConfigBackend}

	var errOutput string
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	se := newOutErr()
	c.Create([]string{"controller", "builder", "router@1"}, &wg, se.out, se.ew)

	wg.Wait()

	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	expectedUnits := []string{"deis-controller.service", "deis-builder.service",
		"deis-router@1.service"}

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
