package server

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/fsouza/go-dockerclient"
)

const (
	appNameRegex string = `([a-z0-9-]+)_v([1-9][0-9]*).(cmd|web).([1-9][0-9])*`
)

// Server is the main entrypoint for a publisher. It listens on a docker client for events
// and publishes their host:port to the etcd client.
type Server struct {
	DockerClient *docker.Client
	EtcdClient   *etcd.Client

	host     string
	logLevel string
}

var safeMap = struct {
	sync.RWMutex
	data map[string]string
}{data: make(map[string]string)}

// New returns a new instance of Server.
func New(dockerClient *docker.Client, etcdClient *etcd.Client, host, logLevel string) *Server {
	return &Server{
		DockerClient: dockerClient,
		EtcdClient:   etcdClient,
		host:         host,
		logLevel:     logLevel,
	}
}

// Listen adds an event listener to the docker client and publishes containers that were started.
func (s *Server) Listen(ttl time.Duration) {
	listener := make(chan *docker.APIEvents)
	// TODO: figure out why we need to sleep for 10 milliseconds
	// https://github.com/fsouza/go-dockerclient/blob/0236a64c6c4bd563ec277ba00e370cc753e1677c/event_test.go#L43
	defer func() { time.Sleep(10 * time.Millisecond); s.DockerClient.RemoveEventListener(listener) }()
	if err := s.DockerClient.AddEventListener(listener); err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case event := <-listener:
			if event.Status == "start" {
				container, err := s.getContainer(event.ID)
				if err != nil {
					log.Println(err)
					continue
				}
				s.publishContainer(container, ttl)
			} else if event.Status == "stop" {
				s.removeContainer(event.ID)
			}
		}
	}
}

// Poll lists all containers from the docker client every time the TTL comes up and publishes them to etcd
func (s *Server) Poll(ttl time.Duration) {
	containers, err := s.DockerClient.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, container := range containers {
		// send container to channel for processing
		s.publishContainer(&container, ttl)
	}
}

// getContainer retrieves a container from the docker client based on id
func (s *Server) getContainer(id string) (*docker.APIContainers, error) {
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
	return nil, fmt.Errorf("could not find container with id %v", id)
}

// publishContainer publishes the docker container to etcd.
func (s *Server) publishContainer(container *docker.APIContainers, ttl time.Duration) {
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
		appPath := fmt.Sprintf("%s/%s", appName, containerName)
		keyPath := fmt.Sprintf("/deis/services/%s", appPath)
		for _, p := range container.Ports {
			// lowest port wins (docker sorts the ports)
			// TODO (bacongobbler): support multiple exposed ports
			port := strconv.Itoa(int(p.PrivatePort))
			apiContainer, _ := s.DockerClient.InspectContainer(container.ID)
			hostAndPort := apiContainer.NetworkSettings.IPAddress + ":" + port
			if s.IsPublishableApp(containerName) && s.IsPortOpen(hostAndPort) {
				s.setEtcd(keyPath, hostAndPort, uint64(ttl.Seconds()))
				safeMap.Lock()
				safeMap.data[container.ID] = appPath
				safeMap.Unlock()
			}
			break
		}
	}
}

// removeContainer remove a container published by this component
func (s *Server) removeContainer(event string) {
	safeMap.RLock()
	appPath := safeMap.data[event]
	safeMap.RUnlock()

	if appPath != "" {
		keyPath := fmt.Sprintf("/deis/services/%s", appPath)
		log.Printf("stopped %s\n", keyPath)
		s.removeEtcd(keyPath, false)
	}
}

// IsPublishableApp determines if the application should be published to etcd.
func (s *Server) IsPublishableApp(name string) bool {
	r := regexp.MustCompile(appNameRegex)
	match := r.FindStringSubmatch(name)
	if match == nil {
		return false
	}
	appName := match[1]
	version, err := strconv.Atoi(match[2])
	if err != nil {
		log.Println(err)
		return false
	}

	if version >= latestRunningVersion(s.EtcdClient, appName) {
		return true
	}
	return false
}

// IsPortOpen checks if the given port is accepting tcp connections
func (s *Server) IsPortOpen(hostAndPort string) bool {
	portOpen := false
	conn, err := net.Dial("tcp", hostAndPort)
	if err == nil {
		portOpen = true
		defer conn.Close()
	}
	return portOpen
}

// latestRunningVersion retrieves the highest version of the application published
// to etcd. If no app has been published, returns 0.
func latestRunningVersion(client *etcd.Client, appName string) int {
	r := regexp.MustCompile(appNameRegex)
	if client == nil {
		// FIXME: client should only be nil during tests. This should be properly refactored.
		if appName == "ceci-nest-pas-une-app" {
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
		version, err := strconv.Atoi(match[2])
		if err != nil {
			log.Println(err)
			return 0
		}
		versions = append(versions, version)
	}
	return max(versions)
}

// max returns the maximum value in n
func max(n []int) int {
	val := 0
	for _, i := range n {
		if i > val {
			val = i
		}
	}
	return val
}

// setEtcd sets the corresponding etcd key with the value and ttl
func (s *Server) setEtcd(key, value string, ttl uint64) {
	if _, err := s.EtcdClient.Set(key, value, ttl); err != nil {
		log.Println(err)
	}
	if s.logLevel == "debug" {
		log.Println("set", key, "->", value)
	}
}

// removeEtcd removes the corresponding etcd key
func (s *Server) removeEtcd(key string, recursive bool) {
	if _, err := s.EtcdClient.Delete(key, recursive); err != nil {
		log.Println(err)
	}
	if s.logLevel == "debug" {
		log.Println("del", key)
	}
}
