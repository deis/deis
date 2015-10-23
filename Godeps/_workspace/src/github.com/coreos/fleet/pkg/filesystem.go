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
	"io/ioutil"
)

// ListDirectory generates a slice of all the file names that both exist in
// the provided directory and pass the filter.
// The returned file names are relative to the directory argument.
// filterFunc is called once for each file found in the directory. If
// filterFunc returns true, the given file will ignored.
func ListDirectory(dir string, filterFunc func(string) bool) ([]string, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	units := make([]string, 0)
	for _, fi := range fis {
		name := fi.Name()
		if filterFunc(name) {
			continue
		}
		units = append(units, name)
	}

	return units, nil
}
