package client

import (
	"errors"
	"fmt"
	"os"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/cmd"

	docopt "github.com/docopt/docopt-go"
)

// DeisCtlClient manages Deis components, configuration, and related tasks.
type DeisCtlClient interface {
	Config(argv []string) error
	Install(argv []string) error
	Journal(argv []string) error
	List(argv []string) error
	RefreshUnits(argv []string) error
	Restart(argv []string) error
	Scale(argv []string) error
	SSH(argv []string) error
	Start(argv []string) error
	Status(argv []string) error
	Stop(argv []string) error
	Uninstall(argv []string) error
}

// Client uses a backend to implement the DeisCtlClient interface.
type Client struct {
	Backend backend.Backend
}

// NewClient returns a Client using the requested backend.
// The only backend currently supported is "fleet".
func NewClient(requestedBackend string) (*Client, error) {
	var backend backend.Backend

	if requestedBackend == "" {
		requestedBackend = "fleet"
	}

	switch requestedBackend {
	case "fleet":
		b, err := fleet.NewClient()
		if err != nil {
			return nil, err
		}
		backend = b
	default:
		return nil, errors.New("invalid backend")
	}
	return &Client{Backend: backend}, nil
}

// Config gets or sets a configuration value from the cluster.
//
// A configuration value is stored and retrieved from a key/value store (in this case, etcd)
// at /deis/<component>/<config>. Configuration values are typically used for component-level
// configuration, such as enabling TLS for the routers.
func (c *Client) Config(argv []string) error {
	usage := `Gets or sets a configuration value from the cluster.

A configuration value is stored and retrieved from a key/value store
(in this case, etcd) at /deis/<component>/<config>. Configuration
values are typically used for component-level configuration, such as
enabling TLS for the routers.

Note: "deisctl config platform set sshPrivateKey=" expects a path
to a private key.

Usage:
  deisctl config <target> get [<key>...]
  deisctl config <target> set <key=val>...
  deisctl config <target> rm [<key>...]

Examples:
  deisctl config platform set domain=mydomain.com
  deisctl config platform set sshPrivateKey=$HOME/.ssh/deis
  deisctl config controller get webEnabled
  deisctl config controller rm webEnabled
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	var action string
	var key []string

	if args["set"] == true {
		action = "set"
		key = args["<key=val>"].([]string)
	} else if args["rm"] == true {
		action = "rm"
		key = args["<key>"].([]string)
	} else {
		action = "get"
		key = args["<key>"].([]string)
	}

	return cmd.Config(args["<target>"].(string), action, key)
}

// Install loads the definitions of components from local unit files.
// After Install, the components will be available to Start.
func (c *Client) Install(argv []string) error {
	usage := `Loads the definitions of components from local unit files.

After install, the components will be available to start.

"deisctl install" looks for unit files in these directories, in this order:
- the $DEISCTL_UNITS environment variable, if set
- $HOME/.deis/units
- /var/lib/deis/units

Usage:
  deisctl install [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Install(args["<target>"].([]string), c.Backend)
}

// Journal prints log output for the specified components.
func (c *Client) Journal(argv []string) error {
	usage := `Prints log output for the specified components.

Usage:
  deisctl journal [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Journal(args["<target>"].([]string), c.Backend)
}

// List prints a summary of installed components.
func (c *Client) List(argv []string) error {
	usage := `Prints a list of installed units.

Usage:
  deisctl list [options]
`
	// parse command-line arguments
	if _, err := docopt.Parse(usage, argv, true, "", false); err != nil {
		return err
	}
	return cmd.ListUnits(c.Backend)
}

// RefreshUnits overwrites local unit files with those requested.
func (c *Client) RefreshUnits(argv []string) error {
	usage := `Overwrites local unit files with those requested.

Downloading from the Deis project GitHub URL by tag or SHA is the only mechanism
currently supported.

"deisctl install" looks for unit files in these directories, in this order:
- the $DEISCTL_UNITS environment variable, if set
- $HOME/.deis/units
- /var/lib/deis/units

Usage:
  deisctl refresh-units [-p <target>] [-t <tag>]

Options:
  -p --path=<target>   where to save unit files [default: $HOME/.deis/units]
  -t --tag=<tag>       git tag, branch, or SHA to use when downloading unit files
                       [default: master]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}

	return cmd.RefreshUnits(args["--path"].(string), args["--tag"].(string))
}

// Restart stops and then starts components.
func (c *Client) Restart(argv []string) error {
	usage := `Stops and then starts the specified components.

Usage:
  deisctl restart [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Restart(args["<target>"].([]string), c.Backend)
}

// Scale grows or shrinks the number of running components.
func (c *Client) Scale(argv []string) error {
	usage := `Grows or shrinks the number of running components.

Currently "router", "registry" and "store-gateway" are the only types that can be scaled.

Usage:
  deisctl scale [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Scale(args["<target>"].([]string), c.Backend)
}

// SSH opens an interactive shell with a machine in the cluster.
func (c *Client) SSH(argv []string) error {
	usage := `Open an interactive shell on a machine in the cluster given a unit or machine id.

Usage:
  deisctl ssh <target>
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.SSH(args["<target>"].(string), c.Backend)
}

// Start activates the specified components.
func (c *Client) Start(argv []string) error {
	usage := `Activates the specified components.

Usage:
  deisctl start [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Start(args["<target>"].([]string), c.Backend)
}

// Status prints the current status of components.
func (c *Client) Status(argv []string) error {
	usage := `Prints the current status of components.

Usage:
  deisctl status [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Status(args["<target>"].([]string), c.Backend)
}

// Stop deactivates the specified components.
func (c *Client) Stop(argv []string) error {
	usage := `Deactivates the specified components.

Usage:
  deisctl stop [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Stop(args["<target>"].([]string), c.Backend)
}

// Uninstall unloads the definitions of the specified components.
// After Uninstall, the components will be unavailable until Install is called.
func (c *Client) Uninstall(argv []string) error {
	usage := `Unloads the definitions of the specified components.

After uninstall, the components will be unavailable until install is called.

Usage:
  deisctl uninstall [<target>...] [options]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Uninstall(args["<target>"].([]string), c.Backend)
}
