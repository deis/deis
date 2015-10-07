package mock

import (
	"fmt"
	"regexp"

	"github.com/deis/deis/deisctl/config/model"
)

// Store for a mocked config backend
type Store []*model.ConfigNode

// ConfigBackend is an in memory "mock" config datastore used for testing
type ConfigBackend struct {
	Expected Store
}

// Get a value by key from an in memory config backend
func (cb ConfigBackend) Get(key string) (value string, err error) {
	for _, expect := range cb.Expected {
		if expect.Key == key {
			return expect.Value, nil
		}
	}
	return "", fmt.Errorf("%s does not exist", cb.Expected)
}

// GetWithDefault gets a value by key from an in memory config backend and
// return a default value if not found
func (cb ConfigBackend) GetWithDefault(key string, defaultValue string) (string, error) {
	for _, expect := range cb.Expected {
		if expect.Key == key {
			return expect.Value, nil
		}
	}
	return defaultValue, nil
}

// Set a value for the specified key in an in memory config backend
func (cb ConfigBackend) Set(key, value string) (returnedValue string, err error) {
	for _, expect := range cb.Expected {
		if expect.Key == key {
			return value, nil
		}
	}
	return "", fmt.Errorf("%s does not exist", cb.Expected)
}

// Delete a key/value pair by key from an in memory config backend
func (cb ConfigBackend) Delete(key string) (err error) {
	for _, expect := range cb.Expected {
		if expect.Key == key {
			return nil
		}
	}
	return fmt.Errorf("%s does not exist", cb.Expected)
}

// GetRecursive returns a slice of all key/value pairs "under" a specified key
// in an in memory config backend (this is assuming some hierarchichal
// order exists wherein the value corresponding to a key may in fact be another
// key/value pair)
func (cb ConfigBackend) GetRecursive(key string) ([]*model.ConfigNode, error) {
	r, _ := regexp.Compile(`^deis/services\s+`)
	var configNodes []*model.ConfigNode

	for _, expect := range cb.Expected {
		if r.MatchString(expect.Key) {
			configNodes = append(configNodes, expect)
		}
	}
	return configNodes, nil
}

// SetWithTTL sets a value for the specified key in an in memory config
// backend-- with a time to live
func (cb ConfigBackend) SetWithTTL(key string, value string, ttl uint64) (string, error) {
	for _, expect := range cb.Expected {
		if expect.Key == key {
			expect.TTL = int64(ttl)
			return expect.Key, nil
		}
	}
	return "", fmt.Errorf("%s does not exist", cb.Expected)
}
