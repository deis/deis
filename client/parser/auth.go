package parser

import (
	"fmt"

	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Auth routes auth commands to the specific function.
func Auth(argv []string) error {
	usage := `
Valid commands for auth:

auth:register          register a new user
auth:login             authenticate against a controller
auth:logout            clear the current user session
auth:passwd            change the password for the current user
auth:whoami            display the current user
auth:cancel            remove the current user account
auth:regenerate        regenerate user tokens

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "auth:register":
		return authRegister(argv)
	case "auth:login":
		return authLogin(argv)
	case "auth:logout":
		return authLogout(argv)
	case "auth:passwd":
		return authPasswd(argv)
	case "auth:whoami":
		return authWhoami(argv)
	case "auth:cancel":
		return authCancel(argv)
	case "auth:regenerate":
		return authRegenerate(argv)
	case "auth":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func authRegister(argv []string) error {
	usage := `
Registers a new user with a Deis controller.

Usage: deis auth:register <controller> [options]

Arguments:
  <controller>
    fully-qualified controller URI, e.g. 'http://deis.local3.deisapp.com/'

Options:
  --username=<username>
    provide a username for the new account.
  --password=<password>
    provide a password for the new account.
  --email=<email>
    provide an email address.
  --ssl-verify=false
    disables SSL certificate verification for API requests
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	email := safeGetValue(args, "--email")
	sslVerify := false

	if args["--ssl-verify"] != nil && args["--ssl-verify"].(string) == "true" {
		sslVerify = true
	}

	return cmd.Register(controller, username, password, email, sslVerify)
}

func authLogin(argv []string) error {
	usage := `
Logs in by authenticating against a controller.

Usage: deis auth:login <controller> [options]

Arguments:
  <controller>
    a fully-qualified controller URI, e.g. "http://deis.local3.deisapp.com/".

Options:
  --username=<username>
    provide a username for the account.
  --password=<password>
    provide a password for the account.
  --ssl-verify=false
    disables SSL certificate verification for API requests
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	controller := safeGetValue(args, "<controller>")
	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	sslVerify := false

	if args["--ssl-verify"] != nil && args["--ssl-verify"].(string) == "true" {
		sslVerify = true
	}

	return cmd.Login(controller, username, password, sslVerify)
}

func authLogout(argv []string) error {
	usage := `
Logs out from a controller and clears the user session.

Usage: deis auth:logout
`

	if _, err := docopt.Parse(usage, argv, true, "", false, true); err != nil {
		return err
	}

	return cmd.Logout()
}

func authPasswd(argv []string) error {
	usage := `
Changes the password for the current user.

Usage: deis auth:passwd [options]

Options:
  --password=<password>
    the current password for the account.
  --new-password=<new-password>
    the new password for the account.
  --username=<username>
    the account's username.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	newPassword := safeGetValue(args, "--new-password")

	return cmd.Passwd(username, password, newPassword)
}

func authWhoami(argv []string) error {
	usage := `
Displays the currently logged in user.

Usage: deis auth:whoami
`

	if _, err := docopt.Parse(usage, argv, true, "", false, true); err != nil {
		return err
	}

	return cmd.Whoami()
}

func authCancel(argv []string) error {
	usage := `
Cancels and removes the current account.

Usage: deis auth:cancel [options]

Options:
  --username=<username>
    provide a username for the account.
  --password=<password>
    provide a password for the account.
  --yes
    force "yes" when prompted.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	password := safeGetValue(args, "--password")
	yes := args["--yes"].(bool)

	return cmd.Cancel(username, password, yes)
}

func authRegenerate(argv []string) error {
	usage := `
Regenerates auth token, defaults to regenerating token for the current user.

Usage: deis auth:regenerate [options]

Options:
  -u --username=<username>
    specify user to regenerate. Requires admin privilages.
  --all
    regenerate token for every user. Requires admin privilages.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	username := safeGetValue(args, "--username")
	all := args["--all"].(bool)

	return cmd.Regenerate(username, all)
}
