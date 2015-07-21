package fleet

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"
)

func init() {
	_, err := exec.Command("etcd", "--version").Output()
	if err != nil {
		log.Fatal(err)
	}
}

var etcdServer *exec.Cmd

func startEtcd() {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "etcd-test")
	if err != nil {
		log.Fatal("creating temp dir:", err)
	}
	log.Debugf("temp dir: %v", tmpDir)

	etcdServer = exec.Command("etcd", "-data-dir="+tmpDir, "-name=default")
	etcdServer.Start()
	time.Sleep(1 * time.Second)
}

func stopEtcd() {
	etcdServer.Process.Kill()
}

func TestGetNodesWithMetadata(t *testing.T) {
	startEtcd()
	defer stopEtcd()

	data, err := ParseMetadata("zookeeper=true")
	if err != nil {
		t.Fatalf("Unexpected error '%v'", err)
	}

	machines, err := GetNodesWithMetadata([]string{"http://172.17.8.100:4001"}, data)
	if err != nil {
		t.Fatalf("Expected '%v' arguments but returned '%v'", "", err)
	}

	if len(machines) <= 0 {
		t.Fatalf("Expected at least one machines but %v were returned ", len(machines))
	}
}
