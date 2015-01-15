package config

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"
)

// TestConfigSSHPrivateKey ensures private keys are base64 encoded from file path
func TestConfigSSHPrivateKey(t *testing.T) {

	f, err := writeTempFile("private-key")
	if err != nil {
		t.Fatal(err)
	}

	val, err := valueForPath("/deis/platform/sshPrivateKey", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte("private-key"))

	if val != encoded {
		t.Fatalf("expected: %v, got: %v", encoded, val)
	}
}

func TestConfigRouterKey(t *testing.T) {

	f, err := writeTempFile("router-key")
	if err != nil {
		t.Fatal(err)
	}

	val, err := valueForPath("/deis/router/sslKey", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if val != "router-key" {
		t.Fatalf("expected: router-key, got: %v", val)
	}

}

func TestConfigRouterCert(t *testing.T) {

	f, err := writeTempFile("router-cert")
	if err != nil {
		t.Fatal(err)
	}

	val, err := valueForPath("/deis/router/sslCert", f.Name())
	if err != nil {
		t.Fatal(err)
	}

	if val != "router-cert" {
		t.Fatalf("expected: router-cert, got: %v", val)
	}

}

func writeTempFile(data string) (*os.File, error) {
	f, err := ioutil.TempFile("", "deisctl")
	if err != nil {
		return nil, err
	}

	f.Write([]byte(data))
	defer f.Close()

	return f, nil
}
