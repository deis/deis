package integration

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "path"
    "strings"
    "testing"

    "github.com/deis/deis/version"
)

var (
    AuthKey string = path.Join(os.Getenv("HOME"), ".ssh", "testkey")
    ExampleAppURL string = "https://github.com/deis/example-ruby-sinatra"
    Hosts string = "172.17.8.100"
    Hostname string = "local.deisapp.com"
    SSHKey string = "~/.vagrant.d/insecure_private_key"
)

func TestSetUp(t *testing.T) {
    // initialize variables if the corresponding envvar exists
    if key := os.Getenv("DEIS_TEST_APP_URL"); key != "" {
        ExampleAppURL = key
    }
    if key := os.Getenv("DEIS_TEST_HOSTS"); key != "" {
        Hosts = key
    }
    if key := os.Getenv("DEIS_TEST_HOSTNAME"); key != "" {
        Hostname = key
    }
    if key := os.Getenv("DEIS_TEST_SSH_KEY"); key != "" {
        SSHKey = key
    }

    _, err := os.Stat(AuthKey)
    if os.IsNotExist(err) {
        // generate a new SSH key
        output, err := exec.Command(
            "ssh-keygen",
            "-q",
            "-t",
            "rsa",
            "-f",
            AuthKey,
            "-N",
            "",
            "-C",
            "deis",
        ).CombinedOutput()
        if err != nil {
            t.Fatalf("Error while generating SSH key: %s", string(output))
        }
    }

    // add the SSH key to the keychain
    output, err := exec.Command(
        "ssh-add",
        AuthKey,
    ).CombinedOutput()
    if err != nil {
        t.Fatalf("Error while adding SSH key to the keychain: %s", string(output))
    }
}

func TestClientVersion(t *testing.T) {
    cmd := exec.Command("deis", "--version")
    output, _ := cmd.CombinedOutput()

    if strings.TrimSpace(string(output)) != fmt.Sprintf("%s", version.Version) {
        t.Fatalf("Received unexpected output for `deis --version`: '%s'", output)
    }
}

func TestRegister(t *testing.T) {
    cmd := exec.Command(
        "deis",
        "register",
        fmt.Sprintf("http://%s:8000", Hostname),
        "--username=test",
        "--email=test@test.co.nz",
        "--password=asdf1234",
    )
    output, _ := cmd.CombinedOutput()

    // Account for the username already being registered
    if !strings.Contains(string(output), "Regist") {
       t.Fatalf("Received unexpected output for `deis register`: '%s'", output)
    }
}

func TestLogin(t *testing.T) {
    cmd := exec.Command(
        "deis",
        "login",
        fmt.Sprintf("http://%s:8000", Hostname),
        "--username=test",
        "--password=asdf1234",
    )
    output, _ := cmd.CombinedOutput()

    if !strings.Contains(string(output), "Logged in as test") {
       t.Fatalf("Received unexpected output for `deis login`: '%s'", output)
    }
}

func TestAddKey(t *testing.T) {
    cmd := exec.Command(
        "deis",
        "keys:add",
        fmt.Sprintf("%s.pub", AuthKey),
    )
    output, _ := cmd.CombinedOutput()

    if !strings.Contains(string(output), "done") &&
       !strings.Contains(string(output), "already exists") {
       t.Fatalf("Received unexpected output for `deis keys:add`: '%s'", output)
    }
}

func TestCreateCluster(t *testing.T) {
    cmd := exec.Command(
        "deis",
        "clusters:create",
        "dev",
        Hostname,
        fmt.Sprintf("--hosts=%s", Hosts),
        fmt.Sprintf("--auth=%s", SSHKey),
    )
    output, _ := cmd.CombinedOutput()

    if !strings.Contains(string(output), "done, created dev") &&
       !strings.Contains(string(output), "already exists") {
       t.Fatalf("Received unexpected output for `deis clusters:create`: '%s'", output)
    }
}

func TestAppDeploy(t *testing.T) {
    var err error
    var output []byte

    // clone the example app to the test_app directory
    _, err = exec.Command("git", "clone", ExampleAppURL, "test_app").CombinedOutput()

    if err != nil {
        t.Fatalf("Error while cloning app: %v", err)
    }

    err = os.Chdir("test_app")

    if err != nil {
        t.Fatalf("Could not change directories: '%s'", err)
    }

    output, _ = exec.Command(
        "deis",
        "apps:create",
    ).CombinedOutput()

    if !strings.Contains(string(output), "done, created") {
       t.Fatalf("Received unexpected output for `deis apps:create`: '%s'", output)
    }

    sl := strings.Split(strings.Replace(string(output), "\n", " ", -1), " ")
    appName := sl[len(sl)-6]

    output, _ = exec.Command(
        "git",
        "push",
        "deis",
        "master",
    ).CombinedOutput()

    if !strings.Contains(string(output), "deployed to Deis") {
       t.Fatalf("Received unexpected output for `git push deis master`: '%s'", output)
    }

    output, _ = exec.Command(
        "deis",
        "scale",
        "web=2",
    ).CombinedOutput()

    if !strings.Contains(string(output), "done") {
       t.Fatalf("Received unexpected output for `deis scale`: '%s'", output)
    }

    output, _ = exec.Command(
        "deis",
        "domains:add",
        fmt.Sprintf("testdomain.%s", Hostname),
    ).CombinedOutput()

    if !strings.Contains(string(output), "done") {
       t.Fatalf("Received unexpected output for `deis domains:add`: '%s'", output)
    }

    _, err = http.Get(fmt.Sprintf("http://%s.%s", appName, Hostname))

    if err != nil {
        t.Fatalf("Could not reach %s.%s: '%s'", appName, Hostname, output)
    }

    _, err = http.Get(fmt.Sprintf("http://testdomain.%s", Hostname))

    if err != nil {
        t.Fatalf("Could not reach testdomain.%s: '%s'", Hostname, output)
    }

    output, _ = exec.Command(
        "deis",
        "apps:destroy",
        fmt.Sprintf("--confirm=%s", appName),
    ).CombinedOutput()

    if !strings.Contains(string(output), "done") {
       t.Fatalf("Received unexpected output for `deis apps:destroy`: '%s'", output)
    }

    // remove the test_app directory
    os.Chdir("..")
    err = os.RemoveAll("test_app")

    if err != nil {
        t.Fatalf("Could not remove test_app: '%s'", err)
    }
}

func TestDestroyCluster(t *testing.T) {
    cmd := exec.Command(
        "deis",
        "clusters:destroy",
        "dev",
        "--confirm=dev",
    )
    output, _ := cmd.CombinedOutput()

    if !strings.Contains(string(output), "done") {
       t.Fatalf("Received unexpected output for `deis clusters:destroy`: '%s'", output)
    }
}

func TestRemoveKey(t *testing.T) {
    cmd := exec.Command(
        "deis",
        "keys:remove",
        "deis",
    )
    output, _ := cmd.CombinedOutput()

    if !strings.Contains(string(output), "done") {
       t.Fatalf("Received unexpected output for `deis keys:remove`: '%s'", output)
    }
}

func TestLogout(t *testing.T) {
    cmd := exec.Command(
        "deis",
        "auth:logout",
    )
    output, _ := cmd.CombinedOutput()

    if strings.TrimSpace(string(output)) != "Logged out" {
       t.Fatalf("Received unexpected output for `deis auth:logout`: '%s'", output)
    }
}
