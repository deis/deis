package fleet

import (
	"strings"
	"testing"
)

func TestGetTunnelFlag(t *testing.T) {
	flag := getTunnelFlag()
	if flag != "" {
		t.Fatalf("got %v, expected \"\"", flag)
	}
	Flags.Tunnel = "hostname:2222"
	flag = getTunnelFlag()
	if flag != "hostname:2222" {
		t.Fatalf("got %v, expected \"hostname:2222\"", flag)
	}
	Flags.Tunnel = "hostname"
	flag = getTunnelFlag()
	if flag != "hostname:22" {
		t.Fatalf("got %v, expected \"hostname:22\"", flag)
	}
}

func TestGetChecker(t *testing.T) {
	Flags.StrictHostKeyChecking = false
	checker := getChecker()
	if checker != nil {
		t.Fatalf("expected nil checker, got %v", checker)
	}
}

func TestFakeClient(t *testing.T) {
	_, err := getFakeClient()
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegistryClient(t *testing.T) {
	Flags.Tunnel = ""
	Flags.Endpoint = "http://127.0.0.1:4001"
	_, err := getRegistryClient()
	if err != nil && !strings.Contains(err.Error(), "no such host") {
		t.Fatal(err)
	}
}
