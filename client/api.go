package client

// Client interface used to interact with the cluster control plane
type Client interface {
	Create(string) error
	Destroy(string) error
	Start(string) error
	Stop(string) error
	Scale(string, int) error
	ListUnits() error
	ListUnitFiles() error
	Status(string) error
	Journal(string) error
}
