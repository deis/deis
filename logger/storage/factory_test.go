package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetUsingInvalidValues(t *testing.T) {
	adapterStrs := []string{"bogus", "memory:", "memory:foo"}
	for _, adapterStr := range adapterStrs {
		_, err := NewAdapter(adapterStr)
		if err == nil || err.Error() != fmt.Sprintf("Unrecognized storage adapter type: '%s'", adapterStr) {
			t.Error("Did not receive expected error message")
		}
	}
}

func TestGetUsingEmptyString(t *testing.T) {
	a, err := NewAdapter("")
	if err != nil {
		t.Error(err)
	}
	// Should use file-based log storage by default
	expected := "*file.adapter"
	aType := reflect.TypeOf(a).String()
	if aType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, aType)
	}
}

func TestGetFileBasedAdapter(t *testing.T) {
	a, err := NewAdapter("file")
	if err != nil {
		t.Error(err)
	}
	expected := "*file.adapter"
	aType := reflect.TypeOf(a).String()
	if aType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, aType)
	}
}

func TestGetMemoryBasedAdapter(t *testing.T) {
	a, err := NewAdapter("memory")
	if err != nil {
		t.Error(err)
	}
	expected := "*ringbuffer.adapter"
	aType := reflect.TypeOf(a).String()
	if aType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, aType)
	}
}

func TestGetMemoryBasedAdapterWithBufferSize(t *testing.T) {
	a, err := NewAdapter("memory:1000")
	if err != nil {
		t.Error(err)
	}
	expected := "*ringbuffer.adapter"
	aType := reflect.TypeOf(a).String()
	if aType != expected {
		t.Errorf("Expected a %s, but got a %s", expected, aType)
	}
}

func TestMain(m *testing.M) {
	LogRoot, _ = ioutil.TempDir("", "log-tests")
	defer os.Remove(LogRoot)
	os.Exit(m.Run())
}
