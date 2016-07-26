package api

import (
	"sort"
	"testing"
)

func TestAppsSorted(t *testing.T) {
	apps := Apps{
		{"2014-01-01T00:00:00UTC", "Zulu", "John", "2016-01-02", "zulu.example.com", "d57be2ba-7ae2-4825-9ace-7c86cb893046"},
		{"2014-01-01T00:00:00UTC", "Alpha", "John", "2016-01-02", "alpha.example.com", "3d501190-1b8e-41ef-94c5-dd9a0bb707bb"},
		{"2014-01-01T00:00:00UTC", "Gamma", "John", "2016-01-02", "gamma.example.com", "41d95133-fd4d-4f4c-92a2-e454857371cc"},
		{"2014-01-01T00:00:00UTC", "Beta", "John", "2016-01-02", "beta.example.com", "222ed1aa-e985-4bec-9966-a88215300661"},
	}

	sort.Sort(apps)
	expectedAppNames := []string{"Alpha", "Beta", "Gamma", "Zulu"}

	for i, app := range apps {
		if expectedAppNames[i] != app.ID {
			t.Errorf("Expected apps to be sorted %v, Got %v at index %v", expectedAppNames[i], app.ID, i)
		}
	}
}
