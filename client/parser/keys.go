package parser

import (
	"github.com/deis/deis/client/cmd"
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

	switch argv[0] {
	case "keys:list":
		return keysList(argv)
	case "keys:add":
		return keyAdd(argv)
	case "keys:remove":
		return keyRemove(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "keys" {
			argv[0] = "keys:list"
			return keysList(argv)
		}

		PrintUsage()
		return nil
	}
}

func keysList(argv []string) error {
	usage := `
Lists SSH keys for the logged in user.

Usage: deis keys:list [options]

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

	return cmd.KeysList(results)
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
