package etcd

import (
	"testing"

	"github.com/Masterminds/cookoo"
	"github.com/coreos/go-etcd/etcd"
)

func TestInterfaces(t *testing.T) {
	// Throughout the codebase, we assume that our interfaces match what the
	// etcd client provides. This is a canary to verify.
	cli := etcd.NewClient([]string{"http://localhost:4001"})
	var _ Getter = cli
	var _ = cli
	var _ GetterSetter = cli
}

func TestCreateClient(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(CreateClient, "res").Using("url").WithDefault("localhost:4100")

	if err := router.HandleRequest("test", cxt, true); err != nil {
		t.Error(err)
	}

	// All we really want to know is whether we got a valid client back.
	_ = cxt.Get("res", nil).(*etcd.Client)
}

func TestGet(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(Get, "res").
		Using("client").WithDefault(&stubClient{}).
		Using("path").WithDefault("/")

	err := router.HandleRequest("test", cxt, true)
	if err != nil {
		t.Error(err)
	}

	if res := cxt.Get("res", nil); res == nil {
		t.Error("Expected an *etcd.Response, not nil.")
	} else if tt, ok := res.(*etcd.Response); !ok {
		t.Errorf("Expected instance of *etcd.Response. Got %T", tt)
	}
}

func TestSet(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(Set, "res").
		Using("client").WithDefault(&stubClient{}).
		Using("key").WithDefault("Hello").
		Using("value").WithDefault("World")

	err := router.HandleRequest("test", cxt, true)
	if err != nil {
		t.Error(err)
	}

	if res := cxt.Get("res", nil); res == nil {
		t.Error("Expected an *etcd.Response, not nil.")
	} else if tt, ok := res.(*etcd.Response); !ok {
		t.Errorf("Expected instance of *etcd.Response. Got %T", tt)
	}
}
func TestMakeDir(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(MakeDir, "res").
		Using("client").WithDefault(&stubClient{}).
		Using("path").WithDefault("/deis/users/foo")

	err := router.HandleRequest("test", cxt, true)
	if err != nil {
		t.Error(err)
	}

	if res := cxt.Get("res", nil); res == nil {
		t.Error("Expected an *etcd.Response, not nil.")
	} else if tt, ok := res.(*etcd.Response); !ok {
		t.Errorf("Expected instance of *etcd.Response. Got %T", tt)
	}
}

// stubClient implements EtcdGetter and EtcdDirCreator
type stubClient struct {
}

func (s *stubClient) Get(key string, sort, recurse bool) (*etcd.Response, error) {
	return s.response("get"), nil
}

func (s *stubClient) CreateDir(key string, ttl uint64) (*etcd.Response, error) {
	return s.response("createdir"), nil
}

func (s *stubClient) Set(key string, value string, ttl uint64) (*etcd.Response, error) {
	return s.response("set"), nil
}

func (s *stubClient) response(a string) *etcd.Response {
	// This is totally fake data. It may or may not reflect what etcd really
	// returns.
	return &etcd.Response{
		Action:    a,
		EtcdIndex: 1,
		RaftIndex: 1,
		RaftTerm:  1,
		Node: &etcd.Node{
			Dir: true,
			Key: "/foo",
		},
	}

}
