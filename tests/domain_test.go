// +build integration

package tests

import (
	"fmt"
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	domainsAddCmd = "domains:add {{.AppDomain}} --app {{.AppName}}"
	domainsRemoveCmd = "domains:remove {{.AppDomain}} --app {{.AppName}}"
)

func TestDomain(t *testing.T) {
	cfg := domainSetup(t)
	domainTest(t, cfg)
	utils.AppsDestroyTest(t, cfg)
}

func domainSetup(t *testing.T) *utils.DeisTestConfig {
	cfg := utils.GetGlobalConfig()
	cfg.AppName = "domainsample"
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

func domainTest(t *testing.T, cfg *utils.DeisTestConfig) {
	utils.Execute(t, domainsAddCmd, cfg, false, "done")
	// ensure both the root domain and the custom domain work
	utils.CurlApp(t, *cfg)
	utils.Curl(t, fmt.Sprintf("http://%s", cfg.AppDomain))
	utils.Execute(t, domainsRemoveCmd, cfg, false, "done")
	// only the root domain should work now
	utils.CurlApp(t, *cfg)
	// TODO (bacongobbler): add test to ensure that the custom domain fails to connect
}
