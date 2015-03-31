package docker

import (
	"github.com/deis/deis/pkg/os"

	"github.com/fsouza/go-dockerclient"
)

// NewDockerClient returns a new docker client
// By default the socket /var/run/docker.sock is used if there is no env DOCKER_HOST
func NewDockerClient() (*docker.Client, error) {
	endpoint := os.Getopt("DOCKER_HOST", "unix:///var/run/docker.sock")
	return docker.NewClient(endpoint)
}
