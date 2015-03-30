// +build integration

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/deis/deis/tests/utils"
)

var (
	certsAddCmd      = "certs:add {{.SSLCertificatePath}} {{.SSLKeyPath}}"
	certsRemoveCmd   = "certs:remove {{.AppDomain}}"
	domainsAddCmd    = "domains:add {{.AppDomain}} --app {{.AppName}}"
	domainsRemoveCmd = "domains:remove {{.AppDomain}} --app {{.AppName}}"
)

func TestDomains(t *testing.T) {
	cfg := domainsSetup(t)
	domainsTest(t, cfg)
	certsTest(t, cfg)
	utils.AppsDestroyTest(t, cfg)
}

func domainsSetup(t *testing.T) *utils.DeisTestConfig {
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

func domainsTest(t *testing.T, cfg *utils.DeisTestConfig) {
	utils.Execute(t, domainsAddCmd, cfg, false, "done")
	// ensure both the root domain and the custom domain work
	utils.CurlApp(t, *cfg)
	utils.Curl(t, fmt.Sprintf("http://%s", cfg.AppDomain))
	utils.Execute(t, domainsRemoveCmd, cfg, false, "done")
	// only the root domain should work now
	utils.CurlApp(t, *cfg)
	// TODO (bacongobbler): add test to ensure that the custom domain fails to connect
}

func certsTest(t *testing.T, cfg *utils.DeisTestConfig) {
	utils.Execute(t, domainsAddCmd, cfg, false, "done")
	utils.Execute(t, certsAddCmd, cfg, false, cfg.AppDomain)
	// wait for the certs to be populated in the router; cron takes up to 1 minute
	fmt.Println("sleeping for 60 seconds until certs are generated...")
	time.Sleep(60 * time.Second)
	fmt.Println("ok")
	// ensure the custom domain's SSL endpoint works
	utils.Curl(t, fmt.Sprintf("https://%s", cfg.AppDomain))
	utils.Execute(t, certsRemoveCmd, cfg, false, "done")
	// only the root domain should work now
	utils.CurlApp(t, *cfg)
	// TODO (bacongobbler): add test to ensure that the custom domain fails to connect
}
