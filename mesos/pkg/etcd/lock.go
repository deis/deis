package etcd

import (
	etcdlock "github.com/leeor/etcd-sync"
)

// AcquireLock creates a pseudo lock in etcd with a specific ttl
func AcquireLock(c *Client, key string, ttl uint64) error {
	c.lock = etcdlock.NewMutexFromClient(c.client, key, ttl)
	return c.lock.Lock()
}

// ReleaseLock releases the existing lock
func ReleaseLock(c *Client) {
	c.lock.Unlock()
}
