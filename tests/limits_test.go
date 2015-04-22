// +build integration

package tests

import (
	"strings"
	"testing"

	"github.com/deis/deis/tests/utils"
)

var (
	limitsListCmd     = "limits:list --app={{.AppName}}"
	limitsSetMemCmd   = "limits:set --app={{.AppName}} web=256MB"
	limitsSetCPUCmd   = "limits:set --app={{.AppName}} -c web=512"
	limitsUnsetMemCmd = "limits:unset --app={{.AppName}} --memory web"
	limitsUnsetCPUCmd = "limits:unset --app={{.AppName}} -c web"
	output1           = `(?s)"CpuShares": 512,.*"Memory": 0,`
	output2           = `(?s)"CpuShares": 512,.*"Memory": 268435456,`
	output3           = `(?s)"CpuShares": 0,.*"Memory": 268435456,`
	output4           = `(?s)"CpuShares": 0,.*"Memory": 0,`
)

func limitsSetTest(t *testing.T, cfg *utils.DeisTestConfig, ver int) {
	cpuCmd, memCmd := limitsSetCPUCmd, limitsSetMemCmd
	// regression test for https://github.com/deis/deis/issues/1563
	// previously the client would throw a stack trace with empty limits
	utils.Execute(t, limitsListCmd, cfg, false, "Unlimited")
	if strings.Contains(cfg.ExampleApp, "dockerfile") {
		cpuCmd = strings.Replace(cpuCmd, "web", "cmd", 1)
		memCmd = strings.Replace(memCmd, "web", "cmd", 1)
	}
	utils.Execute(t, cpuCmd, cfg, false, "512")
	utils.Execute(t, limitsListCmd, cfg, false, "512")
	utils.Execute(t, memCmd, cfg, false, "256M")
	utils.Execute(t, limitsListCmd, cfg, false, "256M")
}

func limitsUnsetTest(t *testing.T, cfg *utils.DeisTestConfig, ver int) {
	cpuCmd, memCmd := limitsUnsetCPUCmd, limitsUnsetMemCmd
	if strings.Contains(cfg.ExampleApp, "dockerfile") {
		cpuCmd = strings.Replace(cpuCmd, "web", "cmd", 1)
		memCmd = strings.Replace(memCmd, "web", "cmd", 1)
	}
	utils.Execute(t, cpuCmd, cfg, false, "Unlimited")
	utils.Execute(t, limitsListCmd, cfg, false, "Unlimited")
	utils.Execute(t, memCmd, cfg, false, "Unlimited")
	utils.Execute(t, limitsListCmd, cfg, false, "Unlimited")
}
