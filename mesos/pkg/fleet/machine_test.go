package fleet

import (
	"reflect"
	"testing"
)

func TestParseMetadata(t *testing.T) {
	data, err := ParseMetadata("zookeeper=true")
	if err != nil {
		t.Fatalf("Unexpected error '%v'", err)
	}

	expected := make(map[string][]string)
	expected["zookeeper"] = append(expected["zookeeper"], "true")

	if !reflect.DeepEqual(data, expected) {
		t.Fatalf("Expected map with zookeeper=true but it returned %v", expected)
	}
}
