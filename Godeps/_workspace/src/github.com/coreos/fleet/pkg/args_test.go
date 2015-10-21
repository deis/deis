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

import "testing"

func TestTrimToDashes(t *testing.T) {
	var argtests = []struct {
		input  []string
		output []string
	}{
		{[]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}},
		{[]string{"abc", "def", "--", "ghi"}, []string{"ghi"}},
		{[]string{"abc", "def", "--"}, []string{}},
		{[]string{"--"}, []string{}},
		{[]string{"--", "abc", "def", "ghi"}, []string{"abc", "def", "ghi"}},
		{[]string{"--", "bar", "--", "baz"}, []string{"bar", "--", "baz"}},
		{[]string{"--flagname", "--", "ghi"}, []string{"ghi"}},
		{[]string{"--", "--flagname", "ghi"}, []string{"--flagname", "ghi"}},
	}
	for _, test := range argtests {
		args := TrimToDashes(test.input)
		if len(test.output) != len(args) {
			t.Fatalf("error trimming dashes: expected %s, got %s", test.output, args)
		}
		for i, v := range args {
			if v != test.output[i] {
				t.Fatalf("error trimming dashes: expected %s, got %s", test.output, args)
			}
		}
	}
}
