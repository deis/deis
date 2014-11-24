package commons

import (
	"fmt"
	"github.com/hecticjeff/procfile"
	"strings"
)

// ParseProcfile2Structure read a procfile and return the structure to store in the db
func ParseProcfile2Structure(content string) {
	var types []string
	proclist := procfile.Parse(content)
	for name, _ := range proclist {
		types = append(types, "\""+name+"\": 1")
	}

	fmt.Printf("{ %v }", strings.Join(types, ", "))
}
