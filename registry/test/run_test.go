package verbose

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dotcloud/docker/api/client"
)

const (
	unitTestImageName        = "docker-test-image"
	unitTestImageID          = "83599e29c455eb719f77d799bc7c51521b9551972f5a850d7ad265bc1b5292f6" // 1.0
	unitTestImageIDShort     = "83599e29c455"
	unitTestNetworkBridge    = "testdockbr0"
	unitTestStoreBase        = "/var/lib/docker/unit-tests"
	testDaemonHttpsProto     = "tcp"
	testDaemonHttpsAddr      = "localhost:4271"
	testDaemonRogueHttpsAddr = "localhost:4272"
)

func daemonAddr() string {
	addr := os.Getenv("TEST_DAEMON_ADDR")
	if addr == "" {
		addr = "/var/run/docker.sock"
	}
	return addr
}

func daemonProto() string {
	proto := os.Getenv("TEST_DAEMON_PROTO")
	if proto == "" {
		proto = "unix"
	}
	return proto
}

func closeWrap(args ...io.Closer) error {
	e := false
	ret := fmt.Errorf("error closing elements")
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

func getIPAddressTest(t *testing.T) string {
	stdin, _ := io.Pipe()
	stdout, stdoutPipe := io.Pipe()
	cli := client.NewDockerCli(nil, stdoutPipe, nil, daemonProto(), daemonAddr(), nil)
	var IPAdress string
	go func() {
		err := cli.CmdInspect("--format", "{{ .NetworkSettings.IPAddress }}", "deis-etcd")
		if err != nil {
			t.Fatalf("getIPAddressTest %s", err)
		}
		if err = closeWrap(stdout, stdoutPipe, stdin); err != nil {
			t.Fatalf("getIPAddressTest %s", err)
		}
	}()
	time.Sleep(2000 * time.Millisecond)
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			IPAdress = cmdBytes
			fmt.Println(cmdBytes)
		} else {
			break
		}
		fmt.Println("get IPAddress")
	}
	return IPAdress
}

func pullEtcdTest(t *testing.T, done chan bool) {
	stdin, _ := io.Pipe()
	stdout, stdoutPipe := io.Pipe()
	fmt.Println("1")
	cli := client.NewDockerCli(nil, stdoutPipe, nil, daemonProto(), daemonAddr(), nil)
	fmt.Println("2")
	go func() {
		err := cli.CmdPull("phife.atribecalledchris.com:5000/deis/etcd:0.3.0")
		if err != nil {
			t.Fatalf("pullEtcdTest %s", err)
		}
		if err = closeWrap(stdout, stdoutPipe, stdin); err != nil {
			t.Fatalf("pullEtcdTest %s", err)
		}
	}()
	time.Sleep(3000 * time.Millisecond)
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			fmt.Println(cmdBytes)
		} else {
			break
		}
		fmt.Println("pulling etcd")
	}
	done <- true

}

func runEtcdTest(t *testing.T, done chan bool) {
	stdin, _ := io.Pipe()
	stdout, stdoutPipe := io.Pipe()
	cli := client.NewDockerCli(nil, stdoutPipe, nil, daemonProto(), daemonAddr(), nil)
	go func() {
		fmt.Println("12")
		err := cli.CmdRun("--name", "deis-etcd", "phife.atribecalledchris.com:5000/deis/etcd:0.3.0")
		if err != nil {
			t.Fatalf("runEtcdTest %s", err)
		}

	}()
	go func() {
		fmt.Println("here")
		time.Sleep(5000 * time.Millisecond)
		if err := closeWrap(stdout, stdoutPipe, stdin); err != nil {
			t.Fatalf("runEtcdTest %s", err)
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			fmt.Println(cmdBytes)
		} else {
			break
		}
		fmt.Println("Running Etcd ")
	}
	done <- true
}

func buildRegistryTest(t *testing.T) {
	stdin, _ := io.Pipe()
	stdout, stdoutPipe := io.Pipe()
	cli := client.NewDockerCli(nil, stdoutPipe, nil, daemonProto(), daemonAddr(), nil)
	go func() {
		err := cli.CmdBuild("../")
		if err != nil {
			t.Fatalf("buildRegistryTest %s", err)
		}
		if err = closeWrap(stdout, stdoutPipe, stdin); err != nil {
			t.Fatalf("buildRegistryTest %s", err)
		}
	}()
	time.Sleep(3000 * time.Millisecond)
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			fmt.Println(cmdBytes)
		} else {
			break
		}
		fmt.Println("building Deis registy Dockerfile")
	}
}

