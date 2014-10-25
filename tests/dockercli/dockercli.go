// Package dockercli provides helper functions for testing with Docker.

package dockercli

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/deis/deis/tests/utils"
	"github.com/docker/docker/api/client"
)

// CloseWrap ensures that an io.Writer is closed.
func CloseWrap(args ...io.Closer) error {
	e := false
	ret := fmt.Errorf("Error closing elements")
	for _, c := range args {
		if err := c.Close(); err != nil {
			e = true
			ret = fmt.Errorf("%s\n%s", ret, err)
		}
	}
	if e {
		return ret
	}
	return nil
}

// DeisServiceTest tries to connect to a container and port using the
// specified protocol.
func DeisServiceTest(
	t *testing.T, container string, port string, protocol string) {
	ipaddr := utils.HostAddress()
	if ipaddr == "" {
		ipaddr = GetInspectData(
			t, "{{ .NetworkSettings.ipaddr }}", container)
	}
	fmt.Println("Running service test for " + container)
	if strings.Contains(ipaddr, "Error") {
		t.Fatalf("wrong IP %s", ipaddr)
	}
	if protocol == "http" {
		url := "http://" + ipaddr + ":" + port
		response, err := http.Get(url)
		if err != nil {
			t.Fatalf("Not reachable %s", err)
		}
		fmt.Println(response)
	}
	if protocol == "tcp" || protocol == "udp" {
		conn, err := net.Dial(protocol, ipaddr+":"+port)
		if err != nil {
			t.Fatalf("Not reachable %s", err)
		}
		_, err = conn.Write([]byte("HEAD"))
		if err != nil {
			t.Fatalf("Not reachable %s", err)
		}
	}
}

// DockerHost returns the protocol and address of the docker server.
func DockerHost() (string, string, error) {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}
	u, err := url.Parse(dockerHost)
	if err != nil {
		return "", "", err
	}
	if u.Scheme == "unix" {
		return u.Scheme, u.Path, nil
	}
	return u.Scheme, u.Host, nil
}

// NewClient returns a new docker test client.
func NewClient() (
	cli *client.DockerCli, stdout *io.PipeReader, stdoutPipe *io.PipeWriter) {
	proto, addr, _ := DockerHost()
	stdout, stdoutPipe = io.Pipe()
	cli = client.NewDockerCli(nil, stdoutPipe, nil, nil, proto, addr, nil)
	return
}

// PrintToStdout prints a string to stdout.
func PrintToStdout(t *testing.T, stdout *io.PipeReader,
	stdoutPipe *io.PipeWriter, stoptag string) string {
	var result string
	r := bufio.NewReader(stdout)
	for {
		cmdBytes, err := r.ReadString('\n')
		if err != nil {
			break
		}
		result = cmdBytes
		fmt.Print(cmdBytes)
		if strings.Contains(cmdBytes, stoptag) == true {
			if err := CloseWrap(stdout, stdoutPipe); err != nil {
				t.Fatal(err)
			}
		}
	}
	return result
}

// GetInspectData prints and returns `docker inspect` data for a container.
func GetInspectData(t *testing.T, format string, container string) string {
	var inspectData string
	cli, stdout, stdoutPipe := NewClient()
	fmt.Println("Getting inspect data :" + format + ":" + container)
	go func() {
		err := cli.CmdInspect("--format", format, container)
		if err != nil {
			fmt.Printf("%s %s", format, err)
		}
		if err = CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("inspect data failed %s", err)
		}
	}()
	go func() {
		time.Sleep(3000 * time.Millisecond)
		if err := CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("Inspect data %s", err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	inspectData = PrintToStdout(t, stdout, stdoutPipe, "get inspect data")
	return strings.TrimSuffix(inspectData, "\n")
}

// RunContainer runs a docker image with the given arguments.
func RunContainer(cli *client.DockerCli, args ...string) error {
	// fmt.Println("--- Run docker container", args[1])
	err := cli.CmdRun(args...)
	if err != nil {
		// Ignore certain errors we see in io handling.
		switch msg := err.Error(); {
		case strings.Contains(msg, "read/write on closed pipe"):
			err = nil
		case strings.Contains(msg, "Code: -1"):
			err = nil
		case strings.Contains(msg, "Code: 2"):
			err = nil
		}
	}
	return err
}

// RunDeisDataTest starts a data container as a prerequisite for a service.
func RunDeisDataTest(t *testing.T, args ...string) {
	done := make(chan bool, 1)
	cli, stdout, stdoutPipe := NewClient()
	var hostname string
	fmt.Println(args[2] + " test")
	hostname = GetInspectData(t, "{{ .Config.Hostname }}", args[1])
	fmt.Println("data container " + hostname)
	done <- true
	if strings.Contains(hostname, "Error") {
		go func() {
			<-done
			if err := RunContainer(cli, args...); err != nil {
				t.Fatal(err)
			}
		}()
		go func() {
			time.Sleep(3000 * time.Millisecond)
			if err := CloseWrap(stdout, stdoutPipe); err != nil {
				t.Fatalf("Inspect Element %s", err)
			}
		}()
		PrintToStdout(t, stdout, stdoutPipe, "running"+args[1])
	}
}

// GetImageID returns the ID of a docker image.
func GetImageID(t *testing.T, repo string) string {
	var imageID string
	cli, stdout, stdoutPipe := NewClient()
	go func() {
		err := cli.CmdImages()
		if err != nil {
			t.Fatalf("GetImageID %s", err)
		}
		if err = CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("GetImageID %s", err)
		}
	}()
	imageID = PrintToStdout(t, stdout, stdoutPipe, repo)
	return strings.Fields(imageID)[2]
}

// RunTestEtcd starts an etcd docker container for testing.
func RunTestEtcd(t *testing.T, name string, port string) {
	var err error
	cli, stdout, stdoutPipe := NewClient()
	etcdImage := "deis/test-etcd:latest"
	ipaddr := utils.HostAddress()
	etcdAddr := ipaddr + ":" + port
	fmt.Printf("--- Running deis/test-etcd at %s\n", etcdAddr)
	done2 := make(chan bool, 1)
	go func() {
		done2 <- true
		_ = cli.CmdRm("-f", name)
		err = RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":"+port,
			"-e", "HOST_IP="+ipaddr,
			"-e", "ETCD_ADDR="+etcdAddr,
			etcdImage)
	}()
	go func() {
		<-done2
		time.Sleep(5000 * time.Millisecond)
		if err := CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("runEtcdTest %s", err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	PrintToStdout(t, stdout, stdoutPipe, "pulling etcd")
	if err != nil {
		t.Fatal(err)
	}
}
