package certs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/version"
)

const certsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "common_name": "test.example.com",
            "expires": "2014-01-01T00:00:00UTC"
        }
    ]
}`

const certFixture string = `
{
    "updated": "2014-01-01T00:00:00UTC",
    "created": "2014-01-01T00:00:00UTC",
    "expires": "2015-01-01T00:00:00UTC",
    "common_name": "test.example.com",
    "owner": "test",
    "id": 1
}`

const certExpected string = `{"certificate":"test","key":"foo","common_name":"test.example.com"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/certs/" && req.Method == "GET" {
		res.Write([]byte(certsFixture))
		return
	}

	if req.URL.Path == "/v1/certs/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != certExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", certExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(certFixture))
		return
	}

	if req.URL.Path == "/v1/certs/test.example.com" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestCertsList(t *testing.T) {
	t.Parallel()

	expected := []api.Cert{
		api.Cert{
			Name:    "test.example.com",
			Expires: "2014-01-01T00:00:00UTC",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	actual, _, err := List(&client, 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestCert(t *testing.T) {
	t.Parallel()

	expected := api.Cert{
		Updated: "2014-01-01T00:00:00UTC",
		Created: "2014-01-01T00:00:00UTC",
		Expires: "2015-01-01T00:00:00UTC",
		Name:    "test.example.com",
		Owner:   "test",
		ID:      1,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	actual, err := New(&client, "test", "foo", "test.example.com")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestCertDeleteion(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	if err = Delete(&client, "test.example.com"); err != nil {
		t.Fatal(err)
	}
}
