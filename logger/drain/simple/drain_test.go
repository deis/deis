package simple

import (
	"fmt"
	"testing"
)

func TestInvalidDrainUrl(t *testing.T) {
	_, err := NewDrain("https://wwww.google.com")
	if err == nil || err.Error() != fmt.Sprintf("Invalid drain url scheme: %s", "https") {
		t.Error("Did not receive expected error message")
	}
}
