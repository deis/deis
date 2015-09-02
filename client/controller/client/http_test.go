package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/deis/deis/version"
)

type fakeHTTPServer struct{}

const limitedFixture string = `
{
    "count": 4,
    "next": "http://replaced.com/limited2/",
    "previous": null,
    "results": [
        {
            "test": "foo"
        },
        {
            "test": "bar"
        }
    ]
}
`

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

	if req.URL.Path == "/limited/" && req.Method == "GET" && req.URL.RawQuery == "page_size=2" {
		res.Write([]byte(limitedFixture))
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

	if err = CheckConnection(httpClient, *u); err != nil {
		t.Error(err)
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

	body, err := client.BasicRequest("POST", "/basic/", []byte("test"))

	if err != nil {
		t.Fatal(err)
	}

	expected := "basic"
	if body != expected {
		t.Errorf("Expected %s, Got %s", expected, body)
	}
}

func TestCheckErrors(t *testing.T) {
	t.Parallel()

	expected := `
404 NOT FOUND
error: This is an error.
error_array: This is an array.
error_array: Foo!
`
	altExpected := `
404 NOT FOUND
error_array: This is an array.
error_array: Foo!
error: This is an error.
`

	body := `
{
	"error": "This is an error.",
	"error_array": [
		"This is an array.",
		"Foo!"
	]
}`

	res := http.Response{
		StatusCode: 404,
		Status:     "404 NOT FOUND",
	}

	actual := checkForErrors(&res, body).Error()

	if actual != expected && actual != altExpected {
		t.Errorf("Expected %s or %s, Got %s", expected, altExpected, actual)
	}

	expected = `
503 Service Temporarily Unavailable
<html>
<head><title>503 Service Temporarily Unavailable</title></head>
<body bgcolor="white">
<center><h1>503 Service Temporarily Unavailable</h1></center>
<hr><center>nginx/1.9.4</center>
</body>
</html>
`

	body = `<html>
<head><title>503 Service Temporarily Unavailable</title></head>
<body bgcolor="white">
<center><h1>503 Service Temporarily Unavailable</h1></center>
<hr><center>nginx/1.9.4</center>
</body>
</html>`

	res = http.Response{
		StatusCode: http.StatusServiceUnavailable,
		Status:     "503 Service Temporarily Unavailable",
	}

	actual = checkForErrors(&res, body).Error()

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestCheckErrorsReturnsNil(t *testing.T) {
	t.Parallel()

	responses := []http.Response{
		http.Response{
			StatusCode: http.StatusOK,
		},
		http.Response{
			StatusCode: http.StatusCreated,
		},
		http.Response{
			StatusCode: http.StatusNoContent,
		},
	}

	for _, res := range responses {
		if err := checkForErrors(&res, ""); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLimitedRequest(t *testing.T) {
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

	expected := `[{"test":"foo"},{"test":"bar"}]`
	expectedC := 4

	actual, count, err := client.LimitedRequest("/limited/", 2)

	if err != nil {
		t.Fatal(err)
	}

	if count != expectedC {
		t.Errorf("Expected %d, Got %d", expectedC, count)
	}

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}
