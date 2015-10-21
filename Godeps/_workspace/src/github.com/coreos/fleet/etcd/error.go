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
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ErrorKeyNotFound       = 100
	ErrorNodeExist         = 105
	ErrorEventIndexCleared = 401
)

type Error struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Cause     string `json:"cause"`
	Index     uint64 `json:"index"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v: %v (%v) [%v]", e.ErrorCode, e.Message, e.Cause, e.Index)
}

func unmarshalFailedResponse(resp *http.Response, body []byte) (*Result, error) {
	var etcdErr Error
	err := json.Unmarshal(body, &etcdErr)
	if err != nil {
		return nil, err
	}

	return nil, etcdErr
}

func IsKeyNotFound(err error) bool {
	e, ok := err.(Error)
	return ok && e.ErrorCode == ErrorKeyNotFound
}

func IsNodeExist(err error) bool {
	e, ok := err.(Error)
	return ok && e.ErrorCode == ErrorNodeExist
}
