package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/deis/deis/client-go/parser"
	"github.com/deis/deis/version"
	docopt "github.com/docopt/docopt-go"
)

// main exits with the return value of Command(os.Args[1:]), deferring all logic to
// a func we can test.
func main() {
	os.Exit(Command(os.Args[1:]))
}

// Command routes deis commands to their proper parser.
func Command(argv []string) int {
	usage := `
The Deis command-line client issues API calls to a Deis controller.

Usage: deis <command> [<args>...]

Auth commands::

  register      register a new user with a controller
  login         login to a controller
  logout        logout from the current controller

Subcommands, use 'deis help [subcommand]' to learn more::

  apps          manage applications used to provide services
  ps            manage processes inside an app container
  config        manage environment variables that define app config
  domains       manage and assign domain names to your applications
  builds        manage builds created using 'git push'
  limits        manage resource limits for your application
  tags          manage tags for application containers
  releases      manage releases of an application
  certs         manage SSL endpoints for an app

  keys          manage ssh keys used for 'git push' deployments
  perms         manage permissions for applications
  git           manage git for applications
  users         manage users

Shortcut commands, use 'deis shortcuts' to see all::

  create        create a new application
  scale         scale processes by type (web=2, worker=1)
  info          view information about the current app
  open          open a URL to the app in a browser
  logs          view aggregated log info for the app
  run           run a command in an ephemeral app container
  destroy       destroy an application
  pull          imports an image and deploys as a new release

Use 'git push deis master' to deploy to an application.
`
	// Reorganize some command line flags and commands.
	argv = parseArgs(argv)
	// Give docopt an optional final false arg so it doesn't call os.Exit().
	args, err := docopt.Parse(usage, argv, false, version.Version, true, false)

	if err != nil {
		fmt.Println(err)
		return 1
	}

	if len(args) == 0 {
		return 0
	}

	command := args["<command>"]
	// Dispatch the command, passing the argv through so subcommands can
	// re-parse it according to their usage strings.
	switch command {
	case "auth":
		err = parser.Auth(argv)
	case "ps":
		err = parser.Ps(argv)
	case "apps":
		err = parser.Apps(argv)
	case "keys":
		err = parser.Keys(argv)
	case "users":
		err = parser.Users(argv)
	case "help":
		fmt.Print(usage)
		return 0
	default:
		parser.PrintUsage()
		return 1
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}
	return 0
}

// parseArgs returns the provided args with "--help" as the last arg if need be,
// expands shortcuts and formats commands to be properly routed.
func parseArgs(argv []string) []string {
	if len(argv) == 1 {
		// rearrange "deisctl --help" as "deisctl help"
		if argv[0] == "--help" || argv[0] == "-h" {
			argv[0] = "help"
		}
	}

	if len(argv) >= 2 {
		// Rearrange "deisctl help <command>" to "deisctl <command> --help".
		if argv[0] == "help" || argv[0] == "--help" || argv[0] == "-h" {
			argv = append(argv[1:], "--help")
		}
	}

	// Split commands with colons [apps:create, ...] -> [apps, create, ...].
	if len(argv) > 0 {
		argv[0] = replaceShortcut(argv[0])

		index := strings.Index(argv[0], ":")

		if index != -1 {
			command := argv[0]
			argv[0] = command[:index]
			argv = append(argv[:1], append([]string{command[(index + 1):]}, argv[1:]...)...)
		}
	}

	return argv
}

func replaceShortcut(command string) string {
	shortcuts := map[string]string{
		"create":         "apps:create",
		"destroy":        "apps:destroy",
		"info":           "apps:info",
		"login":          "auth:login",
		"logout":         "auth:logout",
		"logs":           "apps:logs",
		"open":           "apps:open",
		"passwd":         "auth:passwd",
		"pull":           "builds:create",
		"register":       "auth:register",
		"rollback":       "releases:rollback",
		"run":            "apps:run",
		"scale":          "ps:scale",
		"sharing":        "perms:list",
		"sharing:list":   "perms:list",
		"sharing:add":    "perms:create",
		"sharing:remove": "perms:delete",
		"whoami":         "auth:whoami",
	}

	expandedCommand := shortcuts[command]
	if expandedCommand == "" {
		return command
	}

	return expandedCommand
}
