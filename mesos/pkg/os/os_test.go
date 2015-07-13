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

	command, args = BuildCommandFromString("docker -d -D")
	if command != "docker" {
		t.Fatalf("Expected 'docker' as value of empty env name but %s returned", command)
	}

	if len(args) != 2 {
		t.Fatalf("Expected '%v' arguments but %v returned", 0, len(args))
	}

	command, args = BuildCommandFromString("ls -lat")
	if command != "ls" {
		t.Fatalf("Expected 'ls' as value of empty env name but %s returned", command)
	}

	if len(args) != 1 {
		t.Fatalf("Expected '%v' arguments but %v returned", 1, len(args))
	}
}

func TestRandom(t *testing.T) {
	rnd, err := Random(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(rnd) != 1 {
		t.Fatalf("Expected a string of 1 character but %s returned", rnd)
	}
}

func TestRandomError(t *testing.T) {
	rnd, err := Random(0)
	if err == nil {
		t.Fatalf("Expected an error requiring a random string of length 0 but %s returned", rnd)
	}
}
