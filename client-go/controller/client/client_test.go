package client

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"testing"
)

const sFile string = `{"username":"t","ssl_verify":false,"controller":"http://d.t","token":"a","response_limit": 50}`

func createTempProfile(contents string) error {
	name, err := ioutil.TempDir("", "client")

	if err != nil {
		return err
	}

	os.Unsetenv("DEIS_PROFILE")
	os.Setenv("HOME", name)
	folder := path.Join(name, "/.deis/")
	if err = os.Mkdir(folder, 0755); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(folder, "client.json"), []byte(contents), 0775); err != nil {
		return err
	}

	return nil
}

func TestLoadSave(t *testing.T) {
	if err := createTempProfile(sFile); err != nil {
		t.Fatal(err)
	}

	client, err := New()

	if err != nil {
		t.Fatal(err)
	}

	expectedB := false
	if client.SSLVerify != expectedB {
		t.Errorf("Expected %t, Got %t", expectedB, client.SSLVerify)
	}

	expected := "a"
	if client.Token != expected {
		t.Errorf("Expected %s, Got %s", expected, client.Token)
	}

	expected = "t"
	if client.Username != expected {
		t.Errorf("Expected %s, Got %s", expected, client.Username)
	}

	expected = "http://d.t"
	if client.ControllerURL.String() != expected {
		t.Errorf("Expected %s, Got %s", expected, client.ControllerURL.String())
	}

	expectedI := 50
	if client.ResponseLimit != expectedI {
		t.Errorf("Expected %d, Got %d", expectedI, client.ResponseLimit)
	}

	client.SSLVerify = true
	client.Token = "b"
	client.Username = "c"
	client.ResponseLimit = 0

	u, err := url.Parse("http://deis.test")

	if err != nil {
		t.Fatal(err)
	}

	client.ControllerURL = *u

	if err = client.Save(); err != nil {
		t.Fatal(err)
	}

	client, err = New()

	expectedB = true
	if client.SSLVerify != expectedB {
		t.Errorf("Expected %t, Got %t", expectedB, client.SSLVerify)
	}

	expected = "b"
	if client.Token != expected {
		t.Errorf("Expected %s, Got %s", expected, client.Token)
	}

	expected = "c"
	if client.Username != expected {
		t.Errorf("Expected %s, Got %s", expected, client.Username)
	}

	expected = "http://deis.test"
	if client.ControllerURL.String() != expected {
		t.Errorf("Expected %s, Got %s", expected, client.ControllerURL.String())
	}

	expectedI = 100
	if client.ResponseLimit != expectedI {
		t.Errorf("Expected %d, Got %d", expectedI, client.ResponseLimit)
	}
}

func TestDeleteSettings(t *testing.T) {
	if err := createTempProfile(""); err != nil {
		t.Fatal(err)
	}

	if err := Delete(); err != nil {
		t.Fatal(err)
	}

	file := locateSettingsFile()

	if _, err := os.Stat(file); err == nil {
		t.Errorf("File %s exists, supposed to have been deleted.", file)
	}
}
