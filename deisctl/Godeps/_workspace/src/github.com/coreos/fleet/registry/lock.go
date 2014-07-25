package registry

import (
	"fmt"
	"path"
	"time"

	log "github.com/coreos/fleet/third_party/github.com/golang/glog"

	"github.com/coreos/fleet/etcd"
)

var (
	// ResourceMutexTTL is the duration of time to allow a mutex to be held on a resource
	ResourceMutexTTL = 10 * time.Second
)

const (
	mutexPrefix = "/mutex/"
)

// lockResource will attempt to lock a mutex on a resource defined by the
// provided class and id. The context will be persisted to the Registry to
// track by whom the mutex is currently locked.
func (r *EtcdRegistry) lockResource(class, id, context string) *TimedResourceMutex {
	mutexName := fmt.Sprintf("%s-%s", class, id)
	log.V(1).Infof("Attempting to acquire mutex on %s", mutexName)

	req := etcd.Create{
		Key:   path.Join(r.keyPrefix, mutexPrefix, mutexName),
		Value: context,
		TTL:   ResourceMutexTTL,
	}

	resp, err := r.etcd.Do(&req)
	if err != nil {
		log.V(1).Infof("Failed to acquire mutex on %s", mutexName)
		return nil
	}

	log.V(1).Infof("Successfully acquired mutex on %s", mutexName)
	return &TimedResourceMutex{r.etcd, *resp.Node}
}

// TimedResourceMutex is a proxy to an auto-expiring mutex
// stored in the Registry. It assumes the mutex creator has
// initialized a timer.
type TimedResourceMutex struct {
	etcd etcd.Client
	node etcd.Node
}

// Unlock will attempt to remove the lock held on the mutex
// in the Registry.
func (t *TimedResourceMutex) Unlock() error {
	req := etcd.Delete{
		Key:           t.node.Key,
		PreviousIndex: t.node.CreatedIndex,
	}
	_, err := t.etcd.Do(&req)
	if err != nil {
		err = fmt.Errorf("received error while unlocking mutex: %v", err)
		log.Error(err)
		return err
	}
	return nil
}
