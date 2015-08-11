package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/config/model"
	"github.com/deis/deis/deisctl/test/mock"
	"github.com/deis/deis/deisctl/units"
)

type backendStub struct {
	startedUnits     []string
	stoppedUnits     []string
	installedUnits   []string
	uninstalledUnits []string
	restartedUnits   []string
	expected         bool
}

func (backend *backendStub) Create(targets []string, wg *sync.WaitGroup, out, ew io.Writer) {
	backend.installedUnits = append(backend.installedUnits, targets...)
}
func (backend *backendStub) Destroy(targets []string, wg *sync.WaitGroup, out, ew io.Writer) {
	backend.uninstalledUnits = append(backend.uninstalledUnits, targets...)
}
func (backend *backendStub) Start(targets []string, wg *sync.WaitGroup, out, ew io.Writer) {
	backend.startedUnits = append(backend.startedUnits, targets...)
}
func (backend *backendStub) Stop(targets []string, wg *sync.WaitGroup, out, ew io.Writer) {
	backend.stoppedUnits = append(backend.stoppedUnits, targets...)
}
func (backend *backendStub) Scale(component string, num int, wg *sync.WaitGroup, out, ew io.Writer) {
	switch {
	case component == "router" && num == 3:
		backend.expected = true
	case component == "registry" && num == 4:
		backend.expected = true
	default:
		backend.expected = false
	}
}
func (backend *backendStub) RollingRestart(target string, wg *sync.WaitGroup, out, ew io.Writer) {
	backend.restartedUnits = append(backend.restartedUnits, target)
}

func (backend *backendStub) ListUnits() error {
	return nil
}
func (backend *backendStub) ListUnitFiles() error {
	return nil
}
func (backend *backendStub) Status(target string) error {
	if target == "controller" || target == "builder" {
		return nil
	}
	return errors.New("Test Error")
}
func (backend *backendStub) Journal(target string) error {
	if target == "controller" || target == "builder" {
		return nil
	}
	return errors.New("Test Error")
}
func (backend *backendStub) SSH(target string) error {
	if target == "controller" {
		return nil
	}
	return errors.New("Error")
}
func (backend *backendStub) SSHExec(target, command string) error {
	if target == "controller" && command == "sh" {
		return nil
	}
	return errors.New("Error")
}

func (backend *backendStub) Dock(target string, command []string) error {
	return nil
}

var _ backend.Backend = &backendStub{}

func fakeCheckKeys(cb config.Backend) error {
	return nil
}

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if strings.Split(req.URL.Path, "/")[1] != "v1.7.2" {
		res.WriteHeader(http.StatusNotFound)
	}

	res.Write([]byte("test"))
}

func TestRefreshUnits(t *testing.T) {
	t.Parallel()

	name, err := ioutil.TempDir("", "deisctl")

	if err != nil {
		t.Error(err)
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	err = RefreshUnits(name, "v1.7.2", server.URL+"/")

	if err != nil {
		t.Error(err)
	}

	files, err := ioutil.ReadDir(name)

	// There will be a "decorators" subdirectory and that shouldn't be
	// counted as a unit when making the upcoming assertion.
	numFiles := len(files) - 1

	if len(units.Names) != numFiles {
		t.Error(fmt.Errorf("Expected %d units, Got %d", len(units.Names), numFiles))
	}

	for _, unit := range units.Names {
		found := false

		for _, file := range files {
			if unit+".service" == file.Name() {
				found = true
			}
		}

		if found == false {
			t.Error(fmt.Errorf("Expected to find %s in %v", unit, files))
		}
	}
}

func TestRefreshUnitsError(t *testing.T) {
	t.Parallel()

	name, err := ioutil.TempDir("", "deisctl")

	if err != nil {
		t.Error(err)
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	err = RefreshUnits(name, "foo", server.URL+"/")
	result := err.Error()
	expected := "404 Not Found"

	if result != expected {
		t.Error(fmt.Errorf("Expected %s, Got %s", expected, result))
	}
}

func TestListUnits(t *testing.T) {
	t.Parallel()

	b := backendStub{installedUnits: []string{"router@1", "router@2"}}

	if ListUnits(&b) != nil {
		t.Error("unexpected error")
	}
}

func TestListUnitFiles(t *testing.T) {
	t.Parallel()

	b := backendStub{}

	if ListUnitFiles(&b) != nil {
		t.Error("unexpected error")
	}
}

func TestScaling(t *testing.T) {
	t.Parallel()

	b := backendStub{expected: false}
	scale := []string{"registry=4", "router=3"}

	Scale(scale, &b)

	if b.expected == false {
		t.Error("b.Scale called with unexpected arguements")
	}
}

func TestScalingNonScalableComponent(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := "cannot scale controller component"
	err := Scale([]string{"controller=2"}, &b).Error()

	if err != expected {
		t.Error(fmt.Errorf("Expected '%v', Got '%v'", expected, err))
	}
}

func TestScalingInvalidFormat(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := "Could not parse: controller2"
	err := Scale([]string{"controller2"}, &b).Error()

	if err != expected {
		t.Error(fmt.Errorf("Expected '%v', Got '%v'", expected, err))
	}
}

func TestStart(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@1", "router@2"}

	Start(expected, &b)

	if !reflect.DeepEqual(b.startedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.startedUnits))
	}
}

func TestStartPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"store-monitor", "store-daemon", "store-metadata", "store-gateway@*",
		"store-volume", "logger", "logspout", "database", "registry@*", "controller",
		"builder", "publisher", "router@*", "database", "registry@*", "controller",
		"builder", "publisher", "router@*"}

	Start([]string{"platform"}, &b)

	if !reflect.DeepEqual(b.startedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.startedUnits))
	}
}

func TestStartStatelessPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"logspout", "registry@*", "controller",
		"builder", "publisher", "router@*", "registry@*", "controller",
		"builder", "publisher", "router@*"}

	Start([]string{"stateless-platform"}, &b)

	if !reflect.DeepEqual(b.startedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.startedUnits))
	}
}

func TestStartSwarm(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"swarm-manager", "swarm-node"}

	Start([]string{"swarm"}, &b)

	if !reflect.DeepEqual(b.startedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.startedUnits))
	}
}

func TestRollingRestart(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router"}

	RollingRestart("router", &b)

	if !reflect.DeepEqual(b.restartedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.restartedUnits))
	}
}

func TestUpgradePrep(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"database", "registry@*", "controller", "builder", "logger", "logspout", "store-volume",
		"store-gateway@*", "store-metadata", "store-daemon", "store-monitor"}

	UpgradePrep(&b)

	if !reflect.DeepEqual(b.stoppedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.stoppedUnits))
	}
}

func TestUpgradeTakeover(t *testing.T) {
	t.Parallel()
	testMock := mock.ConfigBackend{Expected: []*model.ConfigNode{{Key: "/deis/services/app1", Value: "foo", TTL: 10},
		{Key: "/deis/services/app2", Value: "8000", TTL: 10}}}

	b := backendStub{}
	expectedRestarted := []string{"router"}
	expectedStarted := []string{"publisher", "store-monitor", "store-daemon", "store-metadata",
		"store-gateway@*", "store-volume", "logger", "logspout", "database", "registry@*",
		"controller", "builder", "publisher", "router@*", "database", "registry@*",
		"controller", "builder", "publisher", "router@*"}

	if err := doUpgradeTakeOver(&b, testMock); err != nil {
		t.Error(fmt.Errorf("Takeover failed: %v", err))
	}

	if !reflect.DeepEqual(b.restartedUnits, expectedRestarted) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expectedRestarted, b.restartedUnits))
	}
	if !reflect.DeepEqual(b.startedUnits, expectedStarted) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expectedStarted, b.startedUnits))
	}
}

func TestStop(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@1", "router@2"}
	Stop(expected, &b)

	if !reflect.DeepEqual(b.stoppedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.stoppedUnits))
	}
}

func TestStopPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@*", "publisher", "controller", "builder", "database",
		"registry@*", "logger", "logspout", "store-volume", "store-gateway@*",
		"store-metadata", "store-daemon", "store-monitor"}
	Stop([]string{"platform"}, &b)

	if !reflect.DeepEqual(b.stoppedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.stoppedUnits))
	}
}

func TestStopStatelessPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@*", "publisher", "controller", "builder",
		"registry@*", "logspout"}
	Stop([]string{"stateless-platform"}, &b)

	if !reflect.DeepEqual(b.stoppedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.stoppedUnits))
	}
}

func TestStopSwarm(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"swarm-node", "swarm-manager"}
	Stop([]string{"swarm"}, &b)

	if !reflect.DeepEqual(b.stoppedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.stoppedUnits))
	}
}

func TestRestart(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@4", "router@5"}

	Restart(expected, &b)

	if !reflect.DeepEqual(b.stoppedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.stoppedUnits))
	}
	if !reflect.DeepEqual(b.startedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.startedUnits))
	}
}

