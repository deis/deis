package unit

import (
	"reflect"
	"testing"
)

func TestNewSystemdUnitFileFromLegacyContents(t *testing.T) {
	legacy := map[string]map[string]string{
		"Unit": {
			"Description": "foobar",
		},
		"Service": {
			"Type":      "oneshot",
			"ExecStart": "/usr/bin/echo bar",
		},
	}

	expected := map[string]map[string][]string{
		"Unit": {
			"Description": {"foobar"},
		},
		"Service": {
			"Type":      {"oneshot"},
			"ExecStart": {"/usr/bin/echo bar"},
		},
	}

	u, err := NewUnitFromLegacyContents(legacy)
	if err != nil {
		t.Fatalf("Unexpected error parsing unit %q: %v", legacy, err)
	}
	actual := u.Contents

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Map func did not produce expected output.\nActual=%v\nExpected=%v", actual, expected)
	}
}
