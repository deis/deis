package confd

import (
	"testing"

	"github.com/Masterminds/cookoo"
)

func TestRunOnce(t *testing.T) {
	reg, _, _ := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(RunOnce, "res").
		Using("node").WithDefault("localhost:4001")

	/*
		err := router.HandleRequest("test", cxt, true)
		if err != nil {
			t.Error(err)
		}
	*/
}
func TestRun(t *testing.T) {
	reg, _, _ := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(Run, "res").
		Using("node").WithDefault("localhost:4001").
		Using("interval").WithDefault(200)

	/*
		err := router.HandleRequest("test", cxt, true)
		if err != nil {
			t.Error(err)
		}
	*/
}
