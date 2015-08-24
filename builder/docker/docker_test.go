package docker

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Masterminds/cookoo"
	docli "github.com/fsouza/go-dockerclient"
)

func TestCleanup(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	tf, err := ioutil.TempFile("/tmp", "junkdockersock")
	if err != nil {
		t.Error("Could not create temp file for testing: " + err.Error())
	}
	DockSock = tf.Name()
	tf.Close()

	reg.Route("test", "Test route").
		Does(Cleanup, "res")

	if err := router.HandleRequest("test", cxt, true); err != nil {
		t.Error(err)
	}

	// From the TempFile docs: "It is the caller's responsibility to remove the temp
	// file." So we can reasonably assume that if the file's gone, it's because
	// it was deleted by Cleanup.
	if _, err := os.Stat(DockSock); err == nil {
		t.Errorf("Expected file %s to be deleted, but got a stat.", DockSock)
	}
}

func TestCreateClient(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(CreateClient, "res").Using("url").WithDefault("http://example.com:4321")

	if err := router.HandleRequest("test", cxt, true); err != nil {
		t.Error(err)
	}

	if cli := cxt.Get("res", nil); cli == nil {
		t.Error("Expected a client")
	} else if _, ok := cli.(*docli.Client); !ok {
		t.Error("Expected client to be a *docker.Cli")
	}

}
