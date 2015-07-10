// +build integration

package tests

import (
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	healthcheckGoodCmd  = "config:set HEALTHCHECK_URL=/ HEALTHCHECK_STATUS_CODE=200 --app={{.AppName}}"
	healthcheckBadCmd   = "config:set HEALTHCHECK_URL=/ HEALTHCHECK_STATUS_CODE=999 --app={{.AppName}}"
	healthcheckUnsetCmd = "config:unset HEALTHCHECK_URL HEALTHCHECK_STATUS_CODE --app={{.AppName}}"
)

func TestHealthcheck(t *testing.T) {
	params := healthcheckSetup(t)
	psScaleTest(t, params, psScaleCmd)
	utils.Execute(t, healthcheckGoodCmd, params, false, "/")
	psScaleTest(t, params, psScaleCmd)
	psScaleTest(t, params, psDownScaleCmd)
	psScaleTest(t, params, psScaleCmd)
	utils.Execute(t, healthcheckUnsetCmd, params, false, "")
	// FIXME (bacongobbler): publisher will not template the app if the status code is incorrect, which means we'll receive
	// a 404 response from the router.
	utils.Execute(t, healthcheckBadCmd, params, true, "aborting, app failed health check (got '404', expected: '999')")
	utils.Execute(t, healthcheckGoodCmd, params, false, "/")
	appsOpenTest(t, params)
	utils.AppsDestroyTest(t, params)
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
