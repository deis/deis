// +build integration

package tests

import (
  "testing"

  "github.com/deis/deis/tests/utils"
)

var (
  CertsAddCmd     = "certs:add {{.AppCertPath}} {{.AppCertKeyPath}} --app={{.AppName}}"
  CertsUpdateCmd  = "certs:update {{.AppCertCn}} {{.AppCertPath}} {{.KeyCertKeyPath}} --app={{.AppName}}"
  CertsListCmd    = "certs:list --app={{.AppName}}"
  CertsRemoveCmd  = "certs:remove {{.AppCertCn}} --app={{.AppName}}"
)

func TestCerts(t *testing.T) {
  params := certsSetup(t)
  CertsAddTest(t, params)
  CertsListTest(t, params, params.AppCertCn)
  CertsRemoveTest(t, params)
  CertsListTest(t, params, "no certificates")
}

// Requires that ~/deis-test.cer and ~/deis-test.key be set up
// see tests/bin/test-setup.sh
func certsSetup(t *testing.T) *utils.DeisTestConfig {
  cfg := utils.GetGlobalConfig()
  cfg.AppName = "certsample"
  utils.Execute(t, authLoginCmd, cfg, false, "")
  utils.Execute(t, gitCloneCmd, cfg, false, "")
  if err := utils.Chdir(cfg.ExampleApp); err != nil {
    t.Fatal(err)
  }
  utils.Execute(t, appsCreateCmd, cfg, false, "")
  return cfg
}

func CertsAddTest(t *testing.T, params *utils.DeisTestConfig) {
  utils.Execute(t, CertsAddCmd, params, false, params.AppCertCn)
  utils.Execute(t, CertsAddCmd, params, true,
    "already exists")
}

func CertsUpdateTest(t *testing.T, params *utils.DeisTestConfig) {
  utils.Execute(t, CertsUpdateCmd, params, false, "updated")
}

func CertsListTest(t *testing.T, params *utils.DeisTestConfig, expect string) {
  utils.Execute(t, CertsListCmd, params, true, expect)
}

func CertsRemoveTest(t *testing.T, params *utils.DeisTestConfig) {
  utils.Execute(t, CertsRemoveCmd, params, false, "removed")
}
