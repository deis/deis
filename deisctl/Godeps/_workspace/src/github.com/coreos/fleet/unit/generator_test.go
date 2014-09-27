package unit

import (
	"reflect"
	"testing"
)

func assertGenerateUnitStateHeartbeats(t *testing.T, um UnitManager, gen *UnitStateGenerator, expect []UnitStateHeartbeat) {
	beatchan, err := gen.Generate()
	if err != nil {
		t.Fatalf("Unexpected error from Generate(): %v", err)
	}

	got := []UnitStateHeartbeat{}
	for beat := range beatchan {
		beat := beat
		got = append(got, *beat)
	}

	if !reflect.DeepEqual(got, expect) {
		t.Fatalf("got %#v, expected %#v", got, expect)
	}
}

func TestUnitStateGeneratorSubscribeLifecycle(t *testing.T) {
	um := NewFakeUnitManager()
	um.Load("foo.service", UnitFile{})

	gen := NewUnitStateGenerator(um)

	// not subscribed to anything yet, so no heartbeats
	assertGenerateUnitStateHeartbeats(t, um, gen, []UnitStateHeartbeat{})

	gen.Subscribe("foo.service")

	// subscribed to foo.service so we should get a heartbeat
	expect := []UnitStateHeartbeat{
		UnitStateHeartbeat{Name: "foo.service", State: &UnitState{"loaded", "active", "running", "", "", "foo.service"}},
	}
	assertGenerateUnitStateHeartbeats(t, um, gen, expect)

	gen.Unsubscribe("foo.service")

	// heartbeat for foo.service should have nil State since we have not called Generate since Unsubscribe
	expect = []UnitStateHeartbeat{
		UnitStateHeartbeat{Name: "foo.service", State: nil},
	}
	assertGenerateUnitStateHeartbeats(t, um, gen, expect)

	// since the nil-State heartbeat for foo.service was sent for the last Generate, it can be forgotten
	assertGenerateUnitStateHeartbeats(t, um, gen, []UnitStateHeartbeat{})
}

func TestUnitStateGeneratorNoState(t *testing.T) {
	um := NewFakeUnitManager()
	gen := NewUnitStateGenerator(um)

	// not subscribed to anything yet, so no heartbeats
	assertGenerateUnitStateHeartbeats(t, um, gen, []UnitStateHeartbeat{})

	gen.Subscribe("foo.service")

	// subscribed to foo.service but no underlying state so no heartbeat
	assertGenerateUnitStateHeartbeats(t, um, gen, []UnitStateHeartbeat{})
}
