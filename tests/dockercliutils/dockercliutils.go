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
		addr = "/var/run/docker.sock"
	}
	return addr
}

// DaemonProto returns the docker server protocol for testing.
func DaemonProto() string {
	proto := os.Getenv("TEST_DAEMON_PROTO")
	if proto == "" {
		proto = "unix"
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
	IPAddress := os.Getenv("HOST_IPADDR")
	if IPAddress == "" {
		IPAddress = GetInspectData(
			t, "{{ .NetworkSettings.IPAddress }}", container)
	}
	fmt.Println("Running service test for " + container)
	if strings.Contains(IPAddress, "Error") {
		t.Fatalf("wrong IP %s", IPAddress)
	}
	if protocol == "http" {
		url := "http://" + IPAddress + ":" + port
		response, err := http.Get(url)
		if err != nil {
			t.Fatalf("Not reachable %s", err)
		}
		fmt.Println(response)
	}
	if protocol == "tcp" || protocol == "udp" {
		conn, err := net.Dial(protocol, IPAddress+":"+port)
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
	var testDaemonAddr, testDaemonProto string
	if utils.GetHostOs() == "darwin" {
		testDaemonAddr = "172.17.8.100:4243"
		testDaemonProto = "tcp"
	} else {
		testDaemonAddr = DaemonAddr()
		testDaemonProto = DaemonProto()
	}
	stdout, stdoutPipe = io.Pipe()
	cli = client.NewDockerCli(
		nil, stdoutPipe, nil, testDaemonProto, testDaemonAddr, nil)
	return
}

// PrintToStdout prints a string to stdout.
func PrintToStdout(t *testing.T, stdout *io.PipeReader,
	stdoutPipe *io.PipeWriter, stoptag string) string {
	var result string
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			result = cmdBytes
			fmt.Print(cmdBytes)
			if strings.Contains(cmdBytes, stoptag) == true {
				if err := CloseWrap(stdout, stdoutPipe); err != nil {
					t.Fatalf("Closewraps %s", err)
				}
			}
		} else {
			break
		}
	}
	return result
}

// BuildDockerfile builds and tags a docker image from a Dockerfile.
func BuildDockerfile(t *testing.T, path string, tag string) {
	cli, stdout, stdoutPipe := GetNewClient()
	fmt.Println("Building docker file :" + tag)
	go func() {
		err := cli.CmdBuild("--tag="+tag, path)
		if err != nil {
			t.Fatalf(" %s", err)
		}
		if err = CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("buildDockerfile %s", err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	PrintToStdout(t, stdout, stdoutPipe, "build docker file")
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
func RunContainer(t *testing.T, cli *client.DockerCli, args ...string) {
	fmt.Println("Running docker container :" + args[1])
	err := cli.CmdRun(args...)
	if err != nil {
		// Ignore certain errors we see in io handling.
		switch msg := err.Error(); {
		case strings.Contains(msg, "read/write on closed pipe"):
			return
		case strings.Contains(msg, "Code: -1"):
			return
		case strings.Contains(msg, "Code: 2"):
			return
		default:
			t.Fatalf("RunContainer failed: %v", err)
		}
	}
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
			RunContainer(t, cli, args...)
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

func removeImages(t *testing.T, sliceImageIds []string) {
	cli, stdout, stdoutPipe := GetNewClient()
	go func() {
		for _, value := range sliceImageIds {
			err := cli.CmdRmi("-f", value)
			if err != nil {
				if !((strings.Contains(fmt.Sprintf("%s", err), "No such image")) || (strings.Contains(fmt.Sprintf("%s", err), "one or more"))) {
					t.Fatalf("removeImages %s", err)
				}
			}
		}
		if err := CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("remove Images %s", err)
		}
	}()
	PrintToStdout(t, stdout, stdoutPipe, "removing container")
}

// ClearTestSession cleans up after a typical test session.
func ClearTestSession(t *testing.T, uid string) {
	fmt.Println("clearing test session", uid)
	sliceContainerIds := getContainerIds(t, uid)
	// sliceImageids := getImageIds(t, uid)
	// //fmt.Println(sliceContainerIds)
	// //fmt.Println(sliceImageids)
	// fmt.Println("removing containers and images for the test session " + uid)
	stopContainers(t, sliceContainerIds)
	// removeImages(t, sliceImageids)
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
	//docker run -t -i --name=deis-etcd -p 4001:4001  -e HOST_IP=172.17.8.100
	// -e ETCD_ADDR=172.17.8.100:4001
	// --entrypoint=/bin/bash phife.atribecalledchris.com:5000/deis/etcd:0.3.0
	// -c /usr/local/bin/etcd
	cli, stdout, stdoutPipe := GetNewClient()
	done2 := make(chan bool, 1)
	IPAddress := utils.GetHostIPAddress()
	go func() {
		PullImage(t, cli, "deis/test-etcd:latest")
		done2 <- true
		RunContainer(t, cli,
			"--name", "deis-etcd-"+uid,
			"--rm",
			"-p", port+":"+port,
			"-e", "HOST_IP="+IPAddress,
			"-e", "ETCD_ADDR="+IPAddress+":"+port,
			"deis/test-etcd:latest")
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
}
