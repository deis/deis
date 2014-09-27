package lock

import (
	"encoding/json"

	etcdError "github.com/coreos/locksmith/third_party/github.com/coreos/etcd/error"
	"github.com/coreos/locksmith/third_party/github.com/coreos/go-etcd/etcd"
)

const (
	keyPrefix       = "coreos.com/updateengine/rebootlock"
	holdersPrefix   = keyPrefix + "/holders"
	SemaphorePrefix = keyPrefix + "/semaphore"
)

// EtcdLockClient is a wrapper around the go-etcd client that provides
// simple primitives to operate on the internal semaphore and holders
// structs through etcd.
type EtcdLockClient struct {
	client *etcd.Client
}

func NewEtcdLockClient(machines []string) (client *EtcdLockClient, err error) {
	ec := etcd.NewClient(machines)
	client = &EtcdLockClient{ec}
	err = client.Init()

	return client, err
}

// Init sets an initial copy of the semaphore if it doesn't exist yet.
func (c *EtcdLockClient) Init() (err error) {
	sem := newSemaphore()
	b, err := json.Marshal(sem)
	if err != nil {
		return err
	}

	_, err = c.client.Create(SemaphorePrefix, string(b), 0)
	if err != nil {
		eerr, ok := err.(*etcd.EtcdError)
		if ok && eerr.ErrorCode == etcdError.EcodeNodeExist {
			return nil
		}
	}

	return err
}

// Get fetches the Semaphore from etcd.
func (c *EtcdLockClient) Get() (sem *Semaphore, err error) {
	resp, err := c.client.Get(SemaphorePrefix, false, false)
	if err != nil {
		return nil, err
	}

	sem = &Semaphore{}
	err = json.Unmarshal([]byte(resp.Node.Value), sem)
	if err != nil {
		return nil, err
	}

	sem.Index = resp.Node.ModifiedIndex

	return sem, nil
}

// Set sets a Semaphore in etcd.
func (c *EtcdLockClient) Set(sem *Semaphore) (err error) {
	b, err := json.Marshal(sem)
	if err != nil {
		return err
	}

	_, err = c.client.CompareAndSwap(SemaphorePrefix, string(b), 0, "", sem.Index)

	return err
}
