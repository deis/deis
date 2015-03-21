package builder

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	dtime "github.com/deis/deis/pkg/time"
)

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() error {
	// we don't have to do anything here, since the buffer is just some data in memory
	return nil
}

func stringInSlice(list []string, s string) bool {
	for _, li := range list {
		if li == s {
			return true
		}
	}
	return false
}

func TestYamlToJSONGood(t *testing.T) {
	goodProcfiles := [][]byte{
		[]byte(`web: while true; do echo hello; sleep 1; done`),

		[]byte(`web: while true; do echo hello; sleep 1; done
worker: while true; do echo hello; sleep 1; done`),
		// test a procfile with quoted strings
		[]byte(`web: /bin/bash -c "while true; do echo hello; sleep 1; done"`),
	}

	goodProcess := "while true; do echo hello; sleep 1; done"

	for _, procfile := range goodProcfiles {
		data, err := YamlToJSON(procfile)
		if err != nil {
			t.Errorf("expected procfile to be valid, got '%v'", err)
		}
		var p ProcessType
		if err := json.Unmarshal([]byte(data), &p); err != nil {
			t.Errorf("expected to be able to unmarshal object, got '%v'", err)
		}
		if !strings.Contains(p["web"], goodProcess) {
			t.Errorf("expected web process == '%s', got '%s'", goodProcess, p["web"])
		}
	}
}

func TestParseConfigGood(t *testing.T) {
	// mock the controller response
	resp := bytes.NewBufferString(`{"owner": "test",
			"app": "example-go",
			"values": {"FOO": "bar", "CAR": 1234},
			"memory": {},
			"cpu": {},
			"tags": {},
			"created": "2014-01-01T00:00:00UTC",
			"updated": "2014-01-01T00:00:00UTC",
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"}`)

	config, err := ParseConfig(resp.Bytes())

	if err != nil {
		t.Error(err)
	}

	if config.Values["FOO"] != "bar" {
		t.Errorf("expected FOO='bar', got FOO='%v'", config.Values["FOO"])
	}

	if car, ok := config.Values["CAR"].(float64); ok {
		if car != 1234 {
			t.Errorf("expected CAR=1234, got CAR=%d", config.Values["CAR"])
		}
	} else {
		t.Error("expected CAR to be of type float64")
	}
}

func TestParseDomainGood(t *testing.T) {
	// mock controller build-hook response
	resp := []byte(`{"release": {"version": 1},
"domains": ["test.example.com", "test2.example.com"]}`)

	domain, err := ParseDomain(resp)
	if err != nil {
		t.Errorf("expected to parse domain, got '%v'", err)
	}
	if domain != "test.example.com" {
		t.Errorf("expected 'test.example.com', got '%s'", domain)
	}
}

func TestParseReleaseVersionGood(t *testing.T) {
	// mock controller build-hook response
	resp := []byte(`{"release": {"version": 1},
"domains": ["test.example.com", "test2.example.com"]}`)

	version, err := ParseReleaseVersion(resp)
	if err != nil {
		t.Errorf("expected to parse version, got '%v'", err)
	}
	if version != 1 {
		t.Errorf("expected '1', got '%d'", version)
	}
}

func TestGetDefaultTypeGood(t *testing.T) {
	goodData := [][]byte{[]byte(`default_process_types:
  web: while true; do echo hello; sleep 1; done`),
		[]byte(`foo: bar
default_process_types:
  web: while true; do echo hello; sleep 1; done`),
		[]byte(``)}

	for _, data := range goodData {
		defaultType, err := GetDefaultType(data)
		if err != nil {
			t.Error(err)
		}
		if defaultType != `{"web":"while true; do echo hello; sleep 1; done"}` && string(data) != "" {
			t.Errorf("incorrect default type, got %s", defaultType)
		}
		if string(data) == "" && defaultType != "{}" {
			t.Errorf("incorrect default type, got %s", defaultType)
		}
	}
}

func TestParseControllerConfigGood(t *testing.T) {
	// mock controller config response
	resp := []byte(`{"owner": "test",
		"app": "example-go",
		"values": {"FOO": "bar", "CAR": "star"},
		"memory": {},
		"cpu": {},
		"tags": {},
		"created": "2014-01-01T00:00:00UTC",
		"updated": "2014-01-01T00:00:00UTC",
		"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
	}`)

	config, err := ParseControllerConfig(resp)

	if err != nil {
		t.Errorf("expected to pass, got '%v'", err)
	}

	if len(config) != 2 {
		t.Errorf("expected 2, got %d", len(config))
	}

	if !stringInSlice(config, " -e CAR=\"star\"") {
		t.Error("expected ' -e CAR=\"star\"' in slice")
	}
}

func TestTimeSerialize(t *testing.T) {
	time, err := json.Marshal(&dtime.Time{Time: time.Now().UTC()})

	if err != nil {
		t.Errorf("expected to be able to serialize time as json, got '%v'", err)
	}

	if !strings.Contains(string(time), "UTC") {
		t.Errorf("could not find 'UTC' in datetime, got '%s'", string(time))
	}
}
