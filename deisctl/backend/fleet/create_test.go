package fleet

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
)

func TestCreate(t *testing.T) {

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

	c := &FleetClient{templatePaths: []string{name}, Fleet: &testFleetClient}

	var errOutput string
	outchan := make(chan string)
	errchan := make(chan error)
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	go logState(outchan, errchan, &errOutput, &logMutex)

	c.Create([]string{"controller", "builder", "router@1"}, &wg, outchan, errchan)

	wg.Wait()
	close(errchan)
	close(outchan)

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
			t.Error(fmt.Errorf("Expected Unit %s not found in Unit States", expectedUnit))
		}
	}
}
