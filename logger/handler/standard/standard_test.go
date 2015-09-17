package standard

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const app string = "test-app"

func TestReadFromNonExistingApp(t *testing.T) {
	logRoot, err := ioutil.TempDir(os.TempDir(), "log-tests")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(logRoot)
	// Initialize a new handler
	h, err := NewHandler(logRoot)
	if err != nil {
		t.Error(err)
	}
	// No logs have been writter; there should be no ringBuffer for app
	messages, err := h.Read(app, 10)
	if messages != nil {
		t.Error("Expected no messages, but got some")
	}
	if err == nil || err.Error() != fmt.Sprintf("Could not find logs for '%s'", app) {
		t.Error("Did not receive expected error message")
	}
}

func TestWithBadLogRoot(t *testing.T) {
	// Initialize with bad log path (doesn't exist)
	bogusLogRoot := "/bogus/path"
	h, err := NewHandler(bogusLogRoot)
	if h != nil {
		t.Error("Expected no handler, but got one")
	}
	if err == nil || err.Error() != fmt.Sprintf("Directory %s does exist", bogusLogRoot) {
		t.Error("Did not receive expected error message")
	}
	// Create a temporary file
	file, err := ioutil.TempFile(os.TempDir(), "log-file")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())
	h, err = NewHandler(file.Name())
	if h != nil {
		t.Error("Expected no handler, but got one")
	}
	if err == nil || err.Error() != fmt.Sprintf("%s is not a directory", file.Name()) {
		t.Error("Did not receive expected error message")
	}
}

func TestLogs(t *testing.T) {
	logRoot, err := ioutil.TempDir(os.TempDir(), "log-tests")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(logRoot)
	h, err := NewHandler(logRoot)
	if err != nil {
		t.Error(err)
	}
	// And write a few logs
	for i := 0; i < 5; i++ {
		if err := h.Write(app, fmt.Sprintf("message %d", i)); err != nil {
			t.Error(err)
		}
	}
	// Read more logs than there are
	messages, err := h.Read(app, 8)
	if err != nil {
		t.Error(err)
	}
	// Should only get as many messages as we actually have
	if len(messages) != 5 {
		t.Error("only expected 5 log messages")
	}
	// Read fewer logs than there are
	messages, err = h.Read(app, 3)
	if err != nil {
		t.Error(err)
	}
	// Should get the 3 MOST RECENT logs
	if len(messages) != 3 {
		t.Error("only expected 5 log messages, got %d", len(messages))
	}
	for i := 0; i < 3; i++ {
		expectedMessage := fmt.Sprintf("message %d", i+2)
		if messages[i] != expectedMessage {
			t.Errorf("expected: \"%s\", got \"%s\"", expectedMessage, messages[i])
		}
	}
}
