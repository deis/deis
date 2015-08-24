package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"

	"io/ioutil"
	"os"
)

func TestBuilder(t *testing.T) {
	var err error
	var errfile error
	setkeys := []string{
		"/deis/registry/protocol",
		"/deis/registry/host",
		"/deis/registry/port",
		"/deis/cache/host",
		"/deis/cache/port",
		"/deis/controller/protocol",
		"/deis/controller/host",
		"/deis/controller/port",
		"/deis/controller/builderKey",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/cache",
		"/deis/database",
		"/deis/registry",
		"/deis/domains",
		"/deis/services",
	}
	setproxy := []byte("HTTP_PROXY=\nhttp_proxy=\n")

	tmpfile, errfile := ioutil.TempFile("/tmp", "deis-test-")
	if errfile != nil {
		t.Fatal(errfile)
	}
	ioutil.WriteFile(tmpfile.Name(), setproxy, 0644)
	defer os.Remove(tmpfile.Name())

	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	imageName := utils.ImagePrefix() + "builder:" + tag
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	handler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
	etcdutils.PublishEtcd(t, handler)
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run %s at %s:%s\n", imageName, host, port)

	// Run a mock registry to test whether the builder can push its initial
	// images.
	regport := utils.RandomPort()
	mockRegistry(host, regport, t)
	setupRegistry("http", host, regport, t, handler)
	// When we switch to Registry v2, we probably want to uncomment this
	// and then remove mockRegistry.
	// dockercli.RunTestRegistry(t, "registry", host, regport)

	name := "deis-builder-" + tag
	defer cli.CmdRm("-f", "-v", name)
	go func() {
		_ = cli.CmdRm("-f", "-v", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":2223",
			"-e", "PORT=2223",
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "EXTERNAL_PORT="+port,
			"--privileged",
			"-v", tmpfile.Name()+":/etc/environment_proxy",
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "Builder is running")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "tcp")
	etcdutils.VerifyEtcdValue(t, "/deis/builder/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/builder/port", port, etcdPort)
}

// mockRegistry mocks a Docker v1 registry.
//
// This is largely derived from the Docker repo's mock:
// https://github.com/docker/docker/blob/84e917b8767c749b9bd1400a5a2253d972635bcf/registry/registry_mock_test.go
func mockRegistry(host, port string, t *testing.T) {
	addr := host + ":" + port

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Mock Registry request: %s %s\n", r.Method, r.RequestURI)

		if strings.Contains(r.RequestURI, "v2") {
			fmt.Println("**WARNING: mockRegistry does not support the v2 API**")
		}

		switch r.RequestURI {
		case "/v2/", "/v1/_ping":
			w.WriteHeader(200)
		case "/v1/repositories/deis/slugrunner/", "/v1/repositories/deis/slugbuilder/":
			w.Header().Add("X-Docker-Endpoints", addr)
			w.Header().Add("X-Docker-Token", fmt.Sprintf("FAKE-SESSION-%d", time.Now().UnixNano()))
			w.WriteHeader(200)
		case "/v1/repositories/deis/slugrunner/images", "/v1/repositories/deis/slugbuilder/images":
			w.WriteHeader(204)
		default:
			w.Header().Add("X-Docker-Size", "2000")
			w.WriteHeader(200)
		}

	})

	fmt.Printf("Starting mock registry on %s\n", addr)
	go http.ListenAndServe(addr, nil)
}

func setupRegistry(proto, host, port string, t *testing.T, handler *etcdutils.EtcdHandle) {
	vals := map[string]string{
		"/deis/registry/protocol": proto,
		"/deis/registry/port":     port,
		"/deis/registry/host":     host,
	}

	for k, v := range vals {
		fmt.Printf("Setting etcd key %s to %s\n", k, v)
		if _, err := handler.C.Set(k, v, 0); err != nil {
			t.Fatalf("Error setting %s to %s: %s\n", k, v, err)
		}
	}

}
