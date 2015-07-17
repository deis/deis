package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/deis/deis/version"
)

const rE string = `{"username":"test","password":"opensesame","email":"test@example.com"}`
const lE string = `{"username":"test","password":"opensesame"}`
const pE string = `{"username":"test","password":"old","new_password":"new"}`
const rAE string = `{"all":true}`
const rUE string = `{"username":"test"}`

type fakeAuthHTTPServer struct {
	regenBodyEmpty    bool
	regenBodyAll      bool
	regenBodyUsername bool
}

func (f fakeAuthHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/auth/register/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != rE {
			fmt.Printf("Expected '%s', Got '%s'\n", rE, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v1/auth/login/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != lE {
			fmt.Printf("Expected '%s', Got '%s'\n", lE, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.Write([]byte(`{"token":"abc"}`))
		return
	}

	if req.URL.Path == "/v1/auth/passwd/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != pE {
			fmt.Printf("Expected '%s', Got '%s'\n", lE, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.Write(nil)
		return
	}

	if req.URL.Path == "/v1/auth/tokens/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == rAE && !f.regenBodyAll {
			f.regenBodyAll = true
			res.Write(nil)
			return
		} else if string(body) == rUE && !f.regenBodyUsername {
			f.regenBodyUsername = true
			res.Write(nil)
			return
		} else if !f.regenBodyEmpty {
			f.regenBodyEmpty = true
			res.Write([]byte(`{"token":"abc"}`))
			return
		}

		fmt.Printf("Expected '%s', Got '%s'\n", lE, body)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v1/auth/cancel/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestRegister(t *testing.T) {
	t.Parallel()

	handler := fakeAuthHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	if err = Register(*u, "test", "opensesame", "test@example.com", false, false); err != nil {
		t.Error(err)
	}
}

func TestLogin(t *testing.T) {
	err := createTempProfile("")

	handler := fakeAuthHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	controllerURL, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	if Login(*controllerURL, "test", "opensesame", false); err != nil {
		t.Error(err)
	}

	client, err := New()

	if err != nil {
		t.Fatal(err)
	}

	if client.ControllerURL.String() != controllerURL.String() {
		t.Errorf("Expected %s, Got %s", controllerURL.String(), client.ControllerURL.String())
	}

	expected := "test"
	if client.Username != expected {
		t.Errorf("Expected %s, Got %s", expected, client.Username)
	}

	expected = "abc"
	if client.Token != expected {
		t.Errorf("Expected %s, Got %s", expected, client.Token)
	}

	expectedB := false
	if client.SSLVerify != expectedB {
		t.Errorf("Expected %t, Got %t", expectedB, client.SSLVerify)
	}
}

func TestLogout(t *testing.T) {
	err := createTempProfile(sFile)

	if err != nil {
		t.Fatal(err)
	}

	if err = Logout(); err != nil {
		t.Fatal(err)
	}

	file := locateSettingsFile()

	if _, err := os.Stat(file); err == nil {
		t.Errorf("File %s exists, supposed to have been deleted.", file)
	}
}

func TestPasswd(t *testing.T) {
	handler := fakeAuthHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	sF := fmt.Sprintf(`{"username":"t","ssl_verify":false,"controller":"%s","token":"a"}`, server.URL)
	err := createTempProfile(sF)

	if err != nil {
		t.Fatal(err)
	}

	if err = Passwd("test", "old", "new"); err != nil {
		t.Error(err)
	}
}

func TestCancel(t *testing.T) {
	handler := fakeAuthHTTPServer{regenBodyEmpty: false, regenBodyAll: false,
		regenBodyUsername: false}
	server := httptest.NewServer(handler)
	defer server.Close()

	sF := fmt.Sprintf(`{"username":"t","ssl_verify":false,"controller":"%s","token":"a"}`, server.URL)
	err := createTempProfile(sF)

	if err != nil {
		t.Fatal(err)
	}

	if err = Regenerate("", true); err != nil {
		t.Error(err)
	}

	if err = Cancel(); err != nil {
		t.Error(err)
	}

	file := locateSettingsFile()

	if _, err := os.Stat(file); err == nil {
		t.Errorf("File %s exists, supposed to have been deleted.", file)
	}
}

func TestRegenerate(t *testing.T) {
	handler := fakeAuthHTTPServer{regenBodyEmpty: false, regenBodyAll: false,
		regenBodyUsername: false}
	server := httptest.NewServer(handler)
	defer server.Close()

	sF := fmt.Sprintf(`{"username":"t","ssl_verify":false,"controller":"%s","token":"a"}`, server.URL)
	err := createTempProfile(sF)

	if err != nil {
		t.Fatal(err)
	}

	if err = Regenerate("", true); err != nil {
		t.Error(err)
	}

	if err = Regenerate("test", false); err != nil {
		t.Error(err)
	}

	if err = Regenerate("", false); err != nil {
		t.Error(err)
	}

	client, err := New()

	if err != nil {
		t.Error(err)
	}

	expected := "abc"
	if client.Token != expected {
		t.Errorf("Expected %s, Got %s", expected, client.Token)
	}
}
