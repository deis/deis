package git

import (
	"testing"
)

func TestRemoteURL(t *testing.T) {
	t.Parallel()

	actual := RemoteURL("example.com", "app")
	expected := "ssh://git@example.com:2222/app.git"

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}

	actual = RemoteURL("deis.10.245.1.3.xip.io:31350", "velcro-underdog")
	expected = "ssh://git@deis.10.245.1.3.xip.io:2222/velcro-underdog.git"

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}
