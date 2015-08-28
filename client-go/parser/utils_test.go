package parser

import "testing"

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

func TestPrintHelp(t *testing.T) {
	t.Parallel()

	usage := ""

	if !printHelp([]string{"ps", "--help"}, usage) {
		t.Error("Expected true")
	}

	if !printHelp([]string{"ps", "-h"}, usage) {
		t.Error("Expected true")
	}

	if printHelp([]string{"ps"}, usage) {
		t.Error("Expected false")
	}

	if printHelp([]string{"ps", "--foo"}, usage) {
		t.Error("Expected false")
	}
}
