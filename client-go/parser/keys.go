package parser

import (
	"fmt"

	"github.com/deis/deis/client-go/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Keys routes key commands to the specific function.
func Keys(argv []string) error {
	usage := `
Valid commands for SSH keys:

keys:list        list SSH keys for the logged in user
keys:add         add an SSH key
keys:remove      remove an SSH key

Use 'deis help [command]' to learn more.
`
	if len(argv) < 2 {
		return keysList([]string{"keys:list"})
	}

	switch argv[1] {
	case "list":
		return keysList(combineCommand(argv))
	case "add":
		return keyAdd(combineCommand(argv))
	case "remove":
		return keyRemove(combineCommand(argv))
	case "--help":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func keysList(argv []string) error {
	usage := `
Lists SSH keys for the logged in user.

Usage: deis keys:list
`

	_, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.KeysList()
}

func keyAdd(argv []string) error {
	usage := `
Adds SSH keys for the logged in user.

Usage: deis keys:add [<key>]

Arguments:
  <key>
    a local file path to an SSH public key used to push application code.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	key := safeGetValue(args, "<key>")

	return cmd.KeyAdd(key)
}

func keyRemove(argv []string) error {
	usage := `
Removes an SSH key for the logged in user.

Usage: deis keys:remove <key>

Arguments:
  <key>
    the SSH public key to revoke source code push access.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	key := safeGetValue(args, "<key>")

	return cmd.KeyRemove(key)
}
