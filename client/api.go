package client

// Client interface used to interact with the cluster control plane
type Client interface {
	Create(string, bool) error
	Destroy(string) error
	Start(string, bool) error
	Stop(string) error
	Scale(string, int) error
	List() error
	Status(string) error
}
