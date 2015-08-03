// +build integration

package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	healthcheckGoodCmd = "config:set HEALTHCHECK_URL=/ --app={{.AppName}}"
)

func TestHealthcheck(t *testing.T) {
	client := utils.HTTPClient()
	cfg := healthcheckSetup(t)
	done := make(chan bool, 1)
	url := fmt.Sprintf("http://%s.%s", cfg.AppName, cfg.Domain)

	utils.Execute(t, healthcheckGoodCmd, cfg, false, "/")
	go func() {
		// there should never be any downtime during these health check operations
		psScaleTest(t, cfg, psScaleCmd)
		cfg.ProcessNum = "1"
		psScaleTest(t, cfg, psScaleCmd)
		// kill healthcheck goroutine
		done <- true
	}()

	// run health checks in parallel while performing operations
	fmt.Printf("starting health checks at %s\n", url)
loop:
	for {
		select {
		case <-done:
			fmt.Println("done performing health checks")
			break loop
		default:
			doHealthCheck(t, client, url)
		}
	}
	utils.AppsDestroyTest(t, cfg)
}

func healthcheckSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "healthchecksample"
	utils.Execute(t, authLoginCmd, cfg, false, "")
	utils.Execute(t, gitCloneCmd, cfg, false, "")
	if err := utils.Chdir(cfg.ExampleApp); err != nil {
		t.Fatal(err)
	}
	utils.Execute(t, appsCreateCmd, cfg, false, "")
	utils.Execute(t, gitPushCmd, cfg, false, "")
	utils.CurlApp(t, *cfg)
	if err := utils.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func doHealthCheck(t *testing.T, client *http.Client, url string) {
	response, err := client.Get(url)
	if err != nil {
		t.Fatalf("could not retrieve response from %s: %v\n", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("app had some downtime while undergoing health checks (got %d response)", response.StatusCode)
	}
}
