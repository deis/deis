package drain

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetUsingInvalidValues(t *testing.T) {
	drainStrs := []string{"bogus"}
	for _, drainStr := range drainStrs {
		_, err := NewDrain(drainStr)
		if err == nil || err.Error() != fmt.Sprintf("Cannot construct a drain for URL: '%s'", drainStr) {
			t.Error("Did not receive expected error message")
		}
	}
}

func TestGetUsingEmptyString(t *testing.T) {
	d, err := NewDrain("")
	if err != nil {
		t.Error(err)
	}
	// Should return nil by default.  nil means no drain-- which is valid.
	if d != nil {
		dType := reflect.TypeOf(d).String()
		t.Errorf("Expected a nil drain, but got a %s", dType)
	}
}

func TestGetUdpDrain(t *testing.T) {
	d, err := NewDrain("udp://my-awesome-log-server:514")
	if err != nil {
		t.Error(err)
	}
	if want, got := "*simple.logDrain", reflect.TypeOf(d).String(); want != got {
		t.Errorf("Expected a %s, but got a %s", want, got)
	}
}

func TestGetSyslogDrain(t *testing.T) {
	d, err := NewDrain("syslog://my-awesome-log-server:514")
	if err != nil {
		t.Error(err)
	}
	if want, got := "*simple.logDrain", reflect.TypeOf(d).String(); want != got {
		t.Errorf("Expected a %s, but got a %s", want, got)
	}
}

func TestGetTcpDrain(t *testing.T) {
	d, err := NewDrain("tcp://my-awesome-log-server:12345")
	if err != nil {
		t.Error(err)
	}
	if want, got := "*simple.logDrain", reflect.TypeOf(d).String(); want != got {
		t.Errorf("Expected a %s, but got a %s", want, got)
	}
}
