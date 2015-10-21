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

package machine

import (
	"reflect"
	"testing"
)

func TestFakeMachine(t *testing.T) {
	ms := MachineState{ID: "XXX"}
	fm := FakeMachine{ms}

	ret := fm.State()
	if !reflect.DeepEqual(ms, ret) {
		t.Fatalf("FakeMachine.State() returned %v, expected %v", ret, ms)
	}
}
