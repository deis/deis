package machine

import (
	"testing"
)

func TestHasMetadataSimpleMatch(t *testing.T) {
	metadata := map[string]string{
		"region": "us-east-1",
	}
	ms := &MachineState{Metadata: metadata}

	match := map[string][]string{
		"region": {"us-east-1"},
	}
	if !HasMetadata(ms, match) {
		t.Errorf("Machine reported it did not have expected state")
	}
}

func TestHasMetadataMultiMatch(t *testing.T) {
	metadata := map[string]string{
		"groups": "ping",
	}
	ms := &MachineState{Metadata: metadata}

	match := map[string][]string{
		"groups": {"ping", "pong"},
	}
	if !HasMetadata(ms, match) {
		t.Errorf("Machine reported it did not have expected state")
	}
}

func TestHasMetadataSingleMatchFail(t *testing.T) {
	metadata := map[string]string{
		"groups": "ping",
	}
	ms := &MachineState{Metadata: metadata}

	match := map[string][]string{
		"groups": {"pong"},
	}
	if HasMetadata(ms, match) {
		t.Errorf("Machine reported a successful match for metadata which it does not have")
	}
}

func TestHasMetadataPartialMatchFail(t *testing.T) {
	metadata := map[string]string{
		"region": "us-east-1",
		"groups": "ping",
	}
	ms := &MachineState{Metadata: metadata}

	match := map[string][]string{
		"region": {"us-east-1"},
		"groups": {"pong"},
	}
	if HasMetadata(ms, match) {
		t.Errorf("Machine reported a successful match for metadata which it does not have")
	}
}
