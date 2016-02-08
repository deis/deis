package parser

import (
	"fmt"

	"github.com/deis/deis/version"
	docopt "github.com/docopt/docopt-go"
)

// Version displays the client version
func Version(argv []string) error {
	usage := `
Displays the client version.

Usage: deis version

Use 'deis help [command]' to learn more.
`
	if _, err := docopt.Parse(usage, argv, true, "", false, true); err != nil {
		return err
	}

	fmt.Println(version.Version)
	return nil
}
