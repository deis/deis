// +build integration

package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/deis/deis/tests/utils"
)

func TestDenyPathsClusterWide(t *testing.T) {
	client := utils.HTTPClient()
	cfg := denyPathsSetup()
	urls := []string{
		fmt.Sprintf("http://%s.%s/varz", cfg.AppName, cfg.Domain),
		fmt.Sprintf("http://%s.%s/statusz", cfg.AppName, cfg.Domain),
	}

	utils.Execute(t, "config:set DENY_PATHS=/varz,/statusz", cfg, false, "/")

	for _, url := range urls {
		response, err := client.Get(url)
		if err != nil {
			t.Fatalf("could not retrieve response from %s: %s", url, err)
		}
		defer response.Body.Close()
		if want, got := http.StatusForbidden, response.StatusCode; want != got {
			t.Errorf("GET %s: want %s got %s", url, http.StatusText(want), http.StatusText(got))
		}
	}
}

func denyPathsSetup() {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "denypathscluster"
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
