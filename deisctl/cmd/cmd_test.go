package cmd

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

type stubBackend struct{}

var (
	startedUnits   []string
	stoppedUnits   []string
	installedUnits []string
)

func (backend stubBackend) Create([]string, *sync.WaitGroup, chan string, chan error) {
	return
}
func (backend stubBackend) Destroy([]string, *sync.WaitGroup, chan string, chan error) {
	return
}
func (backend stubBackend) Start(targets []string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	startedUnits = targets
	return
}
func (backend stubBackend) Stop(targets []string, wg *sync.WaitGroup, outchan chan string, errchan chan error) {
	stoppedUnits = targets
	return
}
func (backend stubBackend) Scale(string, int, *sync.WaitGroup, chan string, chan error) {
	return
}
func (backend stubBackend) ListUnits() error {
	return fmt.Errorf("ListUnits not implemented yet.")
}
func (backend stubBackend) ListUnitFiles() error {
	return fmt.Errorf("ListUnitFiles not implemented yet.")
}
func (backend stubBackend) Status(string) error {
	return fmt.Errorf("Status not implemented yet.")
}
func (backend stubBackend) Journal(string) error {
	return fmt.Errorf("Journal not implemented yet.")
}
func (backend stubBackend) SSH(string) error {
	return fmt.Errorf("SSH not implemented yet.")
}

var b stubBackend

// Start units
func TestStart(t *testing.T) {
	Start([]string{"start", "router@1", "router@2"}, b)

	if !reflect.DeepEqual(startedUnits, []string{"router@1", "router@2"}) {
		t.Error(startedUnits)
	}
}

// Stop units
func TestStop(t *testing.T) {
	Stop([]string{"stop", "router@1", "router@2"}, b)

	if !reflect.DeepEqual(stoppedUnits, []string{"router@1", "router@2"}) {
		t.Error(stoppedUnits)
	}
}

// Restart units
func TestRestart(t *testing.T) {
	Restart([]string{"restart", "router@4", "router@5"}, b)

	if !reflect.DeepEqual(stoppedUnits, []string{"router@4", "router@5"}) {
		t.Error(stoppedUnits)
	}
	if !reflect.DeepEqual(startedUnits, []string{"router@4", "router@5"}) {
		t.Error(startedUnits)
	}
}