func TestSSH(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	err := SSH("controller", []string{}, &b)

	if err != nil {
		t.Error(err)
	}
}
func TestSSHExec(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	err := SSH("controller", []string{"sh"}, &b)

	if err != nil {
		t.Error(err)
	}
}

func TestSSHError(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	err := SSH("registry", []string{}, &b)

	if err == nil {
		t.Error("Error expected")
	}
}

func TestStatus(t *testing.T) {
	t.Parallel()

	b := backendStub{}

	if Status([]string{"controller", "builder"}, &b) != nil {
		t.Error("Unexpected Error")
	}
}

func TestStatusError(t *testing.T) {
	t.Parallel()

	b := backendStub{}

	expected := "Test Error"
	err := Status([]string{"blah"}, &b).Error()

	if err != expected {
		t.Error(fmt.Errorf("Expected '%v', Got '%v'", expected, err))
	}
}

func TestJournal(t *testing.T) {
	t.Parallel()

	b := backendStub{}

	if Journal([]string{"controller", "builder"}, &b) != nil {
		t.Error("Unexpected Error")
	}
}

func TestJournalError(t *testing.T) {
	t.Parallel()

	b := backendStub{}

	expected := "Test Error"
	err := Journal([]string{"blah"}, &b).Error()

	if err != expected {
		t.Error(fmt.Errorf("Expected '%v', Got '%v'", expected, err))
	}
}

func TestInstall(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	cb := mock.ConfigBackend{}

	expected := []string{"router@1", "router@2"}

	Install(expected, &b, &cb, fakeCheckKeys)

	if !reflect.DeepEqual(b.installedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.installedUnits))
	}
}

func TestInstallPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	cb := mock.ConfigBackend{}

	expected := []string{"store-daemon", "store-monitor", "store-metadata", "store-volume",
		"store-gateway@1", "logger", "logspout", "database", "registry@1",
		"controller", "builder", "publisher", "router@1", "router@2", "router@3"}

	Install([]string{"platform"}, &b, &cb, fakeCheckKeys)

	if !reflect.DeepEqual(b.installedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.installedUnits))
	}
}

func TestInstallPlatformWithCustomRouterMeshSize(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	cb := mock.ConfigBackend{}

	expected := []string{"store-daemon", "store-monitor", "store-metadata", "store-volume",
		"store-gateway@1", "logger", "logspout", "database", "registry@1",
		"controller", "builder", "publisher", "router@1", "router@2", "router@3", "router@4", "router@5"}
	RouterMeshSize = 5

	Install([]string{"platform"}, &b, &cb, fakeCheckKeys)
	RouterMeshSize = DefaultRouterMeshSize

	if !reflect.DeepEqual(b.installedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.installedUnits))
	}
}

func TestInstallStatelessPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	cb := mock.ConfigBackend{}

	expected := []string{"logspout", "registry@1",
		"controller", "builder", "publisher", "router@1", "router@2", "router@3"}

	Install([]string{"stateless-platform"}, &b, &cb, fakeCheckKeys)

	if !reflect.DeepEqual(b.installedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.installedUnits))
	}
}

func TestInstallSwarm(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	cb := mock.ConfigBackend{}

	expected := []string{"swarm-manager", "swarm-node"}

	Install([]string{"swarm"}, &b, &cb, fakeCheckKeys)

	if !reflect.DeepEqual(b.installedUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.installedUnits))
	}
}

func TestUninstall(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@3", "router@4"}

	Uninstall(expected, &b)

	if !reflect.DeepEqual(b.uninstalledUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.uninstalledUnits))
	}
}

func TestUninstallPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@*", "publisher", "controller", "builder", "database",
		"registry@*", "logger", "logspout", "store-volume", "store-gateway@*",
		"store-metadata", "store-daemon", "store-monitor"}

	Uninstall([]string{"platform"}, &b)

	if !reflect.DeepEqual(b.uninstalledUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.uninstalledUnits))
	}
}

func TestUninstallStatelessPlatform(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"router@*", "publisher", "controller", "builder",
		"registry@*", "logspout"}

	Uninstall([]string{"stateless-platform"}, &b)

	if !reflect.DeepEqual(b.uninstalledUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.uninstalledUnits))
	}
}

func TestUninstallSwarm(t *testing.T) {
	t.Parallel()

	b := backendStub{}
	expected := []string{"swarm-node", "swarm-manager"}

	Uninstall([]string{"swarm"}, &b)

	if !reflect.DeepEqual(b.uninstalledUnits, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, b.uninstalledUnits))
	}
}
