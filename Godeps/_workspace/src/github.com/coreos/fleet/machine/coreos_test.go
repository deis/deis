package machine

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestReadLocalMachineIDMissing(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "fleet-")
	if err != nil {
		t.Fatalf("Failed creating tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	machID, err := readLocalMachineID(dir)
	if err == nil {
		t.Fatal("Expected error for missing machID, but got nil")
	}
	if machID != "" {
		t.Fatalf("Received incorrect machID: %s", machID)
	}
}

func TestReadLocalMachineIDFound(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "fleet-")
	if err != nil {
		t.Fatalf("Failed creating tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	tmpMachineIDPath := filepath.Join(dir, "/etc/machine-id")
	err = os.MkdirAll(filepath.Dir(tmpMachineIDPath), os.FileMode(0755))
	if err != nil {
		t.Fatalf("Failed setting up fake mach ID path: %v", err)
	}

	err = ioutil.WriteFile(tmpMachineIDPath, []byte("pingpong"), os.FileMode(0644))
	if err != nil {
		t.Fatalf("Failed writing fake mach ID file: %v", err)
	}

	machID, err := readLocalMachineID(dir)
	if err != nil {
		t.Fatalf("Unexpected error reading machID: %v", err)
	}
	if machID != "pingpong" {
		t.Fatalf("Received incorrect machID %q, expected 'pingpong'", machID)
	}
}
