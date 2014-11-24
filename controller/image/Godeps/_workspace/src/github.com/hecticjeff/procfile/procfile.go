package procfile

import (
	"regexp"
	"strings"
)

var procfileRegex = regexp.MustCompile("^([A-Za-z0-9_]+):\\s*(.+)$")

type Process struct {
	Command   string
	Arguments []string
}

// Parse will read the contents of a procfile and returns a map of string names
// to Process structs.
func Parse(procfile string) (procs map[string]Process) {
	procs = make(map[string]Process)

	for _, line := range strings.Split(procfile, "\n") {
		if matches := procfileRegex.FindStringSubmatch(line); matches != nil {
			name, command := matches[1], matches[2]
			args := strings.Split(command, " ")
			procs[name] = Process{args[0], args[1:]}
		}
	}

	return
}
