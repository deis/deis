package apps

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

const appFixture string = `
{
    "created": "2014-01-01T00:00:00UTC",
    "id": "example-go",
    "owner": "test",
    "structure": {},
    "updated": "2014-01-01T00:00:00UTC",
    "url": "example-go.example.com",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`

const appsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "created": "2014-01-01T00:00:00UTC",
            "id": "example-go",
            "owner": "test",
            "structure": {},
            "updated": "2014-01-01T00:00:00UTC",
            "url": "example-go.example.com",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
        }
    ]
}`

const appCreateExpected string = `{"id":"example-go"}`
const appRunExpected string = `{"command":"echo hi"}`
const appTransferExpected string = `{"owner":"test"}`

type fakeHTTPServer struct {
	createID        bool
	createWithoutID bool
}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/apps/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == appCreateExpected && f.createID == false {
			f.createID = true
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(appFixture))
			return
		} else if string(body) == "" && f.createWithoutID == false {
			f.createWithoutID = true
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(appFixture))
			return
		}

		fmt.Printf("Unexpected Body: %s'\n", body)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v1/apps/" && req.Method == "GET" {
		res.Write([]byte(appsFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/" && req.Method == "GET" {
		res.Write([]byte(appFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v1/apps/example-go/logs" && req.URL.RawQuery == "" && req.Method == "GET" {
		res.Write([]byte("test\nfoo\nbar\n"))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/logs" && req.URL.RawQuery == "log_lines=1" && req.Method == "GET" {
		res.Write([]byte("test\n"))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/run" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appRunExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appRunExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.Write([]byte(`[0,"hi\n"]`))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appTransferExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appTransferExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestAppsCreate(t *testing.T) {
	t.Parallel()

	expected := api.App{
		ID:      "example-go",
		Created: "2014-01-01T00:00:00UTC",
		Owner:   "test",
		Updated: "2014-01-01T00:00:00UTC",
		URL:     "example-go.example.com",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	handler := fakeHTTPServer{createID: false, createWithoutID: false}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	for _, id := range []string{"example-go", ""} {
		actual, err := New(&client, id)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %v, Got %v", expected, actual)
		}
	}
}

func TestAppsGet(t *testing.T) {
	t.Parallel()

	expected := api.App{
		ID:      "example-go",
		Created: "2014-01-01T00:00:00UTC",
		Owner:   "test",
		Updated: "2014-01-01T00:00:00UTC",
		URL:     "example-go.example.com",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	actual, err := Get(&client, "example-go")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppsDestroy(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	if err = Delete(&client, "example-go"); err != nil {
		t.Fatal(err)
	}
}

func TestAppsRun(t *testing.T) {
	t.Parallel()

	expected := api.AppRunResponse{
		Output:     "hi\n",
		ReturnCode: 0,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	actual, err := Run(&client, "example-go", "echo hi")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppsList(t *testing.T) {
	t.Parallel()

	expected := []api.App{
		api.App{
			ID:      "example-go",
			Created: "2014-01-01T00:00:00UTC",
			Owner:   "test",
			Updated: "2014-01-01T00:00:00UTC",
			URL:     "example-go.example.com",
			UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
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

type testExpected struct {
	Input    int
	Expected string
}

func TestAppsLogs(t *testing.T) {
	t.Parallel()

	tests := []testExpected{
		testExpected{
			Input:    -1,
			Expected: "test\nfoo\nbar\n",
		},
		testExpected{
			Input:    1,
			Expected: "test\n",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	for _, test := range tests {
		actual, err := Logs(&client, "example-go", test.Input)

		if err != nil {
			t.Error(err)
		}

		if actual != test.Expected {
			t.Errorf("Expected %s, Got %s", test.Expected, actual)
		}
	}
}

func TestAppsTransfer(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	if err = Transfer(&client, "example-go", "test"); err != nil {
		t.Fatal(err)
	}
}
