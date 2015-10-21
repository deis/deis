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
	"time"
)

type Result struct {
	Action   string `json:"action"`
	Node     *Node  `json:"node"`
	PrevNode *Node  `json:"prevNode"`
	Raw      []byte `json:"-"`
}

func (r *Result) String() string {
	return fmt.Sprintf("{Action: %s, Node: %v, PrevNode: %v}", r.Action, r.Node, r.PrevNode)
}

type Nodes []Node
type Node struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	TTL           int    `json:"ttl"`
	Nodes         Nodes  `json:"nodes"`
	ModifiedIndex uint64 `json:"modifiedIndex"`
	CreatedIndex  uint64 `json:"createdIndex"`
}

func (n Node) TTLDuration() time.Duration {
	dur := time.Duration(n.TTL) * time.Second
	if dur < 0 {
		dur = 0
	}
	return dur
}

func (n *Node) String() string {
	return fmt.Sprintf("{Key: %s, CreatedIndex: %d, ModifiedIndex: %d}", n.Key, n.CreatedIndex, n.ModifiedIndex)
}

func unmarshalSuccessfulResponse(resp *http.Response, body []byte) (*Result, error) {
	var res Result
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	res.Raw = body
	return &res, nil
}
