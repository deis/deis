package etcd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Spot-check NewClient can identify good and bad endpoints
func TestNewClient(t *testing.T) {
	tests := []struct {
		endpoints []string
		pass      bool
	}{
		// these should result in the default endpoint being used
		{[]string{}, true},
		{nil, true},

		// simplest good endpoint, just a scheme and IP
		{[]string{"http://192.0.2.3"}, true},

		// multiple valid values
		{[]string{"http://192.0.2.3", "http://192.0.2.4"}, true},

		// completely invalid URL
		{[]string{"://"}, false},

		// bogus endpoint filtered by our own logic
		{[]string{"boots://pants"}, false},

		// good endpoint followed by a bogus endpoint
		{[]string{"http://192.0.2.3", "boots://pants"}, false},
	}

	for i, tt := range tests {
		_, err := NewClient(tt.endpoints, http.Transport{})
		if tt.pass != (err == nil) {
			t.Errorf("case %d %v: expected to pass=%t, err=%v", i, tt.endpoints, tt.pass, err)
		}
	}
}

// client.SetDefaultPath should only overwrite the path if it is unset
func TestSetDefaultPath(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"http://example.com", "http://example.com/"},
		{"http://example.com/", "http://example.com/"},
		{"http://example.com/foo", "http://example.com/foo"},
	}

	for i, tt := range tests {
		u, _ := url.Parse(tt.in)
		if tt.in != u.String() {
			t.Errorf("case %d: url.Parse modified the URL before we could test it")
			continue
		}

		setDefaultPath(u)
		if tt.out != u.String() {
			t.Errorf("case %d: expected output of %s did not match actual value %s", i, tt.out, u.String())
		}
	}
}

// Enumerate the many permutations of an endpoint, asserting whether or
// not they should be acceptable
func TestFilterURL(t *testing.T) {
	tests := []struct {
		endpoint string
		pass     bool
	}{
		// IP & port
		{"http://192.0.2.3:4001/", true},

		// trailing slash
		{"http://192.0.2.3/", true},

		// hostname
		{"http://example.com/", true},

		// no host info
		{"http:///foo/bar", false},

		// empty path
		{"http://192.0.2.3", false},

		// custom path
		{"http://192.0.2.3/foo/bar", false},

		// custom query params
		{"http://192.0.2.3/?foo=bar", false},

		// no scheme
		{"192.0.2.3:4002/", false},

		// non-http scheme
		{"boots://192.0.2.3:4002/", false},

		// https scheme fails for now
		{"https://192.0.2.3:4002/", false},

		// no slash after scheme (url.URL.Opaque)
		{"http:192.0.2.3/", false},

		// user info
		{"http://elroy@192.0.2.3/", false},

		// fragment
		{"http://192.0.2.3/#foo", false},
	}

	for i, tt := range tests {
		u, _ := url.Parse(tt.endpoint)
		if tt.endpoint != u.String() {
			t.Errorf("case %d: url.Parse modified the URL before we could test it")
			continue
		}

		err := filterURL(u)

		if tt.pass != (err == nil) {
			t.Errorf("case %d %v: expected to pass=%t, err=%v", i, tt.endpoint, tt.pass, err)
		}
	}
}

// Ensure the channel passed into c.resolve is actually wired up
func TestClientCancel(t *testing.T) {
	act := Get{Key: "/foo"}
	c, err := NewClient(nil, http.Transport{})
	if err != nil {
		t.Fatalf("Failed building Client: %v", err)
	}

	cancel := make(chan bool)
	sentinel := make(chan bool, 2)

	rf := func(req *http.Request, cancel <-chan bool) (*http.Response, []byte, error) {
		<-cancel
		sentinel <- true
		return nil, nil, errors.New("Cancelled")
	}

	go func() {
		c.resolve(&act, rf, cancel)
		sentinel <- true
	}()

	select {
	case <-sentinel:
		t.Fatalf("sentinel should not be ready")
	default:
	}

	close(cancel)

	for i := 0; i < 2; i++ {
		select {
		case <-sentinel:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("timed out waiting for sentinel value")
		}
	}
}

type clientStep struct {
	method string
	url    string

	resp http.Response
}

