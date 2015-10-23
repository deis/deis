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
	"reflect"
	"testing"
)

func TestUnmarshalFailedResponse(t *testing.T) {
	body := "{\"errorCode\":401,\"message\":\"foo\",\"cause\":\"bar\",\"index\":12}"
	expect := Error{ErrorCode: 401, Message: "foo", Cause: "bar", Index: 12}

	res, err := unmarshalFailedResponse(nil, []byte(body))
	if res != nil {
		t.Errorf("*Result should always be nil")
	}

	etcdErr, ok := err.(Error)
	if !ok {
		t.Fatalf("error should be of type Error")
	}

	if !reflect.DeepEqual(etcdErr, err) {
		t.Fatalf("returned err %v does not match expected %v", etcdErr, expect)
	}
}

func TestUnmarshalFailedResponseGarbage(t *testing.T) {
	res, err := unmarshalFailedResponse(nil, []byte("garbage"))
	if res != nil {
		t.Errorf("*Result should always be nil")
	}

	if _, ok := err.(Error); ok {
		t.Fatalf("error should not be of type Error")
	}
}
