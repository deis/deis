package parser

import (
	"fmt"
	"strconv"

	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Releases routes releases commands to their specific function.
func Releases(argv []string) error {
	usage := `
Valid commands for releases:

releases:list        list an application's release history
releases:info        print information about a specific release
releases:rollback    return to a previous release

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "releases:list":
		return releasesList(argv)
	case "releases:info":
		return releasesInfo(argv)
	case "releases:rollback":
		return releasesRollback(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "releases" {
			argv[0] = "releases:list"
			return releasesList(argv)
		}

		PrintUsage()
		return nil
	}
}

func releasesList(argv []string) error {
	usage := `
Lists release history for an application.

Usage: deis releases:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
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

	return cmd.ReleasesList(safeGetValue(args, "--app"), results)
}

func releasesInfo(argv []string) error {
	usage := `
Prints info about a particular release.

Usage: deis releases:info <version> [options]

Arguments:
  <version>
    the release of the application, such as 'v1'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	version, err := versionFromString(args["<version>"].(string))

	if err != nil {
		return err
	}

	return cmd.ReleasesInfo(safeGetValue(args, "--app"), version)
}

func releasesRollback(argv []string) error {
	usage := `
Rolls back to a previous application release.

Usage: deis releases:rollback [<version>] [options]

Arguments:
  <version>
    the release of the application, such as 'v1'.

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	var version int

	if args["<version>"] == nil {
		version = -1
	} else {
		version, err = versionFromString(args["<version>"].(string))

		if err != nil {
			return err
		}
	}

	return cmd.ReleasesRollback(safeGetValue(args, "--app"), version)
}

func versionFromString(version string) (int, error) {
	if version[:1] == "v" {
		if len(version) < 2 {
			return -1, fmt.Errorf("%s is not in the form 'v#'", version)
		}

		return strconv.Atoi(version[1:])
	}

	return strconv.Atoi(version)
}
