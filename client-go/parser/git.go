package parser

import (
	"errors"
	"fmt"

	"github.com/deis/deis/client-go/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Git routes git commands to their specific function.
func Git(argv []string) error {
	usage := `
Valid commands for git:

git:remote          Adds git remote of application to repository

Use 'deis help [command]' to learn more.
`
	if len(argv) < 2 {
		return errors.New("'deis git' is not a valid command, try 'deis help git'")
	}

	switch argv[1] {
	case "remote":
		return gitRemote(combineCommand(argv))
	case "--help":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func gitRemote(argv []string) error {
	usage := `
Adds git remote of application to repository

Usage: deis git:remote [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -r --remote=REMOTE
    name of remote to create. [default: deis]
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.GitRemote(safeGetValue(args, "--app"), args["--remote"].(string))
}