func assertClientSteps(t *testing.T, c *client, act Action, steps []clientStep, expectSuccess bool) {
	idx := 0
	rf := func(req *http.Request, cancel <-chan bool) (*http.Response, []byte, error) {
		if idx >= len(steps) {
			t.Fatalf("Received too many requests")
		}
		step := steps[idx]
		idx = idx + 1

		if step.method != req.Method {
			t.Fatalf("step %d: request method is %s, expected %s", idx, req.Method, step.method)
		}

		if step.url != req.URL.String() {
			t.Fatalf("step %d: request URL is %s, expected %s", idx, req.URL, step.url)
		}

		var body []byte
		if step.resp.Body != nil {
			var err error
			body, err = ioutil.ReadAll(step.resp.Body)
			if err != nil {
				t.Fatalf("step %d: failed preparing body: %v", idx, err)
			}
		}

		return &step.resp, body, nil
	}

	_, err := c.resolve(act, rf, make(chan bool))
	if expectSuccess != (err == nil) {
		t.Fatalf("expected to pass=%t, err=%v", expectSuccess, err)
	}
}

// Follow all redirects, using the full Location header regardless of how crazy it seems
func TestClientRedirectsFollowed(t *testing.T) {
	steps := []clientStep{
		{
			"GET", "http://192.0.2.1:4001/v2/keys/foo?recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://192.0.2.2:4001/v2/keys/foo?recursive=false&sorted=false"},
				},
			},
		},
		{
			"GET", "http://192.0.2.2:4001/v2/keys/foo?recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://192.0.2.3:4002/pants?recursive=true"},
				},
			},
		},
		{
			"GET", "http://192.0.2.3:4002/pants?recursive=true",
			http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"X-Etcd-Index": {"123"}},
				Body:       ioutil.NopCloser(strings.NewReader("{}")),
			},
		},
	}

	c, err := NewClient([]string{"http://192.0.2.1:4001"}, http.Transport{})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	act := &Get{Key: "/foo"}
	assertClientSteps(t, c, act, steps, true)
}

// Follow a redirect to a failing node, then fall back to the healthy second endpoint
func TestClientRedirectsAndAlternateEndpoints(t *testing.T) {
	steps := []clientStep{
		{
			"GET", "http://192.0.2.1:4001/v2/keys/foo?recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://192.0.2.5:4001/v2/keys/foo?recursive=true"},
				},
			},
		},
		{
			"GET", "http://192.0.2.5:4001/v2/keys/foo?recursive=true",
			http.Response{
				StatusCode: http.StatusGatewayTimeout,
			},
		},
		{
			"GET", "http://192.0.2.2:4002/v2/keys/foo?recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"X-Etcd-Index": {"123"}},
				Body:       ioutil.NopCloser(strings.NewReader("{}")),
			},
		},
	}

	c, err := NewClient([]string{"http://192.0.2.1:4001", "http://192.0.2.2:4002"}, http.Transport{})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	act := &Get{Key: "/foo"}
	assertClientSteps(t, c, act, steps, true)
}

func TestClientRedirectOverLimit(t *testing.T) {
	reqCount := 0
	rf := func(req *http.Request, cancel <-chan bool) (*http.Response, []byte, error) {
		reqCount = reqCount + 1

		if reqCount > 10 {
			t.Fatalf("c.resolve made %d requests, expected max of 10", reqCount)
		}

		resp := http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Header: http.Header{
				"Location": {"http://127.0.0.1:4001/"},
			},
		}

		return &resp, []byte{}, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan bool))
	if req != nil || err != nil {
		t.Errorf("Expected nil resp and nil err, got resp=%v and err=%v", req, err)
	}

	if reqCount != 10 {
		t.Fatalf("c.resolve should have made 10 responses, got %d", reqCount)
	}
}

func TestClientRedirectMax(t *testing.T) {
	count := 0
	rf := func(req *http.Request, cancel <-chan bool) (*http.Response, []byte, error) {
		var resp http.Response
		var body []byte

		count = count + 1

		if count == 10 {
			resp = http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"X-Etcd-Index": {"123"},
				},
			}
			body = []byte("{}")
		} else {
			resp = http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://127.0.0.1:4001/"},
				},
			}
		}

		return &resp, body, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan bool))
	if req == nil || err != nil {
		t.Errorf("Expected non-nil resp and nil err, got resp=%v and err=%v", req, err)
	}
}

func TestClientRequestFuncError(t *testing.T) {
	rf := func(req *http.Request, cancel <-chan bool) (*http.Response, []byte, error) {
		return nil, nil, errors.New("bogus error")
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan bool))
	if req != nil {
		t.Errorf("Expected req=nil, got %v", nil)
	}
	if err != nil {
		t.Errorf("Expected err=nil, got %v", err)
	}
}

func TestClientRedirectNowhere(t *testing.T) {
	rf := func(req *http.Request, cancel <-chan bool) (*http.Response, []byte, error) {
		resp := http.Response{StatusCode: http.StatusTemporaryRedirect}
		return &resp, []byte{}, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:4001")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan bool))
	if req != nil {
		t.Errorf("Expected req=nil, got %v", nil)
	}
	if err != nil {
		t.Errorf("Expected err=nil, got %v", err)
	}
}
