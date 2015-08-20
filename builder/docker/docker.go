package docker

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/log"
	"github.com/Masterminds/cookoo/safely"
	"github.com/deis/deis/builder/etcd"
	docli "github.com/fsouza/go-dockerclient"
)

// Path to the Docker unix socket.
var DockSock = "/var/run/docker.sock"

// Cleanup removes any existing Docker artifacts.
//
// Returns true if the file exists (and was deleted), or false if no file
// was deleted.
func Cleanup(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {

	// If info is returned, then the file is there. If we get an error, we're
	// pretty much not going to be able to remove the file (which probably
	// doesn't exist).
	if _, err := os.Stat(DockSock); err == nil {
		log.Infof(c, "Removing leftover docker socket %s: %v", DockSock)
		return true, os.Remove(DockSock)
	}
	return false, nil
}

// CreateClient creates a new Docker client.
//
// Params:
// 	- url (string): The URI to the Docker daemon. This defaults to the UNIX
// 		socket /var/run/docker.sock.
//
// Returns:
// 	- *docker.Client
//
func CreateClient(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	path := p.Get("url", "unix:///var/run/docker.sock").(string)

	return docli.NewClient(path)
}

// Start starts a Docker daemon.
//
// This assumes the presence of the docker client on the host. It does not use
// the API.
func Start(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dargs := []string{
		"-d",
		"--bip=172.19.42.1/16",
		"--insecure-registry",
		"10.0.0.0/8",
		"--insecure-registry",
		"172.16.0.0/12",
		"--insecure-registry",
		"192.168.0.0/16",
		"--insecure-registry",
		"100.64.0.0/10",
		//"--storage-driver=overlay", // Add this dynamically.
	}

	// There is probably a better way to do this. I'm not sure whether the
	// original intent of the mkdir was to accomplish something, or just check
	// whether overlay is supported.
	if err := os.MkdirAll("/", 0700); err == nil {

		cmd := exec.Command("findmnt", "--noheadings", "--output", "FSTYPE", "--target", "/")

		if out, err := cmd.Output(); err == nil && strings.TrimSpace(string(out)) == "overlay" {
			dargs = append(dargs, "--storage-driver=overlay")
		} else {
			log.Infof(c, "File system type: '%s' (%v)", out, err)
		}
	}

	//# spawn a docker daemon to run builds
	//docker -d --bip=172.19.42.1/16 $DRIVER_OVERRIDE --insecure-registry 10.0.0.0/8 --insecure-registry 172.16.0.0/12
	// --insecure-registry 192.168.0.0/16 --insecure-registry 100.64.0.0/10 &
	//DOCKER_PID=$!

	log.Infof(c, "Starting docker with %s", strings.Join(dargs, " "))
	cmd := exec.Command("docker", dargs...)
	if err := cmd.Start(); err != nil {
		log.Errf(c, "Failed to start Docker. %s", err)
		return -1, err
	}
	// Get the PID and return it.
	return cmd.Process.Pid, nil
}

// WaitForStart delays until Docker appears to be up and running.
//
// Params:
// 	- client (*docker.Client): Docker client.
// 	- timeout (time.Duration): Time after which to give up.
//
// Returns:
// 	- boolean true if the server is up.
func WaitForStart(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	ok, missing := p.RequiresValue("client")
	if !ok {
		return nil, &cookoo.FatalError{"Missing required fields: " + strings.Join(missing, ", ")}
	}
	cli := p.Get("client", nil).(*docli.Client)
	timeout := p.Get("timeout", 30*time.Second).(time.Duration)

	keepon := true
	timer := time.AfterFunc(timeout, func() {
		keepon = false
	})

	for keepon == true {
		if err := cli.Ping(); err == nil {
			timer.Stop()
			log.Infof(c, "Docker is running.")
			return true, nil
		}
		time.Sleep(time.Second)
	}
	return false, fmt.Errorf("Docker timed out after waiting %s for server.", timeout)
}

// BuildImg describes a build image.
type BuildImg struct {
	Path, Tag string
}

// ParallelBuild runs multiple docker builds at the same time.
//
// Params:
//	-images ([]BuildImg): Images to build
// 	-alwaysFetch (bool): Default false. If set to true, this will always fetch
// 		the Docker image even if it already exists in the registry.
//
// Returns:
//
// 	- Waiter: A *sync.WaitGroup that is waiting for the docker downloads to finish.
//
// Context:
//
// This puts 'ParallelBuild.failN" (int) into the context to indicate how many failures
// occurred during fetches.
func ParallelBuild(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	images := p.Get("images", []BuildImg{}).([]BuildImg)
	//client := p.Get("client", nil).(*docli.Client)

	var wg sync.WaitGroup
	var m sync.Mutex
	var fails int

	for _, img := range images {
		img := img
		wg.Add(1)
		safely.GoDo(c, func() {
			log.Infof(c, "Starting build for %s (tag: %s)", img.Path, img.Tag)
			_, err := buildImg(c, img.Path, img.Tag)
			//err := build(c, img.Path, img.Tag, client)
			if err != nil {
				log.Errf(c, "Failed to build docker image: %s", err)
				m.Lock()
				fails++
				m.Unlock()
			}
			wg.Done()
		})

	}

	// Number of failures.
	c.Put("ParallelBuild.failN", fails)

	//wg.Wait()
	//log.Infof(c, "Docker images are built.")
	return &wg, nil
}

// Waiter describes a thing that can wait.
//
// It does not bring you food. I should know. I tried.
type Waiter interface {
	Wait()
}

// Wait waits for a sync.WaitGroup to finish.
//
// Params:
// 	- wg (Waiter): The thing to wait for.
// 	- msg (string): The message to print when done. If this is empty, nothing is sent.
// 	- waiting (string): String to tell what we're waiting for. If empty, nothing is displayed.
// 	- failures (int): The number of failures that occurred while waiting.
//
// Returns:
//  Nothing.
func Wait(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	ok, missing := p.RequiresValue("wg")
	if !ok {
		return nil, &cookoo.FatalError{"Missing required fields: " + strings.Join(missing, ", ")}
	}
	wg := p.Get("wg", nil).(Waiter)
	msg := p.Get("msg", "").(string)
	fails := p.Get("failures", 0).(int)
	waitmsg := p.Get("waiting", "").(string)

	if len(waitmsg) > 0 {
		log.Info(c, waitmsg)
	}

	wg.Wait()
	if len(msg) > 0 {
		log.Info(c, msg)
	}

	if fails > 0 {
		return nil, fmt.Errorf("There were %d failures while waiting.", fails)
	}
	return nil, nil
}

// BuildImage builds a docker image.
//
// Essentially, this executes:
// 	docker build -t TAG PATH
//
// Params:
// 	- path (string): The path to the image. REQUIRED
// 	- tag (string): The tag to build.
func BuildImage(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	path := p.Get("path", "").(string)
	tag := p.Get("tag", "").(string)

	log.Infof(c, "Building docker image %s (tag: %s)", path, tag)

	return buildImg(c, path, tag)
}

func buildImg(c cookoo.Context, path, tag string) ([]byte, error) {
	dargs := []string{"build"}
	if len(tag) > 0 {
		dargs = append(dargs, "-t", tag)
	}

	dargs = append(dargs, path)

	out, err := exec.Command("docker", dargs...).CombinedOutput()
	if len(out) > 0 {
		log.Infof(c, "Docker: %s", out)
	}
	return out, err
}

// Push pushes an image to the registry.
//
// This finds the appropriate registry by looking it up in etcd.
//
// Params:
// - client (etcd.Getter): Client to do etcd lookups.
// - tag (string): Tag to push.
//
// Returns:
//
func Push(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	// docker tag deis/slugrunner:lastest HOST:PORT/deis/slugrunner:latest
	// docker push HOST:PORT/deis/slugrunner:latest
	client := p.Get("client", nil).(etcd.Getter)

	host, err := client.Get("/deis/registry/host", false, false)
	if err != nil {
		return nil, err
	}
	port, err := client.Get("/deis/registry/port", false, false)
	if err != nil {
		return nil, err
	}

	registry := host.Node.Value + ":" + port.Node.Value
	tag := p.Get("tag", "").(string)

	log.Infof(c, "Pushing %s to %s. This may take some time.", tag, registry)
	rem := path.Join(registry, tag)

	out, err := exec.Command("docker", "tag", "-f", tag, rem).CombinedOutput()
	if err != nil {
		log.Warnf(c, "Failed to tag %s on host %s: %s (%s)", tag, rem, err, out)
		//return nil, err
	}
	out, err = exec.Command("docker", "-D", "push", rem).CombinedOutput()

	if err != nil {
		log.Warnf(c, "Failed to push %s to host %s: %s (%s)", tag, rem, err, out)
		return nil, err
		//return out, err
	}
	log.Infof(c, "Finished pushing %s to %s.", tag, registry)
	//return out, err
	return nil, nil
}

/*
 * This function only works for very simple docker files that do not have
 * local resources.
 * Need to suck in all of the files in ADD directives, too.
 */
// build takes a Dockerfile and builds an image.
func build(c cookoo.Context, path, tag string, client *docli.Client) error {
	dfile := filepath.Join(path, "Dockerfile")

	// Stat the file
	info, err := os.Stat(dfile)
	if err != nil {
		return fmt.Errorf("Dockerfile stat: %s", err)
	}
	file, err := os.Open(dfile)
	if err != nil {
		return fmt.Errorf("Dockerfile open: %s", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	tw.WriteHeader(&tar.Header{
		Name:    "Dockerfile",
		Size:    info.Size(),
		ModTime: info.ModTime(),
	})
	io.Copy(tw, file)
	if err := tw.Close(); err != nil {
		return fmt.Errorf("Dockerfile tar: %s", err)
	}

	options := docli.BuildImageOptions{
		Name:         tag,
		InputStream:  &buf,
		OutputStream: os.Stdout,
	}
	return client.BuildImage(options)
}
