package parser

import (
	"fmt"

	"github.com/deis/deis/client-go/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Perms routes perms commands to their specific function.
func Perms(argv []string) error {
	usage := `
Valid commands for perms:

perms:list            list permissions granted on an app
perms:create          create a new permission for a user
perms:delete          delete a permission for a user

Use 'deis help perms:[command]' to learn more.
`
	if len(argv) < 2 {
		return permsList([]string{"perms:list"})
	}

	switch argv[1] {
	case "list":
		return permsList(combineCommand(argv))
	case "create":
		return permCreate(combineCommand(argv))
	case "delete":
		return permDelete(combineCommand(argv))
	case "--help":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func permsList(argv []string) error {
	usage := `
Lists all users with permission to use an app, or lists all users with system
administrator privileges.

Usage: deis perms:list [-a --app=<app>|--admin]

Options:
  -a --app=<app>
    lists all users with permission to <app>. <app> is the uniquely identifiable name
    for the application.
  --admin
    lists all users with system administrator privileges.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	admin := args["--admin"].(bool)

	return cmd.PermsList(safeGetValue(args, "--app"), admin)
}

func permCreate(argv []string) error {
	usage := `
Gives another user permission to use an app, or gives another user
system administrator privileges.

Usage: deis perms:create <username> [-a --app=<app>|--admin]

Arguments:
  <username>
    the name of the new user.

Options:
  -a --app=<app>
    grants <username> permission to use <app>. <app> is the uniquely identifiable name
    for the application.
  --admin
    grants <username> system administrator privileges.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	username := args["<username>"].(string)
	admin := args["--admin"].(bool)

	return cmd.PermCreate(app, username, admin)
}

func permDelete(argv []string) error {
	usage := `
Revokes another user's permission to use an app, or revokes another user's system
administrator privileges.

Usage: deis perms:delete <username> [-a --app=<app>|--admin]

Arguments:
  <username>
    the name of the user.

Options:
  -a --app=<app>
    revokes <username> permission to use <app>. <app> is the uniquely identifiable name
    for the application.
  --admin
    revokes <username> system administrator privileges.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	username := args["<username>"].(string)
	admin := args["--admin"].(bool)

	return cmd.PermDelete(app, username, admin)
}
