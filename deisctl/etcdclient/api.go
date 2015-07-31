package etcdclient

// Client for etcd
type Client interface {
	Get(string) (string, error)
	Set(string, string) (string, error)
	Delete(string) error
	GetRecursive(string) ([]*ServiceKey, error)
	Update(string, string, uint64) (string, error)
}
