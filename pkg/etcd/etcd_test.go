package etcd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/coreos/go-etcd/etcd"
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

func TestGetSetEtcd(t *testing.T) {
	startEtcd()
	defer stopEtcd()

	etcdClient := etcd.NewClient([]string{"http://localhost:4001"})
	SetDefault(etcdClient, "/path", "value")
	value := Get(etcdClient, "/path")

	if value != "value" {
		t.Fatalf("Expected '%v' arguments but returned '%v'", "value", value)
	}

	Set(etcdClient, "/path", "", 0)
	value = Get(etcdClient, "/path")

	if value != "" {
		t.Fatalf("Expected '%v' arguments but returned '%v'", "", value)
	}

	Set(etcdClient, "/path", "value", uint64((1 * time.Second).Seconds()))
	time.Sleep(2 * time.Second)
	value = Get(etcdClient, "/path")

	if value != "" {
		t.Fatalf("Expected '%v' arguments but returned '%v'", "", value)
	}
}

func TestMkdirEtcd(t *testing.T) {
	startEtcd()
	defer stopEtcd()

	etcdClient := etcd.NewClient([]string{"http://localhost:4001"})

	Mkdir(etcdClient, "/directory")
	values := GetList(etcdClient, "/directory")
	if len(values) != 0 {
		t.Fatalf("Expected '%v' arguments but returned '%v'", 0, len(values))
	}

	Set(etcdClient, "/directory/item_1", "value", 0)
	Set(etcdClient, "/directory/item_2", "value", 0)
	values = GetList(etcdClient, "/directory")
	if len(values) != 2 {
		t.Fatalf("Expected '%v' arguments but returned '%v'", 2, len(values))
	}
}

func TestWaitForKeysEtcd(t *testing.T) {
	startEtcd()
	defer stopEtcd()

	etcdClient := etcd.NewClient([]string{"http://localhost:4001"})
	Set(etcdClient, "/key", "value", 0)
	start := time.Now()
	err := WaitForKeys(etcdClient, []string{"/key"}, (10 * time.Second))
	if err != nil {
		t.Fatalf("%v", err)
	}
	end := time.Since(start)
	if end.Seconds() > (2 * time.Second).Seconds() {
		t.Fatalf("Expected '%vs' but returned '%vs'", 2, end.Seconds())
	}

	err = WaitForKeys(etcdClient, []string{"/key2"}, (2 * time.Second))
	if err == nil {
		t.Fatalf("Expected an error")
	}
}
