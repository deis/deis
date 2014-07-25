package pkg

import (
	"os"
	"strings"
	"testing"
)

var pathtests = []struct {
	in     string
	home   string
	prefix string
	suffix string
	full   string
}{
	{"~/a", "", "/", "/a", ""},
	{"~", "", "/", "", ""},
	{"~~", "", "", "", "~~"},
	{"~/a", "/home/foo", "", "", "/home/foo/a"},
}

// TestParseFilepath tests parsing filepath
func TestParseFilepath(t *testing.T) {
	for _, test := range pathtests {
		oldHome := os.Getenv("HOME")
		if err := os.Setenv("HOME", test.home); err != nil {
			t.Fatalf("Failed setting $HOME")
		}
		path := ParseFilepath(test.in)
		if !strings.HasPrefix(path, test.prefix) {
			t.Errorf("Failed parsing out prefix %v for %v", test.prefix, test.in)
		}
		if !strings.HasSuffix(path, test.suffix) {
			t.Errorf("Failed parsing out suffix %v for %v", test.suffix, test.in)
		}
		if test.full != "" && path != test.full {
			t.Errorf("Failed parsing out %v for %v", test.full, test.in)
		}
		if err := os.Setenv("HOME", oldHome); err != nil {
			t.Fatalf("Failed recovering $HOME")
		}
	}
}
