package server

import (
	"testing"
)

func TestIsPublishableApp(t *testing.T) {
	s := &Server{nil, nil}
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
