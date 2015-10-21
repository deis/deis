// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
		_, err := NewClient(tt.endpoints, &http.Transport{}, time.Second)
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
			t.Errorf("case %d: url.Parse modified the URL before we could test it", i)
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
		{"http://192.0.2.3:2379/", true},

		// trailing slash
		{"http://192.0.2.3/", true},

		// hostname
		{"http://example.com/", true},

		// https scheme
		{"https://192.0.2.3:4002/", true},

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
			t.Errorf("case %d: url.Parse modified the URL before we could test it", i)
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
	c, err := NewClient(nil, &http.Transport{}, time.Second)
	if err != nil {
		t.Fatalf("Failed building Client: %v", err)
	}

	cancel := make(chan struct{})
	sentinel := make(chan struct{}, 2)

	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		<-cancel
		sentinel <- struct{}{}
		return nil, nil, errors.New("Cancelled")
	}

	go func() {
		c.resolve(&act, rf, cancel)
		sentinel <- struct{}{}
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
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
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

	_, err := c.resolve(act, rf, make(chan struct{}))
	if expectSuccess != (err == nil) {
		t.Fatalf("expected to pass=%t, err=%v", expectSuccess, err)
	}
}

// Follow all redirects, using the full Location header regardless of how crazy it seems
func TestClientRedirectsFollowed(t *testing.T) {
	steps := []clientStep{
		{
			"GET", "http://192.0.2.1:2379/v2/keys/foo?consistent=true&recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusTemporaryRedirect,
				Header: http.Header{
					"Location": {"http://192.0.2.2:2379/v2/keys/foo?recursive=false&sorted=false"},
				},
			},
		},
		{
			"GET", "http://192.0.2.2:2379/v2/keys/foo?recursive=false&sorted=false",
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

	c, err := NewClient([]string{"http://192.0.2.1:2379"}, &http.Transport{}, time.Second)
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
			"GET", "http://192.0.2.1:4001/v2/keys/foo?consistent=true&recursive=false&sorted=false",
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
			"GET", "http://192.0.2.2:2379/v2/keys/foo?consistent=true&recursive=false&sorted=false",
			http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"X-Etcd-Index": {"123"}},
				Body:       ioutil.NopCloser(strings.NewReader("{}")),
			},
		},
	}

	c, err := NewClient([]string{"http://192.0.2.1:4001", "http://192.0.2.2:2379"}, &http.Transport{}, time.Second)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	act := &Get{Key: "/foo"}
	assertClientSteps(t, c, act, steps, true)
}

func TestClientRedirectOverLimit(t *testing.T) {
	reqCount := 0
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		reqCount = reqCount + 1

		if reqCount > 10 {
			t.Fatalf("c.resolve made %d requests, expected max of 10", reqCount)
		}

		resp := http.Response{
			StatusCode: http.StatusTemporaryRedirect,
			Header: http.Header{
				"Location": {"http://127.0.0.1:2379/"},
			},
		}

		return &resp, []byte{}, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:2379")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req != nil || err != nil {
		t.Errorf("Expected nil resp and nil err, got resp=%v and err=%v", req, err)
	}

	if reqCount != 10 {
		t.Fatalf("c.resolve should have made 10 responses, got %d", reqCount)
	}
}

func TestClientRedirectMax(t *testing.T) {
	count := 0
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
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
					"Location": {"http://127.0.0.1:2379/"},
				},
			}
		}

		return &resp, body, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:2379")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req == nil || err != nil {
		t.Errorf("Expected non-nil resp and nil err, got resp=%v and err=%v", req, err)
	}
}

func TestClientRequestFuncError(t *testing.T) {
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		return nil, nil, errors.New("bogus error")
	}

	endpoint, err := url.Parse("http://192.0.2.1:2379")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req != nil {
		t.Errorf("Expected req=nil, got %v", nil)
	}
	if err != nil {
		t.Errorf("Expected err=nil, got %v", err)
	}
}

