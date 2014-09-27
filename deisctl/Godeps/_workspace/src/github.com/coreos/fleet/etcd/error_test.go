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
