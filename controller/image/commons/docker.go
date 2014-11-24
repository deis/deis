package commons

import (
	"github.com/fsouza/go-dockerclient"
)

// NewClient returns a new docker test client.
func NewDockerClient() (*docker.Client, error) {
	endpoint := Getopt("DOCKER_HOST", "unix:///var/run/docker.sock")
	return docker.NewClient(endpoint)
}
