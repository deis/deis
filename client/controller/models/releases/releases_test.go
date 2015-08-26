package releases

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

const releasesFixture string = `
{
    "count": 3,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "build": null,
            "config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
            "created": "2014-01-01T00:00:00UTC",
            "owner": "test",
            "summary": "test created initial release",
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "version": 1
        }
    ]
}`

const releaseFixture string = `
{
    "app": "example-go",
    "build": null,
    "config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
    "created": "2014-01-01T00:00:00UTC",
    "owner": "test",
    "summary": "test created initial release",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
    "version": 1
}
`

const rollbackFixture string = `
{"version": 5}
`
const rollbackerFixture string = `
{"version": 7}
`

const rollbackExpected string = `{"version":2}`
const rollbackerExpected string = ``

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/apps/example-go/releases/" && req.Method == "GET" {
		res.Write([]byte(releasesFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/releases/v1/" && req.Method == "GET" {
		res.Write([]byte(releaseFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/releases/rollback/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != rollbackExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", rollbackExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(rollbackFixture))
		return
	}

	if req.URL.Path == "/v1/apps/rollbacker/releases/rollback/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != rollbackerExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", rollbackerExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(rollbackerFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestReleasesList(t *testing.T) {
	t.Parallel()

	expected := []api.Release{
		api.Release{
			App:     "example-go",
			Build:   "",
			Config:  "95bd6dea-1685-4f78-a03d-fd7270b058d1",
			Created: "2014-01-01T00:00:00UTC",
			Owner:   "test",
			Summary: "test created initial release",
			Updated: "2014-01-01T00:00:00UTC",
			UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			Version: 1,
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

	actual, _, err := List(&client, "example-go", 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestReleasesGet(t *testing.T) {
	t.Parallel()

	expected := api.Release{
		App:     "example-go",
		Build:   "",
		Config:  "95bd6dea-1685-4f78-a03d-fd7270b058d1",
		Created: "2014-01-01T00:00:00UTC",
		Owner:   "test",
		Summary: "test created initial release",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Version: 1,
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

	actual, err := Get(&client, "example-go", 1)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestRollback(t *testing.T) {
	t.Parallel()

	expected := 5

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	actual, err := Rollback(&client, "example-go", 2)

	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestRollbacker(t *testing.T) {
	t.Parallel()

	expected := 7

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	actual, err := Rollback(&client, "rollbacker", -1)

	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}