func TestClientRedirectNowhere(t *testing.T) {
	rf := func(req *http.Request, cancel <-chan struct{}) (*http.Response, []byte, error) {
		resp := http.Response{StatusCode: http.StatusTemporaryRedirect}
		return &resp, []byte{}, nil
	}

	endpoint, err := url.Parse("http://192.0.2.1:2379")
	if err != nil {
		t.Fatal(err)
	}

	act := &Get{Key: "/foo"}
	ar := newActionResolver(act, endpoint, rf)

	req, err := ar.Resolve(make(chan struct{}))
	if req != nil {
		t.Errorf("Expected req=nil, got %v", nil)
	}
	if err != nil {
		t.Errorf("Expected err=nil, got %v", err)
	}
}

func newTestingRequestAndClient(t *testing.T, handler http.Handler) (*client, *http.Request) {
	ts := httptest.NewServer(handler)
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	c, err := NewClient(nil, &http.Transport{}, time.Second)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	return c, req
}

func TestGoodRequestHTTP(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "testing")
	})
	c, req := newTestingRequestAndClient(t, h)

	cancel := make(chan struct{})
	resp, body, err := c.requestHTTP(req, cancel)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Errorf("unexpected nil response")
	} else {
		// ensure the body was closed
		var b []byte
		if n, err := resp.Body.Read(b); n != 0 || err == nil {
			t.Errorf("resp.Body.Read() returned unexpectedly: want (0, err), got (%d, %v)", n, err)
		}
	}
	if string(body) != "testing" {
		t.Errorf("unexpected body: got %q, want %q", body, "testing")
	}
}

// transport that returns a nil Response and nil error
type nilNilTransport struct{}

func (n *nilNilTransport) RoundTrip(req *http.Request) (*http.Response,
	error) {
	return nil, nil
}
func (n *nilNilTransport) CancelRequest(req *http.Request) {}

// Ensure that any request that somehow returns (nil, nil) propagates an actual error
func TestNilNilRequestHTTP(t *testing.T) {
	c := &client{[]url.URL{}, &nilNilTransport{}, time.Second}
	cancel := make(chan struct{})
	resp, body, err := c.requestHTTP(nil, cancel)
	if err == nil {
		t.Error("unexpected nil error")
	} else if err.Error() != "nil error and nil response" {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %v", resp)
	}
	if body != nil {
		t.Errorf("unexpected non-nil body: %q", body)
	}
}

// Simple implementation of ReadCloser to serve as response.Body
type rc struct{}

func (r *rc) Read(p []byte) (n int, err error) { return 0, nil }
func (r *rc) Close() error                     { return nil }

// transport that returns a non-nil Response and non-nil error
type respAndErrTransport struct{}

func (r *respAndErrTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Body: &rc{},
	}, errors.New("some error")
}
func (r *respAndErrTransport) CancelRequest(req *http.Request) {}

// Ensure that the body of a response is closed even when an error is returned
func TestRespAndErrRequestHTTP(t *testing.T) {
	c := &client{[]url.URL{}, &respAndErrTransport{}, time.Second}
	cancel := make(chan struct{})
	resp, body, err := c.requestHTTP(nil, cancel)
	if err == nil {
		t.Error("unexpected nil error")
	} else if err.Error() == "cancelled" {
		t.Error("unexpected error, should not be cancelled")
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %v", resp)
	}
	if body != nil {
		t.Errorf("unexpected non-nil body: %q", body)
	}
}

func TestCancelledRequestHTTP(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Hour)
	})
	c, req := newTestingRequestAndClient(t, h)

	cancel := make(chan struct{})
	close(cancel)
	resp, body, err := c.requestHTTP(req, cancel)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if err.Error() != "cancelled" {
		t.Errorf("unexpected error: %v", err)
	}
	if resp != nil {
		t.Errorf("unexpected non-nil response: %v", resp)
	}
	if body != nil {
		t.Errorf("unexpected non-nil body: %q", body)
	}
}
