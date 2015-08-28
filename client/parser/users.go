package parser

import (
	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Users routes user commands to the specific function.
func Users(argv []string) error {
	usage := `
Valid commands for users:

users:list        list all registered users

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "users:list":
		return usersList(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "users" {
			argv[0] = "users:list"
			return usersList(argv)
		}

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
