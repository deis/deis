package integration

import (
    "errors"
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "testing"
)

const HOSTNAME string = "local.deisapp.com"

type App struct {
    Name        string // the name of the application
    RepoURL     string // URL to the app's git repository
}

func TestDeployRuby(t *testing.T) {
    err := deployApp(&App{
        Name: "deis-test-1",
        RepoURL: "https://github.com/deis/example-ruby-sinatra",
    })

    if err != nil {
        t.Errorf("failed deploying example-ruby-sinatra: %v", err)
    }
}

func TestDeployGo(t *testing.T) {
    err := deployApp(&App{
        Name: "deis-test-2",
        RepoURL: "https://github.com/deis/example-go",
    })

    if err != nil {
        t.Errorf("failed deploying example-go: %v", err)
    }
}

func TestDeployPython(t *testing.T) {
    err := deployApp(&App{
        Name: "deis-test-3",
        RepoURL: "https://github.com/deis/example-python-django",
    })

    if err != nil {
        t.Errorf("failed deploying example-python-django: %v", err)
    }
}

func TestDeployDockerfile(t * testing.T) {
    err := deployApp(&App{
        Name: "deis-test-4",
        RepoURL: "https://github.com/deis/example-dockerfile-python",
    })

    if err != nil {
        t.Errorf("failed deploying example-dockerfile-python: %v", err)
    }
}

func deployApp(app *App) error {
    var err error
    var output []byte

    err = cloneApp(app.RepoURL, app.Name)

    if err != nil {
        return errors.New("Could not clone repository " + app.RepoURL)
    }

    // switch context to the app repo
    os.Chdir(app.Name)

    output, _ = exec.Command(
        "deis",
        "apps:create",
        app.Name,
    ).CombinedOutput()

    if !strings.Contains(string(output), "done, created " + app.Name) {
        return errors.New("Received unexpected output for `deis apps:create`: " +
            string(output))
    }

    output, _ = exec.Command(
        "git",
        "push",
        "deis",
        "master",
    ).CombinedOutput()

    if !strings.Contains(string(output), "deployed to Deis") {
       return errors.New(
            "Received unexpected output for `git push deis master`: '%s'" +
            string(output))
    }

    _, err = http.Get(fmt.Sprintf("http://%s.%s", app.Name, HOSTNAME))

    if err != nil {
        return errors.New(fmt.Sprintf("Could not reach %s.%s: '%s'", app.Name, HOSTNAME, output))
    }

    err = destroyApp(app.Name)

    // remove the test_app directory
    os.Chdir("..")
    os.RemoveAll(app.Name)

    return err
}

func cloneApp(url string, path string) error {
    // clone the example app to the test_app directory
    _, err := exec.Command("git", "clone", url, path).CombinedOutput()
    return err
}

func destroyApp(name string) error {
    output, err := exec.Command(
        "deis",
        "apps:destroy",
        fmt.Sprintf("--app=%s", name),
        fmt.Sprintf("--confirm=%s", name),
    ).CombinedOutput()

    if !strings.Contains(string(output), "done") {
        return errors.New(
            "Received unexpected output for `deis apps:destroy`: " + string(output))
    }
    return err
}
