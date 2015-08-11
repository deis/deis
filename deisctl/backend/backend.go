package backend

import (
	"io"
	"sync"
)

// Backend interface is used to interact with the cluster control plane
type Backend interface {
	Create([]string, *sync.WaitGroup, io.Writer, io.Writer)
	Destroy([]string, *sync.WaitGroup, io.Writer, io.Writer)
	Start([]string, *sync.WaitGroup, io.Writer, io.Writer)
	Stop([]string, *sync.WaitGroup, io.Writer, io.Writer)
	Scale(string, int, *sync.WaitGroup, io.Writer, io.Writer)
	RollingRestart(string, *sync.WaitGroup, io.Writer, io.Writer)
	SSH(string) error
	SSHExec(string, string) error
	Dock(string, []string) error
	ListUnits() error
	ListUnitFiles() error
	Status(string) error
	Journal(string) error
}
