package git

import (
	"testing"
)

func TestRemoteURL(t *testing.T) {
	t.Parallel()

	expected := "ssh://git@example.com:2222/app.git"

	actual := RemoteURL("example.com", "app")

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}
