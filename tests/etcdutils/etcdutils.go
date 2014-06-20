package etcdutils

import (
	"github.com/coreos/go-etcd/etcd"
	"testing"
	//"fmt"
	"strings"
)

type EtcdHandle struct {
	Dirs []string
	Keys []string
	C    *etcd.Client
}

func getetcdClient() *etcd.Client {
	machines := []string{"http://172.17.8.100:4001"}
	c := etcd.NewClient(machines)
	return c
}

func InitetcdValues(setdir, setkeys []string) *EtcdHandle {
	cli := getetcdClient()
	controllerHandle := new(EtcdHandle)
	controllerHandle.Dirs = setdir
	controllerHandle.Keys = setkeys
	controllerHandle.C = cli
	return controllerHandle
}

func SetEtcdValues(t *testing.T, keys []string, values []string, c *etcd.Client) {
	for i, key := range keys {
		_, err := c.Set(key, values[i], 0)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func PublishControllervalues(t *testing.T, ecli *EtcdHandle) {
	for _, dir := range ecli.Dirs {
		_, err := ecli.C.SetDir(dir, 0)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, key := range ecli.Keys {
		switch true {
		case (strings.Contains(key, "host")):
			_, err := ecli.C.Set(key, "172.17.8.100", 0)
			if err != nil {
				t.Fatal(err)
			}
		case (strings.Contains(key, "port")):
			_, err := ecli.C.Set(key, "10881", 0)
			if err != nil {
				t.Fatal(err)
			}
		default:
			_, err := ecli.C.Set(key, "deis", 0)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

/*
func CleanEtcdValues(t *testing.T, ecli EtcdHandle){
  for _, dir := range  ecli.dirs{
      _, err :=ecli.c.DeleteDir(dir)
        if err != nil {
          t.Fatal(err)
        }
    }
    for _, key := range  ecli.keys{
        _, err :=ecli.c.Delete(key,true)
          if err != nil {
            t.Fatal(err)
          }
    }
}
*/
