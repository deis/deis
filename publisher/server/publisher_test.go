package server

import (
	"strings"
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

func TestDomainNameToSkyDNS(t *testing.T) {
	s := &Server{nil, nil}
	domain := "deis.local"
	if !strings.EqualFold(s.DomainNameToSkyDNS(domain), "local/deis") {
		t.Errorf("%s should be local/deis", domain)
	}

	longDomain := "this.is.a.large.domain.name"
	if !strings.EqualFold(s.DomainNameToSkyDNS(longDomain), "name/domain/large/a/is/this") {
		t.Errorf("%s should be name/domain/large/a/is/this", longDomain)
	}
}
