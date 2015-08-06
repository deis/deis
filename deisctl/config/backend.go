package config

import "github.com/deis/deis/deisctl/config/model"

// Backend is an interface for any sort of underlying key/value config store
type Backend interface {
	Get(string) (string, error)
	GetWithDefault(string, string) (string, error)
	Set(string, string) (string, error)
	SetWithTTL(string, string, uint64) (string, error)
	Delete(string) error
	GetRecursive(string) ([]*model.ConfigNode, error)
}
