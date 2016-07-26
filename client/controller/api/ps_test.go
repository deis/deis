package api

import (
	"sort"
	"testing"
)

func TestProcessesSorted(t *testing.T) {
	processes := Processes{
		{"", "", "", "", "", "", "web", 4, "up"},
		{"", "", "", "", "", "", "web", 2, "up"},
		{"", "", "", "", "", "", "web", 3, "up"},
		{"", "", "", "", "", "", "web", 1, "up"},
	}

	// The API will return this sorted already, just to be sure
	sort.Sort(processes)

	for i, process := range processes {
		if i+1 != process.Num {
			t.Errorf("Expected processes to be sorted %v, Got %v", i+1, process.Num)
		}
	}
}

func TestProcessTypesSorted(t *testing.T) {
	processTypes := ProcessTypes{
		{"worker", Processes{}},
		{"web", Processes{}},
		{"clock", Processes{}},
	}

	sort.Sort(processTypes)
	expectedProcessTypes := []string{"clock", "web", "worker"}

	for i, processType := range processTypes {
		if expectedProcessTypes[i] != processType.Type {
			t.Errorf("Expected apps to be sorted %v, Got %v at index %v", expectedProcessTypes[i], processType.Type, i)
		}
	}
}
