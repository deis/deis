package parser

import (
	"reflect"
	"testing"
)

func TestCommandCombing(t *testing.T) {
	t.Parallel()

	expected := []string{"apps:create", "test", "foo"}
	actual := combineCommand([]string{"apps", "create", "test", "foo"})

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestSafeGet(t *testing.T) {
	t.Parallel()

	expected := "foo"

	test := make(map[string]interface{}, 1)
	test["test"] = "foo"

	actual := safeGetValue(test, "test")

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestSafeGetNil(t *testing.T) {
	t.Parallel()

	expected := ""

	test := make(map[string]interface{}, 1)
	test["test"] = nil

	actual := safeGetValue(test, "test")

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}
