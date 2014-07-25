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
