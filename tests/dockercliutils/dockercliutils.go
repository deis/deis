package dockercliutils

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/deis/deis/tests/utils"
	"github.com/dotcloud/docker/api/client"
)

// DaemonAddr returns the docker server address for testing.
func DaemonAddr() string {
	addr := os.Getenv("TEST_DAEMON_ADDR")
	if addr == "" {
		if utils.GetHostOs() == "darwin" {
			addr = "172.17.8.100:4243"
		} else {
			addr = "/var/run/docker.sock"
		}
	}
	return addr
}

// DaemonProto returns the docker server protocol for testing.
func DaemonProto() string {
	proto := os.Getenv("TEST_DAEMON_PROTO")
	if proto == "" {
		if utils.GetHostOs() == "darwin" {
			proto = "tcp"
		} else {
			proto = "unix"
		}
	}
	return proto
}

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
	ipaddr := os.Getenv("HOST_IPADDR")
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

// GetNewClient returns a new docker test client.
func GetNewClient() (
	cli *client.DockerCli, stdout *io.PipeReader, stdoutPipe *io.PipeWriter) {
	testDaemonAddr := DaemonAddr()
	testDaemonProto := DaemonProto()
	stdout, stdoutPipe = io.Pipe()
	cli = client.NewDockerCli(
		nil, stdoutPipe, nil, testDaemonProto, testDaemonAddr, nil)
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

// BuildImage builds and tags a docker image from a Dockerfile.
func BuildImage(t *testing.T, path string, tag string) error {
	var err error
	cli, stdout, stdoutPipe := GetNewClient()
	fmt.Println("Building docker image", tag)
	go func() {
		if err = cli.CmdBuild("--tag="+tag, path); err != nil {
			return
		}
		err = CloseWrap(stdout, stdoutPipe)
	}()
	PrintToStdout(t, stdout, stdoutPipe, "build docker file")
	return err
}

// GetInspectData prints and returns `docker inspect` data for a container.
func GetInspectData(t *testing.T, format string, container string) string {
	var inspectData string
	cli, stdout, stdoutPipe := GetNewClient()
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

// PullImage pulls a docker image from the docker index.
func PullImage(t *testing.T, cli *client.DockerCli, args ...string) {
	fmt.Println("pulling image :" + args[0])
	err := cli.CmdPull(args...)
	if err != nil {
		t.Fatalf("pulling Image Failed %s", err)
	}
}

// RunContainer runs a docker image with the given arguments.
func RunContainer(cli *client.DockerCli, args ...string) error {
	fmt.Println("--- Run docker container", args[1])
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
	cli, stdout, stdoutPipe := GetNewClient()
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
			fmt.Println("closing read/write pipe")
			time.Sleep(3000 * time.Millisecond)
			if err := CloseWrap(stdout, stdoutPipe); err != nil {
				t.Fatalf("Inspect Element %s", err)
			}
		}()
		PrintToStdout(t, stdout, stdoutPipe, "running"+args[1])
	}
}

func getContainerIds(t *testing.T, uid string) []string {
	sliceContainers := []string{}
	cli, stdout, stdoutPipe := GetNewClient()
	go func() {
		err := cli.CmdPs("-a")
		if err != nil {
			t.Fatalf("getContainerIds %s", err)
		}
		if err = CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("getContainerIds %s", err)
		}
	}()
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			fmt.Print(cmdBytes)
			if strings.Contains(cmdBytes, uid) {
				sliceContainers = utils.Append(sliceContainers, strings.Fields(cmdBytes)[0])
			}
		} else {
			break
		}

	}
	return sliceContainers
}

func getImageIds(t *testing.T, uid string) []string {
	sliceImageids := []string{}
	cli, stdout, stdoutPipe := GetNewClient()
	go func() {
		err := cli.CmdImages()
		if err != nil {
			t.Fatalf("getImages Ids %s", err)
		}
		if err = CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("getImages Ids %s", err)
		}
	}()
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			fmt.Print(cmdBytes)
			if strings.Contains(cmdBytes, uid) {
				sliceImageids = utils.Append(sliceImageids, strings.Fields(cmdBytes)[2])
			}
		} else {
			break
		}

	}
	return sliceImageids
}

func stopContainers(t *testing.T, sliceContainerIds []string) {
	cli, stdout, stdoutPipe := GetNewClient()
	go func() {
		for _, value := range sliceContainerIds {
			err := cli.CmdStop(value)
			if err != nil {
				t.Log("stop container failed:", err)
			}
		}
		if err := CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("stop Container %s", err)
		}
	}()
	PrintToStdout(t, stdout, stdoutPipe, "removing container")
}

// ClearTestSession cleans up after a typical test session.
func ClearTestSession(t *testing.T, uid string) {
	fmt.Println("--- Clear test session", uid)
	sliceContainerIds := getContainerIds(t, uid)
	stopContainers(t, sliceContainerIds)
}

// GetImageID returns the ID of a docker image.
func GetImageID(t *testing.T, repo string) string {
	var imageID string
	cli, stdout, stdoutPipe := GetNewClient()
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

// RunEtcdTest starts an etcd docker container for testing.
func RunEtcdTest(t *testing.T, uid string, port string) {
	var err error
	cli, stdout, stdoutPipe := GetNewClient()
	etcdImage := "deis/test-etcd:latest"
	done2 := make(chan bool, 1)
	go func() {
		done2 <- true
		ipaddr := utils.GetHostIPAddress()
		err = RunContainer(cli,
			"--name", "deis-etcd-"+uid,
			"--rm",
			"-p", port+":"+port,
			"-e", "HOST_IP="+ipaddr,
			"-e", "ETCD_ADDR="+ipaddr+":"+port,
			etcdImage)
	}()
	go func() {
		<-done2
		fmt.Println("closing read/write pipe")
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
