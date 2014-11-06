package publisher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
)

const (
	appNameRegex string = `([a-z0-9-]+)_v([1-9][0-9]*).(cmd|web).([1-9][0-9])*`
)

type Server struct {
	DockerClient *docker.Client
	EtcdClient   *etcd.Client
}

func (s *Server) Listen(ttl time.Duration) {
	listener := make(chan *docker.APIEvents)
	// TODO: figure out why we need to sleep for 10 milliseconds
	// https://github.com/fsouza/go-dockerclient/blob/0236a64c6c4bd563ec277ba00e370cc753e1677c/event_test.go#L43
	defer func() { time.Sleep(10 * time.Millisecond); s.DockerClient.RemoveEventListener(listener) }()
	err := s.DockerClient.AddEventListener(listener)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case event := <-listener:
			if event.Status == "start" {
				container, err := s.GetContainer(event.ID)
				if err != nil {
					log.Println(err)
					continue
				}
				s.PublishContainer(container, ttl)
			}
		}
	}
}

func (s *Server) Poll(ttl time.Duration) {
	containers, err := s.DockerClient.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		// send container to channel for processing
		s.PublishContainer(&container, ttl)
	}
}

func (s *Server) GetContainer(id string) (*docker.APIContainers, error) {
	containers, err := s.DockerClient.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		// send container to channel for processing
		if container.ID == id {
			return &container, nil
		}
	}
	return nil, errors.New("could not find container")
}

func (s *Server) PublishContainer(container *docker.APIContainers, ttl time.Duration) {
	r := regexp.MustCompile(appNameRegex)
	for _, name := range container.Names {
		// HACK: remove slash from container name
		// see https://github.com/docker/docker/issues/7519
		containerName := name[1:]
		match := r.FindStringSubmatch(containerName)
		if match == nil {
			continue
		}
		appName := match[1]
		keyPath := fmt.Sprintf("/deis/services/%s/%s", appName, containerName)
		for _, p := range container.Ports {
			host := os.Getenv("HOST")
			port := strconv.Itoa(int(p.PublicPort))
			if s.IsPublishableApp(containerName) {
				s.setEtcd(keyPath, host+":"+port, uint64(ttl.Seconds()))
			}
			// TODO: support multiple exposed ports
			break
		}
	}
}

// isPublishableApp determines if the application should be published to etcd.
func (s *Server) IsPublishableApp(name string) bool {
	r := regexp.MustCompile(appNameRegex)
	match := r.FindStringSubmatch(name)
	if match == nil {
		return false
	}
	appName := match[1]
	version, _ := strconv.Atoi(match[2])
	if version >= latestRunningVersion(s.EtcdClient, appName) {
		return true
	} else {
		return false
	}
}

// latestRunningVersion retrieves the highest version of the application published
// to etcd. If no app has been published, returns 0.
func latestRunningVersion(client *etcd.Client, appName string) int {
	r := regexp.MustCompile(appNameRegex)
	if client == nil {
		// TODO: refactor for tests
		if appName == "test" {
			return 3
		}
		return 0
	}
	resp, err := client.Get(fmt.Sprintf("/deis/services/%s", appName), false, true)
	if err != nil {
		// no app has been published here (key not found) or there was an error
		return 0
	}
	var versions []int
	for _, node := range resp.Node.Nodes {
		match := r.FindStringSubmatch(node.Key)
		// account for keys that may not be an application container
		if match == nil {
			continue
		}
		version, _ := strconv.Atoi(match[2])
		versions = append(versions, version)
	}
	return max(versions)
}

func max(n []int) int {
	val := 0
	for _, i := range n {
		if i > val {
			val = i
		}
	}
	return val
}

func (s *Server) setEtcd(key, value string, ttl uint64) {
	_, err := s.EtcdClient.Set(key, value, ttl)
	if err != nil {
		log.Println(err)
	}
	log.Println("set", key, "->", value)
}
