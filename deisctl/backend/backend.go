package backend

import "sync"

// Backend interface is used to interact with the cluster control plane
type Backend interface {
	Create([]string, *sync.WaitGroup, chan string, chan error)
	Destroy([]string, *sync.WaitGroup, chan string, chan error)
	Start([]string, *sync.WaitGroup, chan string, chan error)
	Stop([]string, *sync.WaitGroup, chan string, chan error)
	Scale(string, int, *sync.WaitGroup, chan string, chan error)
	SSH(string) error
	ListUnits() error
	ListUnitFiles() error
	Status(string) error
	Journal(string) error
}
