package config

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

const configFixture string = `
{
    "owner": "test",
    "app": "example-go",
    "values": {
      "TEST": "testing",
      "FOO": "bar"
    },
    "memory": {
      "web": "1G"
    },
    "cpu": {
      "web": "1000"
    },
    "tags": {
      "test": "tests"
    },
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const configUnsetFixture string = `
{
    "owner": "test",
    "app": "unset-test",
    "values": {},
    "memory": {},
    "cpu": {},
    "tags": {},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const configSetExpected string = `{"values":{"FOO":"bar","TEST":"testing"},"memory":{"web":"1G"},"cpu":{"web":"1000"},"tags":{"test":"tests"}}`
const configUnsetExpected string = `{"values":{"FOO":null,"TEST":null},"memory":{"web":null},"cpu":{"web":null},"tags":{"test":null}}`

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/apps/example-go/config/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != configSetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", configSetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(configFixture))
		return
	}

	if req.URL.Path == "/v1/apps/unset-test/config/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != configUnsetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", configUnsetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(configUnsetFixture))
		return
	}

	if req.URL.Path == "/v1/apps/example-go/config/" && req.Method == "GET" {
		res.Write([]byte(configFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestConfigSet(t *testing.T) {
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

	expected := api.Config{
		Owner: "test",
		App:   "example-go",
		Values: map[string]interface{}{
			"TEST": "testing",
			"FOO":  "bar",
		},
		Memory: map[string]interface{}{
			"web": "1G",
		},
		CPU: map[string]interface{}{
			"web": "1000",
		},
		Tags: map[string]interface{}{
			"test": "tests",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	configVars := api.Config{
		Values: map[string]interface{}{
			"TEST": "testing",
			"FOO":  "bar",
		},
		Memory: map[string]interface{}{
			"web": "1G",
		},
		CPU: map[string]interface{}{
			"web": "1000",
		},
		Tags: map[string]interface{}{
			"test": "tests",
		},
	}

	actual, err := Set(&client, "example-go", configVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestConfigUnset(t *testing.T) {
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

	expected := api.Config{
		Owner:   "test",
		App:     "unset-test",
		Values:  map[string]interface{}{},
		Memory:  map[string]interface{}{},
		CPU:     map[string]interface{}{},
		Tags:    map[string]interface{}{},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	configVars := api.Config{
		Values: map[string]interface{}{
			"TEST": nil,
			"FOO":  nil,
		},
		Memory: map[string]interface{}{
			"web": nil,
		},
		CPU: map[string]interface{}{
			"web": nil,
		},
		Tags: map[string]interface{}{
			"test": nil,
		},
	}

	actual, err := Set(&client, "unset-test", configVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestConfigList(t *testing.T) {
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

	expected := api.Config{
		Owner: "test",
		App:   "example-go",
		Values: map[string]interface{}{
			"TEST": "testing",
			"FOO":  "bar",
		},
		Memory: map[string]interface{}{
			"web": "1G",
		},
		CPU: map[string]interface{}{
			"web": "1000",
		},
		Tags: map[string]interface{}{
			"test": "tests",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	actual, err := List(&client, "example-go")

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
