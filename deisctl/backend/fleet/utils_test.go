package fleet

import (
	"reflect"
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
)

func TestNextComponent(t *testing.T) {
	t.Parallel()

	// test first component
	num, err := nextUnitNum([]string{})
	if err != nil {
		t.Fatal(err)
	}
	expected := 1
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	// test next component
	num, err = nextUnitNum([]string{"deis-router@1.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 2
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	// test last component
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@2.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 3
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	// test middle component
	num, err = nextUnitNum([]string{"deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 1
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 2
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@2.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 4
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
}

func TestLastComponent(t *testing.T) {
	t.Parallel()

	num, err := lastUnitNum([]string{})
	errorf := err.Error()
	expectedErr := "Component not found"
	if errorf != expectedErr {
		t.Fatalf("Expected %d, Got %d", expectedErr, errorf)
	}
	num, err = lastUnitNum([]string{"deis-router@1.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected := 1
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	num, err = lastUnitNum([]string{"deis-router@1.service", "deis-router@2.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 2
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	// test middle component
	num, err = lastUnitNum([]string{"deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 3
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	num, err = lastUnitNum([]string{"deis-router@1.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 3
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
	num, err = lastUnitNum([]string{"deis-router@1.service", "deis-router@2.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	expected = 3
	if num != expected {
		t.Fatalf("Expected %d, Got %d", expected, num)
	}
}

func TestSplitJobName(t *testing.T) {
	t.Parallel()

	c, num, err := splitJobName("deis-router@1.service")
	if err != nil {
		t.Fatal(err)
	}
	if c != "router" || num != 1 {
		t.Fatalf("Invalid values: %v %v", c, num)
	}
}

func TestSplitTarget(t *testing.T) {
	t.Parallel()

	c, num, err := splitTarget("router")
	if err != nil {
		t.Fatal(err)
	}
	if c != "router" && num != 0 {
		t.Fatalf("Invalid split on \"%v\": %v %v", "router", c, num)
	}

	c, num, err = splitTarget("router@3")
	if err != nil {
		t.Fatal(err)
	}
	if c != "router" || num != 3 {
		t.Fatalf("Invalid split on \"%v\": %v %v", "router@3", c, num)
	}

	c, num, err = splitTarget("store-gateway")
	if err != nil {
		t.Fatal(err)
	}
	if c != "store-gateway" || num != 0 {
		t.Fatalf("Invalid split on \"%v\": %v %v", "store-gateway", c, num)
	}

}

func TestExpandTargets(t *testing.T) {
	t.Parallel()

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name: "deis-router@1.service",
		},
		&schema.Unit{
			Name: "deis-router@2.service",
		},
		&schema.Unit{
			Name: "deis-store-gateway@1.service",
		},
		&schema.Unit{
			Name: "deis-store-gateway@2.service",
		},
		&schema.Unit{
			Name: "deis-controller.service",
		},
		&schema.Unit{
			Name: "deis-registry@1.service",
		},
		&schema.Unit{
			Name: "registry_v2.cmd.1.service",
		},
	}
	c := &FleetClient{Fleet: &stubFleetClient{testUnits: testUnits, unitsMutex: &sync.Mutex{}}}

	targets := []string{"router@*", "deis-store-gateway@1", "deis-controller", "registry@*"}
	expandedTargets, err := c.expandTargets(targets)
	if err != nil {
		t.Fatal(err)
	}
	expectedTargets := []string{
		"deis-router@1.service",
		"deis-router@2.service",
		"deis-store-gateway@1",
		"deis-controller",
		"deis-registry@1.service",
	}
	if !reflect.DeepEqual(expandedTargets, expectedTargets) {
		t.Fatal(expandedTargets)
	}
}
