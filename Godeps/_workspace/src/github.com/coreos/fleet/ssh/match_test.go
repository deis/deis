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

package ssh

import "testing"

func TestMatchHost(t *testing.T) {
	tests := []struct {
		host    string
		pattern string
		success bool
	}{
		{"foo.com", "foo.com", true},
		{"foo.com", "foo.com,bar,baz", true},
		{"foo.com", "bar,foo.com,baz", true},
		{"foo.com", "!foo.com,bar,baz", false},
		{"foo.com", "foo.com,!foo.com,bar,baz", false},
	}
	for _, test := range tests {
		if matchHost(test.host, test.pattern) != test.success {
			t.Errorf("matching %v against %v: got %v, want %v!", test.host, test.pattern, test.success, !test.success)
		}
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern string
		good    []string
		bad     []string
	}{
		{"foo",
			[]string{"foo"},
			[]string{"fob", "FOO"},
		},
		{"f*o",
			[]string{"foo", "fooo", "fasdfo"},
			[]string{"foc"},
		},
		{"f*",
			[]string{"fasdf", "f0939", "fa"},
			[]string{"g", "asdfadsff"},
		},
		{"f**",
			[]string{"fasdf", "f0939", "fa"},
			[]string{"g", "asdfasdfff"},
		},
		{"a*c?e",
			[]string{"abcde", "accce", "a2c3e", "abbcde", "acde"},
			[]string{"abce", "ace", "abbbbbcdde"},
		},
		{"a*c*e",
			[]string{"abcde", "accce", "a2c3e", "abbcde", "acde", "abbbcdddde"},
			[]string{"abc", "ae"},
		},
		{"a**c*e",
			[]string{"abcde", "accce", "a2c3e", "abbcde", "acde", "abbbcdddde"},
			[]string{"abc", "ae"},
		},
		{"f?b",
			[]string{"fob", "fab"},
			[]string{"fb", "feeb"},
		},
		{"h??d",
			[]string{"hood", "heed", "h12d"},
			[]string{"heldd", "hoob"},
		},
		{"a?c?e",
			[]string{"abcde", "accce", "a2c3e"},
			[]string{"abbcde", "abce", "ace"},
		},
	}

	for _, test := range tests {
		for _, s := range test.good {
			if !match(s, test.pattern) {
				t.Errorf("%v failed match against %v!", s, test.pattern)
			}
		}
		for _, s := range test.bad {
			if match(s, test.pattern) {
				t.Errorf("%v did not fail match against %v!", s, test.pattern)
			}
		}
	}

}
