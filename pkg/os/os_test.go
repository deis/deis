package os

import (
	"testing"
)

func TestGetoptEmpty(t *testing.T) {
	value := Getopt("", "")
	if value != "" {
		t.Fatalf("Expected '' as value of empty env name but %s returned", value)
	}
}

func TestGetoptValid(t *testing.T) {
	value := Getopt("valid", "value")
	if value != "value" {
		t.Fatalf("Expected 'value' as value of 'valid' but %s returned", value)
	}
}

func TestGetoptDefault(t *testing.T) {
	value := Getopt("", "default")
	if value != "default" {
		t.Fatalf("Expected 'default' as value of empty env name but %s returned", value)
	}
}

func TestBuildCommandFromStringSingle(t *testing.T) {
	command, args := BuildCommandFromString("ls")
	if command != "ls" {
		t.Fatalf("Expected 'ls' as value of empty env name but %s returned", command)
	}

	if len(args) != 0 {
		t.Fatalf("Expected '%v' arguments but %v returned", 0, len(args))
	}
}

func TestBuildCommandFromStringNoArgs(t *testing.T) {
	command, args := BuildCommandFromString("ls -lat")
	if command != "ls" {
		t.Fatalf("Expected 'ls' as value of empty env name but %s returned", command)
	}

	if len(args) != 1 {
		t.Fatalf("Expected '%v' arguments but %v returned", 1, len(args))
	}
}
