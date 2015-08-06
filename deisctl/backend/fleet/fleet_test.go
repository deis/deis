package fleet

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/deis/deis/deisctl/config/model"
	"github.com/deis/deis/deisctl/test/mock"

	"github.com/coreos/fleet/machine"
	"github.com/coreos/fleet/schema"
)

type stubFleetClient struct {
	testUnits         []*schema.Unit
	testUnitStates    []*schema.UnitState
	testMachineStates []machine.MachineState
	unitStatesMutex   *sync.Mutex
	unitsMutex        *sync.Mutex
}

func (c *stubFleetClient) Machines() ([]machine.MachineState, error) {
	return c.testMachineStates, nil
}
func (c *stubFleetClient) Unit(name string) (*schema.Unit, error) {
	c.unitsMutex.Lock()
	defer c.unitsMutex.Unlock()

	for _, unit := range c.testUnits {
		if unit.Name == name {
			return unit, nil
		}
	}

	return nil, fmt.Errorf("Unit %s not found", name)
}
func (c *stubFleetClient) Units() ([]*schema.Unit, error) {
	c.unitsMutex.Lock()
	defer c.unitsMutex.Unlock()

	return c.testUnits, nil
}
func (c *stubFleetClient) UnitStates() ([]*schema.UnitState, error) {
	c.unitStatesMutex.Lock()
	defer c.unitStatesMutex.Unlock()

	return c.testUnitStates, nil
}

type failingFleetClient struct {
	stubFleetClient
}

func (c *failingFleetClient) SetUnitTargetState(name, target string) error {
	if err := c.stubFleetClient.SetUnitTargetState(name, target); err != nil {
		return err
	}

	last := len(c.testUnitStates) - 1
	c.testUnitStates[last] = &schema.UnitState{
		Name:               name,
		SystemdSubState:    "failed",
		SystemdActiveState: "failed",
	}

	return nil
}

func (c *stubFleetClient) SetUnitTargetState(name, target string) error {

	var activeState string
	var subState string

	switch target {
	case "loaded":
		activeState = "inactive"
		subState = "dead"
	case "launched":
		activeState = "active"
		subState = "running"
	}

	unit, err := c.Unit(name)

	if err != nil {
		return err
	}

	c.unitsMutex.Lock()
	unit.DesiredState = target
	c.unitsMutex.Unlock()

	c.unitStatesMutex.Lock()
	defer c.unitStatesMutex.Unlock()

	for _, unitState := range c.testUnitStates {
		if name == unitState.Name {
			unitState.SystemdSubState = subState
			unitState.SystemdActiveState = activeState
			return nil
		}
	}

	c.testUnitStates = append(c.testUnitStates, &schema.UnitState{Name: name, SystemdSubState: subState, SystemdActiveState: activeState})

	return nil
}

func (c *stubFleetClient) CreateUnit(unit *schema.Unit) error {
	c.unitsMutex.Lock()
	c.testUnits = append(c.testUnits, unit)
	c.unitsMutex.Unlock()

	return nil
}

func (c *stubFleetClient) DestroyUnit(name string) error {
	c.unitsMutex.Lock()
	for i := len(c.testUnits) - 1; i >= 0; i-- {
		if c.testUnits[i].Name == name {
			c.testUnits = append(c.testUnits[:i], c.testUnits[i+1:]...)
		}
	}
	c.unitsMutex.Unlock()

	return nil
}

func newOutErr() *outErr {
	return &outErr{
		&syncBuffer{},
		&syncBuffer{},
	}
}

// Wrap output and error streams for ease of testing.
type outErr struct {
	out, ew buffer
}

// buffer represents a buffer for collecting written test output.
//
// This is used only in testing, so add more bytes.Buffer methods as needed.
type buffer interface {
	Bytes() []byte
	String() string
	Write([]byte) (int, error)
}

// syncBuffer simply synchronizes writes on a bytes.Buffer.
type syncBuffer struct {
	bytes.Buffer
	mx sync.RWMutex
}

func (s *syncBuffer) Write(b []byte) (int, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.Buffer.Write(b)
}

func (s *syncBuffer) String() string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Buffer.String()
}

func (s *syncBuffer) Bytes() []byte {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.Buffer.Bytes()
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	// set required flags
	Flags.Endpoint = "http://127.0.0.1:4001"

	testConfigBackend := mock.ConfigBackend{Expected: []*model.ConfigNode{}}

	// instantiate client
	_, err := NewClient(testConfigBackend)
	if err != nil {
		t.Fatal(err)
	}
}
