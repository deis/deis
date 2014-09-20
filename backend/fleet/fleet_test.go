package fleet

import "testing"

func TestNewClient(t *testing.T) {
	// set required flags
	Flags.Endpoint = "http://127.0.0.1:4001"

	// instantiate client
	_, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
}
