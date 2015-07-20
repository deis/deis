package parser

import (
	"fmt"

	"github.com/deis/deis/client-go/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Builds routes build commands to their specific function.
func Builds(argv []string) error {
	usage := `
Valid commands for builds:

builds:list        list build history for an application
builds:create      imports an image and deploys as a new release

Use 'deis help [command]' to learn more.
`
	if len(argv) < 2 {
		return buildsList([]string{"builds:list"})
	}

	switch argv[1] {
	case "list":
		return buildsList(combineCommand(argv))
	case "create":
		return buildsCreate(combineCommand(argv))
	case "--help":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func buildsList(argv []string) error {
	usage := `
Lists build history for an application.

Usage: deis builds:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.BuildsList(safeGetValue(args, "--app"))
}

func buildsCreate(argv []string) error {
	usage := `
Creates a new build of an application. Imports an <image> and deploys it to Deis
as a new release. If a Procfile is present in the current directory, it will be used
as the default process types for this application.

Usage: deis builds:create <image> [options]

Arguments:
  <image>
    A fully-qualified docker image, either from Docker Hub (e.g. deis/example-go:latest)
    or from an in-house registry (e.g. myregistry.example.com:5000/example-go:latest).
    This image must include the tag.

Options:
  -a --app=<app>
    The uniquely identifiable name for the application.
  -p --procfile=<procfile>
    A YAML string used to supply a Procfile to the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	image := safeGetValue(args, "<image>")
	procfile := safeGetValue(args, "--procfile")

	return cmd.BuildsCreate(app, image, procfile)
}
