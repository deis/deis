package dockercliutils

import (
	"bufio"
	"fmt"
	"github.com/deis/deis/tests/utils"
	"github.com/dotcloud/docker/api/client"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	unitTestStoreBase    = "/var/lib/docker/unit-tests"
	testDaemonAddr       = "172.17.8.100:4243"
	testDaemonProto      = "tcp"
	testDaemonHttpsProto = "tcp"
)

func DaemonAddr() string {
	addr := os.Getenv("TEST_DAEMON_ADDR")
	if addr == "" {
		addr = "/var/run/docker.sock"
	}
	return addr
}

func DaemonProto() string {
	proto := os.Getenv("TEST_DAEMON_PROTO")
	if proto == "" {
		proto = "unix"
	}
	return proto
}

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

func GetNewClient() (cli *client.DockerCli, stdout *io.PipeReader, stdoutPipe *io.PipeWriter) {
	stdout, stdoutPipe = io.Pipe()
	cli = client.NewDockerCli(nil, stdoutPipe, nil, DaemonProto(), DaemonAddr(), nil)
	return
}

func PrintToStdout(t *testing.T, stdout *io.PipeReader, stdoutPipe *io.PipeWriter, stoptag string) string {
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

func BuildDockerfile(t *testing.T, path string, tag string) {
	cli, stdout, stdoutPipe := GetNewClient()
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
	PrintToStdout(t, stdout, stdoutPipe, "Building docker file")
}

func GetInspectData(t *testing.T, format string, container string) string {
	var inspectData string
	cli, stdout, stdoutPipe := GetNewClient()
	go func() {
		err := cli.CmdInspect("--format", format, container)
		if err != nil {
			fmt.Printf("%s %s", format, err)
		}
		if err = CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("getIPAdressTest %s", err)
		}
	}()
	go func() {
		fmt.Println("here")
		time.Sleep(3000 * time.Millisecond)
		if err := CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("Inspect Element %s", err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	inspectData = PrintToStdout(t, stdout, stdoutPipe, "IPAddress")
	return strings.TrimSuffix(inspectData, "\n")

}

func PullImage(t *testing.T, cli *client.DockerCli, args ...string) {
	err := cli.CmdPull(args...)
	if err != nil {
		t.Fatalf("pulling Image Failed %s", err)
	}
}

func RunContainer(t *testing.T, cli *client.DockerCli, args ...string) {
	err := cli.CmdRun(args...)
	if err != nil {
		if strings.Contains(fmt.Sprintf("%s", err), "read/write on closed pipe") == false {
			t.Fatalf("running Image failed %s", err)
		}
	}
}

func RunDeisDataTest(t *testing.T, args ...string) {
	done := make(chan bool, 1)
	cli, stdout, stdoutPipe := GetNewClient()
	var hostname string
	hostname = GetInspectData(t, "{{ .Config.Hostname }}", args[1])
	fmt.Println("data container " + hostname)
	done <- true
	if strings.Contains(hostname, "Error") {
		go func() {
			<-done
			RunContainer(t, cli, args...)
		}()
		go func() {
			fmt.Println("here")
			time.Sleep(3000 * time.Millisecond)
			if err := CloseWrap(stdout, stdoutPipe); err != nil {
				t.Fatalf("Inspect Element %s", err)
			}
		}()
		PrintToStdout(t, stdout, stdoutPipe, "running"+args[1])
		fmt.Println("pulling Deis registry data")
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

func stop_rmContainers(t *testing.T, sliceContainerIds []string) {
	cli, stdout, stdoutPipe := GetNewClient()
	done := make(chan bool, 1)
	go func() {
		for _, value := range sliceContainerIds {
			err := cli.CmdStop(value)
			if err != nil {
				t.Fatalf("stop Container %s", err)
			}
		}
		done <- true
	}()
	go func() {
		<-done
		for _, value := range sliceContainerIds {
			err := cli.CmdRm(value)
			if err != nil {
				t.Fatalf("stop Container %s", err)
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

func ClearTestSession(t *testing.T, uid string) {
	sliceContainerIds := getContainerIds(t, uid)
	sliceImageids := getImageIds(t, uid)
	//fmt.Println(sliceContainerIds)
	//fmt.Println(sliceImageids)
	stop_rmContainers(t, sliceContainerIds)
	removeImages(t, sliceImageids)
}

func GetImageId(t *testing.T, repo string) string {
	var imageId string
	cli, stdout, stdoutPipe := GetNewClient()
	go func() {
		err := cli.CmdImages()
		if err != nil {
			t.Fatalf("getImageId %s", err)
		}
		if err = CloseWrap(stdout, stdoutPipe); err != nil {
			t.Fatalf("getImageId %s", err)
		}
	}()
	imageId = PrintToStdout(t, stdout, stdoutPipe, repo)
	return strings.Fields(imageId)[2]
}

func RunEtcdTest(t *testing.T, uid string) {
	cli, stdout, stdoutPipe := GetNewClient()
	done := make(chan bool, 1)
	done1 := make(chan bool, 1)
	done2 := make(chan bool, 1)
	var imageId string
	var imageTag string
	go func() {
		fmt.Println("inside pull etcd")
		PullImage(t, cli, "phife.atribecalledchris.com:5000/deis/etcd:0.3.0")
		done <- true
	}()
	go func() {
		<-done
		fmt.Println("inside getting imageId")
		imageId = GetImageId(t, "phife.atribecalledchris.com:5000/deis/etcd")
		imageTag = "deis/etcd:" + uid
		cli.CmdTag(imageId, imageTag)
		done1 <- true
	}()
	go func() {
		<-done1
		done2 <- true
		fmt.Println("inside run etcd")
		RunContainer(t, cli, "--name", "deis-etcd-"+uid, imageTag)
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
