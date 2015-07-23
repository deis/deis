package parser

import (
	"fmt"

	"github.com/deis/deis/client-go/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Limits routes limits commands to their specific function
func Limits(argv []string) error {
	usage := `
Valid commands for limits:

limits:list        list resource limits for an app
limits:set         set resource limits for an app
limits:unset       unset resource limits for an app

Use 'deis help [command]' to learn more.
`
	if len(argv) < 2 {
		return limitsList([]string{"limits:list"})
	}

	switch argv[1] {
	case "list":
		return limitsList(combineCommand(argv))
	case "set":
		return limitSet(combineCommand(argv))
	case "unset":
		return limitUnset(combineCommand(argv))
	case "--help":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func limitsList(argv []string) error {
	usage := `
Lists resource limits for an application.

Usage: deis limits:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.LimitsList(safeGetValue(args, "--app"))
}

func limitSet(argv []string) error {
	usage := `
Sets resource limits for an application.

A resource limit is a finite resource within a container which we can apply
restrictions to either through the scheduler or through the Docker API. This limit
is applied to each individual container, so setting a memory limit of 1G for an
application means that each container gets 1G of memory.

Usage: deis limits:set [options] <type>=<limit>...

Arguments:
  <type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.
  <limit>
    The limit to apply to the process type. By default, this is set to --memory.
    You can only set one type of limit per call.

    With --memory, units are represented in Bytes (B), Kilobytes (K), Megabytes
    (M), or Gigabytes (G). For example, 'deis limit:set cmd=1G' will restrict all
    "cmd" processes to a maximum of 1 Gigabyte of memory each.

    With --cpu, units are represented in the number of cpu shares. For example,
    'deis limit:set --cpu cmd=1024' will restrict all "cmd" processes to a
    maximum of 1024 cpu shares.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -c --cpu
    limits cpu shares.
  -m --memory
    limits memory. [default: true]
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	limits := args["<type>=<limit>"].([]string)
	limitType := "memory"

	if args["--cpu"].(bool) {
		limitType = "cpu"
	}

	return cmd.LimitsSet(app, limits, limitType)
}

func limitUnset(argv []string) error {
	usage := `
Unsets resource limits for an application.

Usage: deis limits:unset [options] [--memory | --cpu] <type>...

Arguments:
  <type>
    the process type as defined in your Procfile, such as 'web' or 'worker'.
    Note that Dockerfile apps have a default 'cmd' process type.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -c --cpu
    limits cpu shares.
  -m --memory
    limits memory. [default: true]
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	limits := args["<type>"].([]string)
	limitType := "memory"

	if args["--cpu"].(bool) {
		limitType = "cpu"
	}

	return cmd.LimitsUnset(app, limits, limitType)
}
