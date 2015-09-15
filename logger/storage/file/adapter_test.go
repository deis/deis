package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

const app string = "test-app"

func TestReadFromNonExistingApp(t *testing.T) {
	logRoot, err := ioutil.TempDir("", "log-tests")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(logRoot)
	// Initialize a new storage adapter
	a, err := NewStorageAdapter(logRoot)
	if err != nil {
		t.Error(err)
	}
	// No logs have been writter; there should be no ringBuffer for app
	messages, err := a.Read(app, 10)
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
	a, err := NewStorageAdapter(bogusLogRoot)
	if a != nil {
		t.Error("Expected no storage adapter, but got one")
	}
	if err == nil || err.Error() != fmt.Sprintf("Directory %s does not exist", bogusLogRoot) {
		t.Error("Did not receive expected error message")
	}
	// Create a temporary file
	file, err := ioutil.TempFile("", "log-file")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())
	a, err = NewStorageAdapter(file.Name())
	if a != nil {
		t.Error("Expected no storage adapter, but got one")
	}
	if err == nil || err.Error() != fmt.Sprintf("%s is not a directory", file.Name()) {
		t.Error("Did not receive expected error message")
	}
}

func TestLogs(t *testing.T) {
	logRoot, err := ioutil.TempDir("", "log-tests")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(logRoot)
	a, err := NewStorageAdapter(logRoot)
	if err != nil {
		t.Error(err)
	}
	// And write a few logs
	for i := 0; i < 5; i++ {
		if err := a.Write(app, fmt.Sprintf("message %d", i)); err != nil {
			t.Error(err)
		}
	}
	// Read more logs than there are
	messages, err := a.Read(app, 8)
	if err != nil {
		t.Error(err)
	}
	// Should only get as many messages as we actually have
	if len(messages) != 5 {
		t.Error("only expected 5 log messages")
	}
	// Read fewer logs than there are
	messages, err = a.Read(app, 3)
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

func TestDestroy(t *testing.T) {
	logRoot, err := ioutil.TempDir("", "log-tests")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(logRoot)
	a, err := NewStorageAdapter(logRoot)
	if err != nil {
		t.Error(err)
	}
	// Write a log to create the file
	if err := a.Write(app, "Hello, log!"); err != nil {
		t.Error(err)
	}
	filename := path.Join(logRoot, fmt.Sprintf("%s.log", app))
	// Test log file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Log file was expected to exist, but doesn't.")
	}
	// Now destroy it
	if err := a.Destroy(app); err != nil {
		t.Error(err)
	}
	// Now check that the file no longer exists
	if _, err := os.Stat(filename); err == nil {
		t.Error("Log file still exists, but was expected not to.")
	}
	// And that we have no reference to it
	if _, ok := a.files[app]; ok {
		t.Error("Log fiel reference still exist, but was expected not to.")
	}
}

func TestReopen(t *testing.T) {
	logRoot, err := ioutil.TempDir("", "log-tests")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(logRoot)
	a, err := NewStorageAdapter(logRoot)
	if err != nil {
		t.Error(err)
	}
	// Write a log to create the file
	if err := a.Write(app, "Hello, log!"); err != nil {
		t.Error(err)
	}
	// At least one file reference should exist
	if len(a.files) == 0 {
		t.Error("At least one log file reference expected to exist, but doesn't.")
	}
	// Now "reopen" logs
	if err := a.Reopen(); err != nil {
		t.Error(err)
	}
	// Now check that no file references exist
	if len(a.files) != 0 {
		t.Error("At least one log file reference still exists, but was expected not to.")
	}
}
