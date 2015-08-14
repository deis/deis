package builder

import (
	"testing"
	"time"

	"github.com/Masterminds/cookoo"
)

func TestSleep(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	reg.Route("test", "Test route").
		Does(Sleep, "res").Using("duration").WithDefault(3 * time.Second)

	start := time.Now()
	err := router.HandleRequest("test", cxt, true)
	if err != nil {
		t.Error(err)
	}
	end := time.Now()

	if end.Sub(start) < 3*time.Second {
		t.Error("Expected elapsed time to be 3 seconds.")
	}

}
