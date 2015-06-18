package main

import (
	"reflect"
	"testing"
)

func TestHelpReformatting(t *testing.T) {
	t.Parallel()

	tests := []string{"--help", "-h", "help"}
	expected := "help"

	for _, test := range tests {
		actual := parseArgs([]string{test})

		if actual[0] != expected {
			t.Errorf("Expected %s, Got %s", expected, actual[0])
		}
	}
}

func TestHelpReformattingWithCommand(t *testing.T) {
	t.Parallel()

	tests := []string{"--help", "-h", "help"}
	expected := "--help"

	for _, test := range tests {
		actual := parseArgs([]string{test, "test"})

		if actual[1] != expected {
			t.Errorf("Expected %s, Got %s", expected, actual[1])
		}
	}
}

func TestCommandSplitting(t *testing.T) {
	t.Parallel()

	expected := []string{"apps", "create", "test", "foo"}
	actual := parseArgs([]string{"apps:create", "test", "foo"})

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestReplaceShortcutRepalce(t *testing.T) {
	t.Parallel()

	expected := "apps:create"
	actual := replaceShortcut("create")

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestReplaceShortcutUnchanged(t *testing.T) {
	t.Parallel()

	expected := "users:list"
	actual := replaceShortcut(expected)

	if expected != actual {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}
