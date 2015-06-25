package utils

import (
	"os"
	"testing"
)

// TestResolvePath ensures ResolvePath replaces $HOME and ~ as expected
func TestResolvePath(t *testing.T) {
	t.Parallel()

	paths := []string{"$HOME/foo/bar", "~/foo/bar"}
	expected := os.Getenv("HOME") + "/foo/bar"
	for _, path := range paths {
		resolved := ResolvePath(path)
		if resolved != expected {
			t.Errorf("Test failed: expected %s, got %s", expected, resolved)
		}
	}
}
