package parser

import (
	"fmt"

	"github.com/deis/deis/client-go/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Users routes user commands to the specific function.
func Users(argv []string) error {
	usage := `
Valid commands for users:

users:list        list all registered users

Use 'deis help [command]' to learn more.
`
	if len(argv) < 2 {
		return usersList([]string{"users:list"})
	}

	switch argv[1] {
	case "list":
		return usersList(combineCommand(argv))
	case "--help":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func usersList(argv []string) error {
	usage := `
Lists all registered users.
Requires admin privilages.

Usage: deis users:list [options]

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}

	return cmd.UsersList(results)
}
