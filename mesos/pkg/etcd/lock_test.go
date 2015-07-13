package etcd

import (
	"os/exec"
	"testing"
)

func init() {
	_, err := exec.Command("etcd", "--version").Output()
	if err != nil {
		log.Fatal(err)
	}
}

func TestAcquireReleaseLock(t *testing.T) {
	startEtcd()
	defer stopEtcd()

	etcdClient := NewClient([]string{"http://localhost:4001"})

	err := AcquireLock(etcdClient, "/lock", 10)
	if err != nil {
		t.Fatalf("Unexpected error '%v'", err)
	}

	value := Get(etcdClient, "/lock")
	if value == "" {
		t.Fatalf("Expected '%v' arguments but returned '%v'", "locked", value)
	}

	if value != "locked" {
		t.Fatalf("Expected '%v' arguments but returned '%v'", "locked", value)
	}

	ReleaseLock(etcdClient)

	value = Get(etcdClient, "/lock")
	if value != "released" {
		t.Fatalf("Expected '%v' arguments but returned '%v'", "released", value)
	}

}
