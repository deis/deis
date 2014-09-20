package fleet

import "testing"

func TestFormatUnitName(t *testing.T) {
	unitName, err := formatUnitName("router", 1)
	if err != nil {
		t.Fatal(err)
	}
	if unitName != "deis-router@1.service" {
		t.Fatalf("invalid unit name: %v", unitName)
	}

}
