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

package pkg

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestListDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "fleet-testing-")
	if err != nil {
		t.Fatal(err.Error())
	}

	defer os.RemoveAll(dir)

	for _, name := range []string{"ping", "pong", "foo", "bar", "baz"} {
		err := ioutil.WriteFile(path.Join(dir, name), []byte{}, 0400)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	filterFunc := func(name string) bool {
		return name == "foo" || name == "bar"
	}

	got, err := ListDirectory(dir, filterFunc)
	if err != nil {
		t.Fatal(err.Error())
	}

	want := []string{"baz", "ping", "pong"}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("ListDirectory output incorrect: want=%v, got=%v", want, got)
	}
}
