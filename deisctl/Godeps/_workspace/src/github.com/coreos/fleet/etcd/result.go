package etcd

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	Nodes         Nodes  `json:"nodes"`
	ModifiedIndex uint64 `json:"modifiedIndex"`
	CreatedIndex  uint64 `json:"createdIndex"`
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