func runDeisRegistryDataTest(t *testing.T) {
	stdin, _ := io.Pipe()
	stdout, stdoutPipe := io.Pipe()

	cli := client.NewDockerCli(nil, stdoutPipe, nil, daemonProto(), daemonAddr(), nil)

	go func() {
		err := cli.CmdInspect("--format", "'{{ .Config.Hostname }}'", "deis-registry-data")
		/*if err != nil {
			t.Fatalf("runDeisRegistryDataTest %s",err)
		}*/
		if err = closeWrap(stdout, stdoutPipe, stdin); err != nil {
			t.Fatalf("runDeisRegistryDataTest %s", err)
		}
	}()
	go func() {
		//fmt.Println("here1")
		time.Sleep(2000 * time.Millisecond)
		if err := closeWrap(stdout, stdoutPipe, stdin); err != nil {
			t.Fatalf("runEtcdTest %s", err)
		}
	}()
	var hostname string
	time.Sleep(1000 * time.Millisecond)
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			fmt.Println(cmdBytes)
			hostname = cmdBytes
		} else {
			break
		}
		fmt.Println("inspecting deis registry data")
	}

	if strings.Contains(hostname, "Error") == true {
		go func() {
			err := cli.CmdRun("--name", "deis-registry-data", "-v", "/data", "deis/base", "/bin/true")
			if err != nil {
				t.Fatalf("%s", err)
			}
			if err := closeWrap(stdout, stdoutPipe, stdin); err != nil {
				t.Fatalf("%s", err)
			}
		}()
		for {
			if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
				fmt.Println(cmdBytes)
			} else {
				break
			}
			fmt.Println("pulling Deis registry data")
		}
	}
}

func runDeisRegistryTest(t *testing.T, IPAddress string) {
	stdin, _ := io.Pipe()
	stdout, stdoutPipe := io.Pipe()
	cli := client.NewDockerCli(nil, stdoutPipe, nil, daemonProto(), daemonAddr(), nil)
	go func() {
		//docker run --name deis-registry -p 5000:5000 -e PUBLISH=5000 -e HOST=10.0.0.37 --volumes-from deis-registry-data deis/registry

		err := cli.CmdRun("--name", "deis-registry", "-p", "5000:5000", "-e", "PUBLISH=5000", "-e", "HOST="+IPAddress, "--volumes-from", "deis-registry-data", "deis/registry")
		if err != nil {
			t.Fatalf("runDeisRegistryTest %s", err)
		}
		if err := closeWrap(stdout, stdoutPipe, stdin); err != nil {
			t.Fatalf("runDeisRegistryTest %s", err)
		}
	}()
	time.Sleep(3000 * time.Millisecond)
	for {
		if cmdBytes, err := bufio.NewReader(stdout).ReadString('\n'); err == nil {
			if strings.Contains(cmdBytes, "Booting") == true {
				if err := closeWrap(stdout, stdoutPipe, stdin); err != nil {
					t.Fatalf("runDeisRegistryTest %s", err)
				}
			}
			fmt.Println(cmdBytes)
		} else {
			break
		}
		fmt.Println("pulling Deis Registry")
	}

}

func TestBuild(t *testing.T) {
	done := make(chan bool, 1)
	done2 := make(chan bool, 1)
	pullEtcdTest(t, done)
	fmt.Println("1st")
	go func() {
		fmt.Println("2nd")
		<-done
		fmt.Println("3rd")
		runEtcdTest(t, done2)
		fmt.Println("4th")
	}()
	<-done2
	fmt.Println("5th")
	buildRegistryTest(t)
	fmt.Println("6th")
	runDeisRegistryDataTest(t)
	fmt.Println("7th")
	IPAddress := strings.TrimSuffix(getIPAddressTest(t), "\n")
	fmt.Println(IPAddress)
	fmt.Println("8th")
	runDeisRegistryTest(t, IPAddress)
}
