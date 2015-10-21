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

package pkg

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/coreos/fleet/log"
)

// ParseFilepath expands ~ and ~user constructions.
// If user or $HOME is unknown, do nothing.
func ParseFilepath(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	i := strings.Index(path, "/")
	if i < 0 {
		i = len(path)
	}
	var home string
	if i == 1 {
		if home = os.Getenv("HOME"); home == "" {
			usr, err := user.Current()
			if err != nil {
				log.Debugf("Failed to get current home directory: %v", err)
				return path
			}
			home = usr.HomeDir
		}
	} else {
		usr, err := user.Lookup(path[1:i])
		if err != nil {
			log.Debugf("Failed to get %v's home directory: %v", path[1:i], err)
			return path
		}
		home = usr.HomeDir
	}
	path = filepath.Join(home, path[i:])
	return path
}
