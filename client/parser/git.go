package parser

import (
	"fmt"

	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Git routes git commands to their specific function.
func Git(argv []string) error {
	usage := `
Valid commands for git:

git:remote          Adds git remote of application to repository

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "git:remote":
		return gitRemote(argv)
	case "git":
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
