package api

import (
	"sort"
	"testing"
)

func TestKeysSorted(t *testing.T) {
	keys := Keys{
		{"", "Delta", "", "", "", ""},
		{"", "Alpha", "", "", "", ""},
		{"", "Gamma", "", "", "", ""},
		{"", "Zeta", "", "", "", ""},
	}

	sort.Sort(keys)
	expectedKeys := []string{"Alpha", "Delta", "Gamma", "Zeta"}

	for i, key := range keys {
		if expectedKeys[i] != key.ID {
			t.Errorf("Expected domains to be sorted %v, Got %v at index %v", expectedKeys[i], key.ID, i)
		}
	}
}
