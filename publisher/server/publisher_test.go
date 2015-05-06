package server

import (
	"net"
	"testing"
)

func TestIsPublishableApp(t *testing.T) {
	s := &Server{}
	appName := "go_v2.web.1"
	if !s.IsPublishableApp(appName) {
		t.Errorf("%s should be publishable", appName)
	}
	badAppName := "go_v2"
	if s.IsPublishableApp(badAppName) {
		t.Errorf("%s should not be publishable", badAppName)
	}
	// publisher assumes that an app name of "test" with a null etcd client has v3 running
	oldVersion := "ceci-nest-pas-une-app_v2.web.1"
	if s.IsPublishableApp(oldVersion) {
		t.Errorf("%s should not be publishable", oldVersion)
	}
	currentVersion := "ceci-nest-pas-une-app_v3.web.1"
	if !s.IsPublishableApp(currentVersion) {
		t.Errorf("%s should be publishable", currentVersion)
	}
	futureVersion := "ceci-nest-pas-une-app_v4.web.1"
	if !s.IsPublishableApp(futureVersion) {
		t.Errorf("%s should be publishable", futureVersion)
	}
}

func TestIsPortOpen(t *testing.T) {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	defer ln.Close()

	s := &Server{}
	if !s.IsPortOpen(ln.Addr().String()) {
		t.Errorf("Port should be open")
	}
	if s.IsPortOpen("127.0.0.1:-1") {
		t.Errorf("Port should be closed")
	}
}
