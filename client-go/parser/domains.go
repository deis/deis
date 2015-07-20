package parser

import (
	"fmt"

	"github.com/deis/deis/client-go/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Domains routes domain commands to their specific function.
func Domains(argv []string) error {
	usage := `
Valid commands for domains:

domains:add           bind a domain to an application
domains:list          list domains bound to an application
domains:remove        unbind a domain from an application

Use 'deis help [command]' to learn more.
`
	if len(argv) < 2 {
		return domainsList([]string{"domains:list"})
	}

	switch argv[1] {
	case "add":
		return domainsAdd(combineCommand(argv))
	case "list":
		return domainsList(combineCommand(argv))
	case "remove":
		return domainsRemove(combineCommand(argv))
	case "--help":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage()
		return nil
	}
}

func domainsAdd(argv []string) error {
	usage := `
Binds a domain to an application.

Usage: deis domains:add <domain> [options]

Arguments:
  <domain>
    the domain name to be bound to the application, such as 'domain.deisapp.com'.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DomainsAdd(safeGetValue(args, "--app"), safeGetValue(args, "<domain>"))
}

func domainsList(argv []string) error {
	usage := `
Lists domains bound to an application.

Usage: deis domains:list [options]

Options:
	-a --app=<app>
		the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DomainsList(safeGetValue(args, "--app"))
}

func domainsRemove(argv []string) error {
	usage := `
Unbinds a domain for an application.

Usage: deis domains:remove <domain> [options]

Arguments:
  <domain>
    the domain name to be removed from the application.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.DomainsRemove(safeGetValue(args, "--app"), safeGetValue(args, "<domain>"))
}
