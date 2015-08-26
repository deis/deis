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
		actual, argv := parseArgs([]string{test})

		if actual != expected {
			t.Errorf("Expected %s, Got %s", expected, actual)
		}

		if len(argv) != 1 {
			t.Errorf("Expected length of 1, Got %d", len(argv))
		}
	}
}

func TestHelpReformattingWithCommand(t *testing.T) {
	t.Parallel()

	tests := []string{"--help", "-h", "help"}
	expected := "test"
	expectedArgv := []string{"test", "--help"}

	for _, test := range tests {
		actual, argv := parseArgs([]string{test, "test"})

		if actual != expected {
			t.Errorf("Expected %s, Got %s", expected, actual)
		}

		if !reflect.DeepEqual(expectedArgv, argv) {
			t.Errorf("Expected %v, Got %v", expectedArgv, argv)
		}
	}
}

func TestCommandSplitting(t *testing.T) {
	t.Parallel()

	expected := "apps"
	expectedArgv := []string{"apps:create", "test", "foo"}

	actual, argv := parseArgs(expectedArgv)

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}

	if !reflect.DeepEqual(expectedArgv, argv) {
		t.Errorf("Expected %v, Got %v", expectedArgv, argv)
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
