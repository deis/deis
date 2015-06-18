package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/deis/deis/version"
)

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	eA := "Deis Client v" + version.Version

	if req.Header.Get("User-Agent") != eA {
		fmt.Printf("User Agent Wrong: Expected %s, Got %s\n", eA, req.Header.Get("User-Agent"))
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v1/" {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/raw/" && req.Method == "POST" {
		res.Write([]byte("test"))
		return
	}

	if req.URL.Path == "/basic/" && req.Method == "POST" {
		eT := "token abc"
		if req.Header.Get("Authorization") != eT {
			fmt.Printf("Token Wrong: Expected %s, Got %s\n", eT, req.Header.Get("Authorization"))
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		eC := "application/json"
		if req.Header.Get("Content-Type") != eC {
			fmt.Printf("Content Type Wrong: Expected %s, Got %s\n", eC, req.Header.Get("Content-Type"))
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		eB := "test"
		if string(body) != eB {
			fmt.Printf("Body Wrong: Expected %s, Got %s\n", eB, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.Write([]byte("basic"))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestCheckConnection(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := CreateHTTPClient(false)

	if err = CheckConection(httpClient, *u); err != nil {
		t.Error(err)
	}
}

func TestRawRequest(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)
	u.Path = "/raw/"

	if err != nil {
		t.Fatal(err)
	}

	client := CreateHTTPClient(false)

	headers := http.Header{}
	addUserAgent(&headers)

	res, err := rawRequest(client, "POST", u.String(), bytes.NewBuffer(nil), headers, 200)

	if err != nil {
		t.Fatal(err)
	}

	actual, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	res.Body.Close()

	expected := "test"
	if string(actual) != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestBasicRequest(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := CreateHTTPClient(false)

	client := Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	body, status, err := client.BasicRequest("POST", "/basic/", []byte("test"))

	if err != nil {
		t.Fatal(err)
	}

	expectedStatus := 200

	if status != expectedStatus {
		t.Errorf("Expected %d, got %d", expectedStatus, status)
	}

	expected := "basic"
	if body != expected {
		t.Errorf("Expected %s, Got %s", expected, body)
	}
}
