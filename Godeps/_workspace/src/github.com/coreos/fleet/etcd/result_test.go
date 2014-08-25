package etcd

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshalSuccessfulResponseNoNodes(t *testing.T) {
	tests := []struct {
		resp        http.Response
		res         *Result
		expectError bool
	}{
		// Neither PrevNode or Node
		{
			http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"action":"delete"}`)),
			},
			&Result{Action: "delete"},
			false,
		},

		// PrevNode
		{
			http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"action":"delete", "prevNode": {"key": "/foo", "value": "bar", "modifiedIndex": 12, "createdIndex": 10}}`)),
			},
			&Result{Action: "delete", PrevNode: &Node{Key: "/foo", Value: "bar", ModifiedIndex: 12, CreatedIndex: 10}},
			false,
		},

		// Node
		{
			http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"action":"get", "node": {"key": "/foo", "value": "bar", "modifiedIndex": 12, "createdIndex": 10}}`)),
			},
			&Result{Action: "get", Node: &Node{Key: "/foo", Value: "bar", ModifiedIndex: 12, CreatedIndex: 10}},
			false,
		},

		// PrevNode and Node
		{
			http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"action":"update", "prevNode": {"key": "/foo", "value": "baz", "modifiedIndex": 10, "createdIndex": 10}, "node": {"key": "/foo", "value": "bar", "modifiedIndex": 12, "createdIndex": 10}}`)),
			},
			&Result{Action: "update", PrevNode: &Node{Key: "/foo", Value: "baz", ModifiedIndex: 10, CreatedIndex: 10}, Node: &Node{Key: "/foo", Value: "bar", ModifiedIndex: 12, CreatedIndex: 10}},
			false,
		},

		// Garbage in body
		{
			http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`garbage`)),
			},
			nil,
			true,
		},
	}

	for i, tt := range tests {
		var body []byte
		if tt.resp.Body != nil {
			var err error
			body, err = ioutil.ReadAll(tt.resp.Body)
			if err != nil {
				t.Fatalf("case %d: failed preparing body: %v", i, err)
			}
		}
		res, err := unmarshalSuccessfulResponse(&tt.resp, body)
		if tt.expectError != (err != nil) {
			t.Errorf("case %d: expectError=%t, err=%v", i, tt.expectError, err)
		}

		if (res == nil) != (tt.res == nil) {
			t.Errorf("case %d: received res==%v, but expected res==%v", i, res, tt.res)
			continue
		} else if tt.res == nil {
			// expected and succesfully got nil response
			continue
		}

		if res.Action != tt.res.Action {
			t.Errorf("case %d: Action=%s, expected %s", i, res.Action, tt.res.Action)
		}

		if !reflect.DeepEqual(res.Node, tt.res.Node) {
			t.Errorf("case %d: Node=%v, expected %v", i, res.Node, tt.res.Node)
		}
	}
}
