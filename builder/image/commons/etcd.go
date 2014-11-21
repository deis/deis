package commons

import (
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/deis/deis/boot/logger"
)

func SetDefaultEtcd(client *etcd.Client, key, value string) {
	_, err := client.Set(key, value, 0)
	if err != nil {
		logger.Log.Error(err)
	}
}

func MkdirEtcd(client *etcd.Client, path string) {
	_, err := client.CreateDir(path, 0)
	if err != nil {
		logger.Log.Error(err)
	}
}

// wait for the required keys for the
// backup/restore of ssh configuration
func WaitForEtcdKeys(client *etcd.Client, keys []string) {
	wait := true

	for {
		for _, key := range keys {
			_, err := client.Get(key, false, false)
			if err != nil {
				logger.Log.Debugf("key \"%s\" error %v", key, err)
				wait = true
			}
		}

		if !wait {
			break
		}

		logger.Log.Info("waiting for missing etcd keys...")
		time.Sleep(1 * time.Second)
		wait = false
	}
}

func GetEtcd(client *etcd.Client, key string) string {
	result, err := client.Get(key, false, false)
	if err != nil {
		return ""
	}

	return result.Node.Value
}

func GetListEtcd(client *etcd.Client, key string) etcd.Nodes {
	result, err := client.Get(key, true, false)
	if err != nil {
		return nil
	}

	return result.Node.Nodes
}

func SetEtcd(client *etcd.Client, key, value string, ttl uint64) {
	_, err := client.Set(key, value, ttl)
	if err != nil {
		logger.Log.Debugf("%v", err)
	}
}

// Publish a service to etcd periodcally
func PublishService(
	client *etcd.Client,
	host string,
	etcdPath string,
	externalPort string,
	ttl uint64,
	timeout time.Duration) {

	for {
		SetEtcd(client, etcdPath+"/host", host, ttl)
		SetEtcd(client, etcdPath+"/port", externalPort, ttl)
		time.Sleep(timeout)
	}
}
