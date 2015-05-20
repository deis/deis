package fleet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

func TestNextComponent(t *testing.T) {
	// test first component
	num, err := nextUnitNum([]string{})
	if err != nil {
		t.Fatal(err)
	}
	if num != 1 {
		t.Fatal("Invalid component number")
	}
	// test next component
	num, err = nextUnitNum([]string{"deis-router@1.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatal("Invalid component number")
	}
	// test last component
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@2.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 3 {
		t.Fatal("Invalid component number")
	}
	// test middle component
	num, err = nextUnitNum([]string{"deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 1 {
		t.Fatal("Invalid component number")
	}
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("Invalid component number: %v", num)
	}
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@2.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 4 {
		t.Fatal("Invalid component number")
	}
}

func TestSplitJobName(t *testing.T) {
	c, num, err := splitJobName("deis-router@1.service")
	if err != nil {
		t.Fatal(err)
	}
	if c != "router" || num != 1 {
		t.Fatalf("Invalid values: %v %v", c, num)
	}
}

func TestSplitTarget(t *testing.T) {
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

type stubFleetClient struct{}

var (
	testUnits []*schema.Unit
)

func (c stubFleetClient) Machines() ([]machine.MachineState, error) {
	return []machine.MachineState{}, nil
}
func (c stubFleetClient) Unit(string) (*schema.Unit, error) {
	return nil, nil
}
func (c stubFleetClient) Units() ([]*schema.Unit, error) {
	return testUnits, nil
}
func (c stubFleetClient) UnitStates() ([]*schema.UnitState, error) {
	return []*schema.UnitState{}, nil
}
func (c stubFleetClient) SetUnitTargetState(name, target string) error {
	return fmt.Errorf("SetUnitTargetState not implemented yet.")
}
func (c stubFleetClient) CreateUnit(*schema.Unit) error {
	return fmt.Errorf("CreateUnit not implemented yet.")
}
func (c stubFleetClient) DestroyUnit(string) error {
	return fmt.Errorf("DestroyUnit not implemented yet.")
}

var fc stubFleetClient

func TestExpandTargets(t *testing.T) {
	testUnits = []*schema.Unit{
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
	c := &FleetClient{Fleet: fc}

	targets := []string{"router@*", "deis-store-gateway@1", "deis-controller", "registry@*"}
	expandedTargets, err := expandTargets(c, targets)
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
