package parser

import (
	"github.com/deis/deis/client/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Certs routes certs commands to their specific function.
func Certs(argv []string) error {
	usage := `
Valid commands for certs:

certs:list            list SSL certificates for an app
certs:add             add an SSL certificate to an app
certs:remove          remove an SSL certificate from an app

Use 'deis help [command]' to learn more.
`

	switch argv[0] {
	case "certs:list":
		return certsList(argv)
	case "certs:add":
		return certAdd(argv)
	case "certs:remove":
		return certRemove(argv)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "certs" {
			argv[0] = "certs:list"
			return certsList(argv)
		}

		PrintUsage()
		return nil
	}
}

func certsList(argv []string) error {
	usage := `
Show certificate information for an SSL application.

Usage: deis certs:list [options]

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

	return cmd.CertsList(results)
}

func certAdd(argv []string) error {
	usage := `
Binds a certificate/key pair to an application.

Usage: deis certs:add <cert> <key> [options]

Arguments:
  <cert>
    The public key of the SSL certificate.
  <key>
    The private key of the SSL certificate.

Options:
  --common-name=<cname>
    The common name of the certificate. If none is provided, the controller will
    interpret the common name from the certificate.
  --subject-alt-names=<sans>
    The subject alternate names (SAN) of the certificate, separated by commas. This will
    create multiple Certificate objects in the controller, one for each SAN.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	cert := args["<cert>"].(string)
	key := args["<key>"].(string)
	commonName := safeGetValue(args, "--common-name")
	sans := safeGetValue(args, "--subject-alt-names")

	return cmd.CertAdd(cert, key, commonName, sans)
}

func certRemove(argv []string) error {
	usage := `
removes a certificate/key pair from the application.

Usage: deis certs:remove <cn> [options]

Arguments:
  <cn>
    the common name of the cert to remove from the app.
`

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmd.CertRemove(safeGetValue(args, "<cn>"))
}
