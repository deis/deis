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

package unit

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/coreos/go-systemd/unit"
)

func NewUnitFile(raw string) (*UnitFile, error) {
	reader := strings.NewReader(raw)
	opts, err := unit.Deserialize(reader)
	if err != nil {
		return nil, err
	}

	return NewUnitFromOptions(opts), nil
}

func NewUnitFromOptions(opts []*unit.UnitOption) *UnitFile {
	return &UnitFile{mapOptions(opts), opts}
}

func mapOptions(opts []*unit.UnitOption) map[string]map[string][]string {
	contents := make(map[string]map[string][]string)
	for _, opt := range opts {
		if _, ok := contents[opt.Section]; !ok {
			contents[opt.Section] = make(map[string][]string)
		}

		if _, ok := contents[opt.Section][opt.Name]; !ok {
			contents[opt.Section][opt.Name] = make([]string, 0)
		}

		var vals []string
		if opt.Section == "X-Fleet" {
			// The go-systemd parser does not know that our X-Fleet options support
			// multivalue options, so we have to manually parse them here
			vals = parseMultivalueLine(opt.Value)
		} else {
			vals = []string{opt.Value}
		}

		contents[opt.Section][opt.Name] = append(contents[opt.Section][opt.Name], vals...)
	}

	return contents
}

// parseMultivalueLine parses a line that includes several quoted values separated by whitespaces.
// Example: MachineMetadata="foo=bar" "baz=qux"
// If the provided line is not a multivalue string, the line is returned as the sole value.
func parseMultivalueLine(line string) (values []string) {
	if !strings.HasPrefix(line, `"`) || !strings.HasSuffix(line, `"`) {
		return []string{line}
	}

	var v bytes.Buffer
	w := false // check whether we're within quotes or not
	for _, e := range []byte(line) {
		// ignore quotes
		if e == '"' {
			w = !w
			continue
		}

		if e == ' ' {
			if !w { // between quoted values, keep the previous value and reset.
				values = append(values, v.String())
				v.Reset()
				continue
			}
		}

		v.WriteByte(e)
	}

	values = append(values, v.String())

	return
}

// A UnitFile represents a systemd configuration which encodes information about any of the unit
// types that fleet supports (as defined in SupportedUnitTypes()).
// UnitFiles are linked to Units by the Hash of their contents.
// Similar to systemd, a UnitFile configuration has no inherent name, but is rather
// named through the reference to it; in the case of systemd, the reference is
// the filename, and in the case of fleet, the reference is the name of the Unit
// that references this UnitFile.
type UnitFile struct {
	// Contents represents the parsed unit file.
	// This field must be considered readonly.
	Contents map[string]map[string][]string

	Options []*unit.UnitOption
}

// Description returns the first Description option found in the [Unit] section.
// If the option is not defined, an empty string is returned.
func (u *UnitFile) Description() string {
	if values := u.Contents["Unit"]["Description"]; len(values) > 0 {
		return values[0]
	}
	return ""
}

func (u *UnitFile) Bytes() []byte {
	b, _ := ioutil.ReadAll(unit.Serialize(u.Options))
	return b
}

func (u *UnitFile) String() string {
	return string(u.Bytes())
}

// Hash returns the SHA1 hash of the raw contents of the Unit
func (u *UnitFile) Hash() Hash {
	return Hash(sha1.Sum(u.Bytes()))
}

// RecognizedUnitType determines whether or not the given unit name represents
// a recognized unit type.
func RecognizedUnitType(name string) bool {
	types := []string{"service", "socket", "timer", "path", "device", "mount", "automount"}
	for _, t := range types {
		suffix := fmt.Sprintf(".%s", t)
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}
	return false
}

// DefaultUnitType appends the default unit type to a given unit name, ignoring
// any file extensions that already exist.
func DefaultUnitType(name string) string {
	return fmt.Sprintf("%s.service", name)
}

// SHA1 sum
type Hash [sha1.Size]byte

func (h Hash) String() string {
	return fmt.Sprintf("%x", h[:])
}

func (h Hash) Short() string {
	return fmt.Sprintf("%.7s", h)
}

func (h *Hash) Empty() bool {
	return *h == Hash{}
}

// UnitState encodes the current state of a unit loaded into a fleet agent
type UnitState struct {
	LoadState   string
	ActiveState string
	SubState    string
	MachineID   string
	UnitHash    string
	UnitName    string
}

func NewUnitState(loadState, activeState, subState, mID string) *UnitState {
	return &UnitState{
		LoadState:   loadState,
		ActiveState: activeState,
		SubState:    subState,
		MachineID:   mID,
	}
}

// UnitNameInfo exposes certain interesting items about a Unit based on its
// name. For example, a unit with the name "foo@.service" constitutes a
// template unit, and a unit named "foo@1.service" would represent an instance
// unit of that template.
type UnitNameInfo struct {
	FullName string // Original complete name of the unit (e.g. foo.socket, foo@bar.service)
	Name     string // Name of the unit without suffix (e.g. foo, foo@bar)
	Prefix   string // Prefix of the template unit (e.g. foo)

	// If the unit represents an instance or a template, the following values are set
	Template string // Name of the canonical template unit (e.g. foo@.service)
	Instance string // Instance name (e.g. bar)
}

// IsInstance returns a boolean indicating whether the UnitNameInfo appears to be
// an Instance of a Template unit
func (nu UnitNameInfo) IsInstance() bool {
	return len(nu.Instance) > 0
}

// NewUnitNameInfo generates a UnitNameInfo from the given name. If the given string
// is not a correct unit name, nil is returned.
func NewUnitNameInfo(un string) *UnitNameInfo {

	// Everything past the first @ and before the last . is the instance
	s := strings.LastIndex(un, ".")
	if s == -1 {
		return nil
	}

	nu := &UnitNameInfo{FullName: un}
	name := un[:s]
	suffix := un[s:]
	nu.Name = name

	a := strings.Index(name, "@")
	if a == -1 {
		// This does not appear to be a template or instance unit.
		nu.Prefix = name
		return nu
	}

	nu.Prefix = name[:a]
	nu.Template = fmt.Sprintf("%s@%s", name[:a], suffix)
	nu.Instance = name[a+1:]
	return nu
}
