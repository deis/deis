package fleet

import "testing"

func TestNextComponent(t *testing.T) {
	// test first component
	num, err := nextUnitNum([]string{})
	if err != nil {
		t.Fatal(err)
	}
	if num != 1 {
		t.Fatal("Invalid component number")
	}
	// test next component
	num, err = nextUnitNum([]string{"deis-router@1.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatal("Invalid component number")
	}
	// test last component
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@2.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 3 {
		t.Fatal("Invalid component number")
	}
	// test middle component
	num, err = nextUnitNum([]string{"deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 1 {
		t.Fatal("Invalid component number")
	}
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 2 {
		t.Fatalf("Invalid component number: %v", num)
	}
	num, err = nextUnitNum([]string{"deis-router@1.service", "deis-router@2.service", "deis-router@3.service"})
	if err != nil {
		t.Fatal(err)
	}
	if num != 4 {
		t.Fatal("Invalid component number")
	}
}

func TestSplitJobName(t *testing.T) {
	c, num, err := splitJobName("deis-router@1.service")
	if err != nil {
		t.Fatal(err)
	}
	if c != "router" || num != 1 {
		t.Fatalf("Invalid values: %v %v", c, num)
	}
}

func TestSplitTarget(t *testing.T) {
	c, num, err := splitTarget("router")
	if err != nil {
		t.Fatal(err)
	}
	if c != "router" && num != 0 {
		t.Fatalf("Invalid split on \"%v\": %v %v", "router", c, num)
	}

	c, num, err = splitTarget("router@3")
	if err != nil {
		t.Fatal(err)
	}
	if c != "router" || num != 3 {
		t.Fatalf("Invalid split on \"%v\": %v %v", "router@3", c, num)
	}

	c, num, err = splitTarget("database-data")
	if err != nil {
		t.Fatal(err)
	}
	if c != "database-data" || num != 0 {
		t.Fatalf("Invalid split on \"%v\": %v %v", "database-data", c, num)
	}

}
