package parser

import (
	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Tags routes tags commands to their specific function
func Tags(argv []string) error {
	usage := `
Valid commands for tags:

tags:list        list tags for an app
tags:set         set tags for an app
tags:unset       unset tags for an app

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "tags:list":
		return tagsList(argv)
	case "tags:set":
		return tagsSet(argv)
	case "tags:unset":
		return tagsUnset(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "tags" {
			argv[0] = "tags:list"
			return tagsList(argv)
		}

		PrintUsage()
		return nil
	}
}

func tagsList(argv []string) error {
	usage := `
Lists tags for an application.

Usage: deis tags:list [options]

Options:
  -a --app=<app>
    the uniquely identifiable name of the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.TagsList(safeGetValue(args, "--app"))
}

func tagsSet(argv []string) error {
	usage := `
Sets tags for an application.

A tag is a key/value pair used to tag an application's containers and is passed to the
scheduler. This is often used to restrict workloads to specific hosts matching the
scheduler-configured metadata.

Usage: deis tags:set [options] <key>=<value>...

Arguments:
  <key> the tag key, for example: "environ" or "rack"
  <value> the tag value, for example: "prod" or "1"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	tags := args["<key>=<value>"].([]string)

	return cmd.TagsSet(app, tags)
}

func tagsUnset(argv []string) error {
	usage := `
Unsets tags for an application.

Usage: deis tags:unset [options] <key>...

Arguments:
  <key> the tag key to unset, for example: "environ" or "rack"

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	tags := args["<key>"].([]string)

	return cmd.TagsUnset(app, tags)
}
