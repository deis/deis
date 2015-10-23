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
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestGetHTTPRequest(t *testing.T) {
	tests := []actionTestCase{
		{
			&Get{Key: "/foo"},
			"GET",
			"/v2/keys/foo?consistent=true&recursive=false&sorted=false",
			"",
		},
		{
			&Get{Key: "/foo", Sorted: true, Recursive: true},
			"GET",
			"/v2/keys/foo?consistent=true&recursive=true&sorted=true",
			"",
		},
	}

	driveActionTestCases(t, tests)
}

func TestCreateHTTPRequest(t *testing.T) {
	tests := []actionTestCase{
		{
			&Create{Key: "/foo"},
			"PUT",
			"/v2/keys/foo?prevExist=false",
			"value=",
		},
		{
			&Create{Key: "/foo", TTL: 5 * time.Minute},
			"PUT",
			"/v2/keys/foo?prevExist=false",
			"ttl=300&value=",
		},
		{
			&Create{Key: "/foo", Value: "bar"},
			"PUT",
			"/v2/keys/foo?prevExist=false",
			"value=bar",
		},
		{
			&Create{Key: "/foo", TTL: 5 * time.Minute, Value: "bar"},
			"PUT",
			"/v2/keys/foo?prevExist=false",
			"ttl=300&value=bar",
		},
	}

	driveActionTestCases(t, tests)
}

func TestUpdateHTTPRequest(t *testing.T) {
	tests := []actionTestCase{
		{
			&Update{Key: "/foo"},
			"PUT",
			"/v2/keys/foo?prevExist=true",
			"value=",
		},
		{
			&Update{Key: "/foo", TTL: 5 * time.Minute},
			"PUT",
			"/v2/keys/foo?prevExist=true",
			"ttl=300&value=",
		},
		{
			&Update{Key: "/foo", Value: "bar"},
			"PUT",
			"/v2/keys/foo?prevExist=true",
			"value=bar",
		},
		{
			&Update{Key: "/foo", TTL: 5 * time.Minute, Value: "bar"},
			"PUT",
			"/v2/keys/foo?prevExist=true",
			"ttl=300&value=bar",
		},
	}

	driveActionTestCases(t, tests)
}

func TestSetHTTPRequest(t *testing.T) {
	tests := []actionTestCase{
		{
			&Set{Key: "/foo"},
			"PUT",
			"/v2/keys/foo",
			"value=",
		},
		{
			&Set{Key: "/foo", TTL: 5 * time.Minute},
			"PUT",
			"/v2/keys/foo",
			"ttl=300&value=",
		},
		{
			&Set{Key: "/foo", Value: "bar"},
			"PUT",
			"/v2/keys/foo",
			"value=bar",
		},
		{
			&Set{Key: "/foo", TTL: 5 * time.Minute, Value: "bar"},
			"PUT",
			"/v2/keys/foo",
			"ttl=300&value=bar",
		},
		{
			&Set{Key: "/foo", Value: "bar", PreviousIndex: 13},
			"PUT",
			"/v2/keys/foo?prevIndex=13",
			"value=bar",
		},
		{
			&Set{Key: "/foo", Value: "bar", PreviousValue: "baz"},
			"PUT",
			"/v2/keys/foo?prevValue=baz",
			"value=bar",
		},
	}

	driveActionTestCases(t, tests)
}

func TestWatchHTTPRequest(t *testing.T) {
	tests := []actionTestCase{
		{
			&Watch{Key: "/foo"},
			"GET",
			"/v2/keys/foo?consistent=true&recursive=false&wait=true",
			"",
		},
		{
			&Watch{Key: "/foo", WaitIndex: 12},
			"GET",
			"/v2/keys/foo?consistent=true&recursive=false&wait=true&waitIndex=12",
			"",
		},
		{
			&Watch{Key: "/foo", Recursive: true},
			"GET",
			"/v2/keys/foo?consistent=true&recursive=true&wait=true",
			"",
		},
	}

	driveActionTestCases(t, tests)
}

func TestDeleteHTTPRequest(t *testing.T) {
	tests := []actionTestCase{
		{
			&Delete{Key: "/foo"},
			"DELETE",
			"/v2/keys/foo?recursive=false",
			"",
		},
		{
			&Delete{Key: "/foo", PreviousValue: "bar"},
			"DELETE",
			"/v2/keys/foo?prevValue=bar&recursive=false",
			"",
		},
		{
			&Delete{Key: "/foo", PreviousIndex: uint64(12)},
			"DELETE",
			"/v2/keys/foo?prevIndex=12&recursive=false",
			"",
		},
		{
			&Delete{Key: "/foo", PreviousValue: "bar", PreviousIndex: uint64(12)},
			"DELETE",
			"/v2/keys/foo?prevIndex=12&prevValue=bar&recursive=false",
			"",
		},
	}

	driveActionTestCases(t, tests)
}

type actionTestCase struct {
	act    Action
	method string
	url    string
	body   string
}

func driveActionTestCases(t *testing.T, tests []actionTestCase) {
	for i, test := range tests {
		req, err := test.act.HTTPRequest()
		if err != nil {
			t.Errorf("%d: HTTPRequest returned unexpected error: %v", i, err)
			continue
		}

		if req.Method != test.method {
			t.Errorf("%d: request method is %s, expected %s", i, req.Method, test.method)
		}

		if req.URL.String() != test.url {
			t.Errorf("%d: request URL is %s, expected %s", i, req.URL.String(), test.url)
		}

		var body string
		if req.Body != nil {
			bbody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Errorf("%d: failed reading request body: %v", i, err)
				continue
			}
			body = fmt.Sprintf("%s", bbody)
		}

		if body != test.body {
			t.Errorf("%d: request body is %q, expected %q", i, body, test.body)
		}
	}
}

func TestV2URLHelper(t *testing.T) {
	tests := []struct {
		input  string
		expect url.URL
	}{
		{"", url.URL{Path: "/v2/keys"}},
		{"/", url.URL{Path: "/v2/keys"}},
		{"/foo", url.URL{Path: "/v2/keys/foo"}},
		{"/foo/bar", url.URL{Path: "/v2/keys/foo/bar"}},
		{"/space space/literal%", url.URL{Path: "/v2/keys/space space/literal%"}},
	}

	for i, tt := range tests {
		output := v2URL(tt.input)
		if !reflect.DeepEqual(output, tt.expect) {
			t.Errorf("case %d: output=%s, expect=%s", i, output.String(), tt.expect.String())
		}
	}
}
