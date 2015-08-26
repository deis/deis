package ps

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

const processesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "owner": "test",
            "app": "example-go",
            "release": "v2",
            "created": "2014-01-01T00:00:00UTC",
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "type": "web",
            "num": 1,
            "state": "up"
        }
    ]
}`

const restartAllFixture string = `[
    {
        "owner": "test",
        "app": "example-go",
        "release": "v2",
        "created": "2014-01-01T00:00:00UTC",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
        "type": "web",
        "num": 1,
        "state": "up"
    }
]
`

const restartWorkerFixture string = `[
    {
        "owner": "test",
        "app": "example-go",
        "release": "v2",
        "created": "2014-01-01T00:00:00UTC",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
        "type": "worker",
        "num": 1,
        "state": "up"
    }
]
`

const restartWebTwoFixture string = `[
    {
        "owner": "test",
        "app": "example-go",
        "release": "v2",
        "created": "2014-01-01T00:00:00UTC",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
        "type": "web",
        "num": 2,
        "state": "up"
    }
]
`

const scaleExpected string = `{"web":2}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/apps/example-go/containers/" && req.Method == "GET" {
		res.Write([]byte(processesFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/containers/restart/" && req.Method == "POST" {
		res.Write([]byte(restartAllFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/containers/worker/restart/" && req.Method == "POST" {
		res.Write([]byte(restartWorkerFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/containers/web/2/restart/" && req.Method == "POST" {
		res.Write([]byte(restartWebTwoFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/scale/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != scaleExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", scaleExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestProcessesList(t *testing.T) {
	t.Parallel()

	expected := []api.Process{
		api.Process{
			Owner:   "test",
			App:     "example-go",
			Release: "v2",
			Created: "2014-01-01T00:00:00UTC",
			Updated: "2014-01-01T00:00:00UTC",
			UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			Type:    "web",
			Num:     1,
			State:   "up",
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

type testExpected struct {
	Num      int
	Type     string
	Expected []api.Process
}

func TestAppsRestart(t *testing.T) {
	t.Parallel()

	tests := []testExpected{
		testExpected{
			Num:  -1,
			Type: "",
			Expected: []api.Process{
				api.Process{
					Owner:   "test",
					App:     "example-go",
					Release: "v2",
					Created: "2014-01-01T00:00:00UTC",
					Updated: "2014-01-01T00:00:00UTC",
					UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
					Type:    "web",
					Num:     1,
					State:   "up",
				},
			},
		},
		testExpected{
			Num:  -1,
			Type: "worker",
			Expected: []api.Process{
				api.Process{
					Owner:   "test",
					App:     "example-go",
					Release: "v2",
					Created: "2014-01-01T00:00:00UTC",
					Updated: "2014-01-01T00:00:00UTC",
					UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
					Type:    "worker",
					Num:     1,
					State:   "up",
				},
			},
		},
		testExpected{
			Num:  2,
			Type: "web",
			Expected: []api.Process{
				api.Process{
					Owner:   "test",
					App:     "example-go",
					Release: "v2",
					Created: "2014-01-01T00:00:00UTC",
					Updated: "2014-01-01T00:00:00UTC",
					UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
					Type:    "web",
					Num:     2,
					State:   "up",
				},
			},
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
		actual, err := Restart(&client, "example-go", test.Type, test.Num)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(test.Expected, actual) {
			t.Error(fmt.Errorf("Expected %v, Got %v", test.Expected, actual))
		}
	}
}

func TestScale(t *testing.T) {
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

	if err = Scale(&client, "example-go", map[string]int{"web": 2}); err != nil {
		t.Fatal(err)
	}
}
