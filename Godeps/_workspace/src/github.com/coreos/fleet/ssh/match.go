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

package ssh

import "strings"

// matchHost tries to match the given host name against a comma-separated
// sequence of subpatterns s (each possibly preceded by ! to indicate negation).
// It returns a boolean indicating whether or not a positive match was made.
// Any matched negations take precedence over any other possible matches in the
// pattern.
func matchHost(host, pattern string) bool {
	subpatterns := strings.Split(pattern, ",")
	found := false
	for _, s := range subpatterns {
		// If the host name matches a negated pattern, it is not
		// accepted even if it matches other patterns on that line.
		if strings.HasPrefix(s, "!") && match(host, s[1:]) {
			return false
		}
		// Otherwise, check for a normal match
		if match(host, s) {
			found = true
		}
	}
	// Return success if we found a positive match.  If there was a negative
	// match, we have already returned false and never get here.
	return found
}

// match compares the input string s to the pattern p, which may contain
// single and multi-character wildcards (? and * respectively). It returns a
// boolean indicating whether the string matches the pattern.
func match(s, p string) bool {
	var i, j int
	for i < len(p) {
		if p[i] == '*' {
			// Skip the asterisk.
			i++

			// If at end of pattern, accept immediately.
			if i == len(p) {
				return true
			}

			// If next character in pattern is known, optimize.
			if p[i] != '?' && p[i] != '*' {
				// Look for instances of the next character in
				// pattern, and try to match starting from those.
				for ; j < len(s); j++ {
					if s[j] == p[i] && match(s[j:], p[i:]) {
						return true
					}
				}
				// Failed.
				return false
			}

			// Move ahead one character at a time and try to
			// match at each position.
			for ; j < len(s); j++ {
				if match(s[j:], p[i:]) {
					return true
				}
			}
			// Failed.
			return false
		}

		// There must be at least one more character in the string.
		// If we are at the end, fail.
		if j == len(s) {
			return false
		}

		// Check if the next character of the string is acceptable.
		if p[i] != '?' && p[i] != s[j] {
			return false
		}

		// Move to the next character, both in string and in pattern.
		i++
		j++
	}
	// If at end of pattern, accept if also at end of string.
	return j == len(s)
}
