package procfile

import (
	"testing"
)

func TestParseBasic(t *testing.T) {
	basicProcfile := "web: bundle exec rackup"
	procs := Parse(basicProcfile)

	if procs["web"].Command != "bundle" {
		t.Error("Expected command to be bundle but got ", procs["web"].Command)
	}

	if len(procs["web"].Arguments) != 2 {
		t.Error("Expected arguments to be 2 in length, got ",
			len(procs["web"].Arguments))
	}

}

func TestParseMultiProcess(t *testing.T) {
	multilineProcfile := `web: coffee server.coffee
worker: node worker.js`

	procs := Parse(multilineProcfile)

	if procs["web"].Command != "coffee" {
		t.Error("Expected web command to be coffee")
	}

	if procs["worker"].Command != "node" {
		t.Error("Expected worker command to be node")
	}
}
