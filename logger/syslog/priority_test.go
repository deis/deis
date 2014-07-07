package syslog

import (
    "testing"
)

func TestFacilityToString(t *testing.T) {
    var fac Facility = 100
    if fac.String() != "unknown" {
        t.Errorf("facility != unknown; got %s", fac.String())
    }
    fac = 20
    if fac.String() != "local4" {
        t.Errorf("facility != local4; got %s", fac.String())
    }
}

func TestSeverityToString(t *testing.T) {
    var sev Severity = 100
    if sev.String() != "unknown" {
        t.Errorf("severity != unknown; got %s", sev.String())
    }
    sev = 5
    if sev.String() != "notice" {
        t.Errorf("severity != local4; got %s", sev.String())
    }
}
