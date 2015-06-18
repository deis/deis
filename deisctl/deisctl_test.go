package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/deis/deis/version"
)

// commandOutput returns stdout for a deisctl command line as a string.
func commandOutput(args []string) (output string) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Command(args)

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	os.Stdout = old
	output = <-outC
	return
}

// TestHelp tests that deisctl is flexible when being asked to print built-in help.
func TestHelp(t *testing.T) {
	allArgs := [][]string{{"-h"}, {"--help"}, {"help"}}
	out := ""
	for _, args := range allArgs {
		out = commandOutput(args)
		if !strings.Contains(out, "Usage: deisctl [options] <command> [<args>...]") ||
			!strings.Contains(out, "Commands, use \"deisctl help <command>\" to learn more") {
			t.Error(out)
		}
	}
}

// TestUsage ensures that deisctl prints a short usage string when no arguments were provided.
func TestUsage(t *testing.T) {
	out := commandOutput(nil)
	expected := "Usage: deisctl [options] <command> [<args>...]\n"
	if out != expected {
		t.Error(fmt.Errorf("Expected '%s', Got '%s'", expected, out))
	}
}

// TestVersion verifies that "deisctl --version" prints the current version string.
func TestVersion(t *testing.T) {
	args := []string{"--version"}
	out := commandOutput(args)
	if !strings.HasPrefix(out, version.Version) {
		t.Error(out)
	}
}
