package mock

import (
	"fmt"
	"regexp"

	"github.com/deis/deis/deisctl/etcdclient"
)

// Store for a mocked etcd
type Store []*etcdclient.ServiceKey

// Client for the mocked etcd
type Client struct {
	Expected Store
}

// GetRecursive mocks returning a slice of all nodes under an etcd directory
func (m Client) GetRecursive(key string) ([]*etcdclient.ServiceKey, error) {
	r, _ := regexp.Compile(`^deis/services\s+`)
	var serviceKeys []*etcdclient.ServiceKey

	for _, expect := range m.Expected {
		if r.MatchString(expect.Key) {
			serviceKeys = append(serviceKeys, expect)
		}
	}
	return serviceKeys, nil
}

// Update a mocked etcd key
func (m Client) Update(key string, value string, ttl uint64) (string, error) {
	for _, expect := range m.Expected {
		if expect.Key == key {
			expect.TTL = int64(ttl)
			return expect.Key, nil
		}
	}
	return "", fmt.Errorf("%s does not exist", m.Expected)
}

// Get a mocked etcd key
func (m Client) Get(key string) (value string, err error) {
	for _, expect := range m.Expected {
		if expect.Key == key {
			return expect.Value, nil
		}
	}
	return "", fmt.Errorf("%s does not exist", m.Expected)
}

// Set a mocked etcd key
func (m Client) Set(key, value string) (returnedValue string, err error) {
	for _, expect := range m.Expected {
		if expect.Key == key {
			return value, nil
		}
	}
	return "", fmt.Errorf("%s does not exist", m.Expected)
}

// Delete a mocked etcd key
func (m Client) Delete(key string) (err error) {
	for _, expect := range m.Expected {
		if expect.Key == key {
			return nil
		}
	}
	return fmt.Errorf("%s does not exist", m.Expected)
}
