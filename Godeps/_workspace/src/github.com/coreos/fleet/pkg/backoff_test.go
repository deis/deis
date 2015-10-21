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
	"testing"
	"time"
)

func TestExpBackoff(t *testing.T) {
	tests := []struct {
		last time.Duration
		max  time.Duration
		next time.Duration
	}{
		{1 * time.Second, 10 * time.Second, 2 * time.Second},
		{8 * time.Second, 10 * time.Second, 10 * time.Second},
		{10 * time.Second, 10 * time.Second, 10 * time.Second},
		{20 * time.Second, 10 * time.Second, 10 * time.Second},
	}

	for i, tt := range tests {
		next := ExpBackoff(tt.last, tt.max)
		if next != tt.next {
			t.Errorf("case %d: last=%v, max=%v, next=%v; got next=%v", i, tt.last, tt.max, tt.next, next)
		}
	}
}
