// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package machine

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestReadLocalMachineIDMissing(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "fleet-")
	if err != nil {
		t.Fatalf("Failed creating tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	machID, err := readLocalMachineID(dir)
	if err == nil {
		t.Fatal("Expected error for missing machID, but got nil")
	}
	if machID != "" {
		t.Fatalf("Received incorrect machID: %s", machID)
	}
}

func TestReadLocalMachineIDFound(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "fleet-")
	if err != nil {
		t.Fatalf("Failed creating tempdir: %v", err)
	}
	defer os.RemoveAll(dir)

	tmpMachineIDPath := filepath.Join(dir, "/etc/machine-id")
	err = os.MkdirAll(filepath.Dir(tmpMachineIDPath), os.FileMode(0755))
	if err != nil {
		t.Fatalf("Failed setting up fake mach ID path: %v", err)
	}

	err = ioutil.WriteFile(tmpMachineIDPath, []byte("pingpong"), os.FileMode(0644))
	if err != nil {
		t.Fatalf("Failed writing fake mach ID file: %v", err)
	}

	machID, err := readLocalMachineID(dir)
	if err != nil {
		t.Fatalf("Unexpected error reading machID: %v", err)
	}
	if machID != "pingpong" {
		t.Fatalf("Received incorrect machID %q, expected 'pingpong'", machID)
	}
}

func TestUsableAddress(t *testing.T) {
	tests := []struct {
		ip net.IP
		ok bool
	}{
		// unicast IPv4 usable
		{net.ParseIP("192.168.1.12"), true},

		// unicast IPv6 unusable
		{net.ParseIP("2001:DB8::3"), false},

		// loopback IPv4/6 unusable
		{net.ParseIP("127.0.0.12"), false},
		{net.ParseIP("::1"), false},

		// link-local IPv4/6 unusable
		{net.ParseIP("169.254.4.87"), false},
		{net.ParseIP("fe80::12"), false},

		// unspecified (all zeros) IPv4/6 unusable
		{net.ParseIP("0.0.0.0"), false},
		{net.ParseIP("::"), false},

		// multicast IPv4/6 unusable
		{net.ParseIP("239.255.255.250"), false},
		{net.ParseIP("ffx2::4"), false},
	}

	for i, tt := range tests {
		ok := usableAddress(tt.ip)
		if tt.ok != ok {
			t.Errorf("case %d: expected %v usable %t, got %t", i, tt.ip, tt.ok, ok)
		}
	}
}
