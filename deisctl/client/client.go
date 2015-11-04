package client

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/deis/deis/deisctl/backend"
	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/cmd"
	"github.com/deis/deis/deisctl/config"
	"github.com/deis/deis/deisctl/config/etcd"
	"github.com/deis/deis/deisctl/units"

	docopt "github.com/docopt/docopt-go"
)

// DeisCtlClient manages Deis components, configuration, and related tasks.
type DeisCtlClient interface {
	Config(argv []string) error
	Install(argv []string) error
	Journal(argv []string) error
	List(argv []string) error
	Machines(argv []string) error
	RefreshUnits(argv []string) error
	Restart(argv []string) error
	Scale(argv []string) error
	SSH(argv []string) error
	Start(argv []string) error
	Status(argv []string) error
	Stop(argv []string) error
	Uninstall(argv []string) error
	UpgradePrep(argv []string) error
	UpgradeTakeover(argv []string) error
	RollingRestart(argv []string) error
}

// Client uses a backend to implement the DeisCtlClient interface.
type Client struct {
	Backend       backend.Backend
	configBackend config.Backend
}

// NewClient returns a Client using the requested backend.
// The only backend currently supported is "fleet".
func NewClient(requestedBackend string) (*Client, error) {
	var backend backend.Backend

	cb, err := etcd.NewConfigBackend()
	if err != nil {
		return nil, err
	}

	if requestedBackend == "" {
		requestedBackend = "fleet"
	}

	switch requestedBackend {
	case "fleet":
		b, err := fleet.NewClient(cb)
		if err != nil {
			return nil, err
		}
		backend = b
	default:
		return nil, errors.New("invalid backend")
	}

	return &Client{Backend: backend, configBackend: cb}, nil
}

// UpgradePrep prepares a running cluster to be upgraded
func (c *Client) UpgradePrep(argv []string) error {
	usage := `Prepare platform for graceful upgrade.

Usage:
  deisctl upgrade-prep [--stateless]

Options:
  --stateless  Use when the target platform is stateless
`
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	stateless, _ := args["--stateless"].(bool)

	return cmd.UpgradePrep(stateless, c.Backend)
}

// UpgradeTakeover gracefully restarts a cluster prepared with upgrade-prep
func (c *Client) UpgradeTakeover(argv []string) error {
	usage := `Complete the upgrade of a prepped cluster.

Usage:
  deisctl upgrade-takeover [--stateless]

Options:
  --stateless  Use when the target platform is stateless
`
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	stateless, _ := args["--stateless"].(bool)

	return cmd.UpgradeTakeover(stateless, c.Backend, c.configBackend)
}

// RollingRestart attempts a rolling restart of an instance unit
func (c *Client) RollingRestart(argv []string) error {
	usage := `Perform a rolling restart of an instance unit.

Usage:
  deisctl rolling-restart <target>
`
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.RollingRestart(args["<target>"].(string), c.Backend)
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

	switch {
	case args["set"] == true:
		action = "set"
		key = args["<key=val>"].([]string)
	case args["rm"] == true:
		action = "rm"
		key = args["<key>"].([]string)
	default:
		action = "get"
		key = args["<key>"].([]string)
	}

	return cmd.Config(args["<target>"].(string), action, key, c.configBackend)
}

// Install loads the definitions of components from local unit files.
// After Install, the components will be available to Start.
func (c *Client) Install(argv []string) error {
	usage := fmt.Sprintf(`Loads the definitions of components from local unit files.

After install, the components will be available to start.

"deisctl install" looks for unit files in these directories, in this order:
- the $DEISCTL_UNITS environment variable, if set
- $HOME/.deis/units
- /var/lib/deis/units

Usage:
  deisctl install [<target>...] [options]

Options:
  --router-mesh-size=<num>  Number of routers to be loaded when installing the platform [default: %d].
`, cmd.DefaultRouterMeshSize)
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	meshSizeArg, _ := args["--router-mesh-size"].(string)
	parsedValue, err := strconv.ParseUint(meshSizeArg, 0, 8)
	if err != nil || parsedValue < 1 {
		fmt.Print("Error: argument --router-mesh-size: invalid value, make sure the value is an integer between 1 and 255.\n")
		return err
	}
	cmd.RouterMeshSize = uint8(parsedValue)

	return cmd.Install(args["<target>"].([]string), c.Backend, c.configBackend, cmd.CheckRequiredKeys)
}

// Journal prints log output for the specified components.
func (c *Client) Journal(argv []string) error {
	usage := `Prints log output for the specified components.

Usage:
  deisctl journal [<target>...]
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
  deisctl list
`
	// parse command-line arguments
	if _, err := docopt.Parse(usage, argv, true, "", false); err != nil {
		return err
	}
	return cmd.ListUnits(c.Backend)
}

func (c *Client) Machines(argv []string) error {
	usage := `List the current hosts in the cluster


Usage:
  deisctl machines
`
	// parse command-line arguments
	if _, err := docopt.Parse(usage, argv, true, "", false); err != nil {
		return err
	}
	return cmd.ListMachines(c.Backend)
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

	return cmd.RefreshUnits(args["--path"].(string), args["--tag"].(string), units.URL)
}

// Restart stops and then starts components.
func (c *Client) Restart(argv []string) error {
	usage := `Stops and then starts the specified components.

Usage:
  deisctl restart [<target>...]
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
  deisctl scale [<target>...]
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

If an optional <command> is provided, that command is run remotely, and the results returned.

Usage:
  deisctl ssh <target> [<command>...]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", true)
	if err != nil {
		return err
	}

	target := args["<target>"].(string)
	// handle help explicitly since docopt parsing is relaxed
	if target == "--help" {
		fmt.Println(usage)
		os.Exit(0)
	}

	var vargs []string
	if v, ok := args["<command>"]; ok {
		vargs = v.([]string)
	}

	return cmd.SSH(target, vargs, c.Backend)
}

func (c *Client) Dock(argv []string) error {
	usage := `Connect to the named docker container and run commands on it.

This is equivalent to running 'docker exec -it <target> <command>'.

Usage:
  deisctl dock <target> [<command>...]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", true)
	if err != nil {
		return err
	}

	target := args["<target>"].(string)
	// handle help explicitly since docopt parsing is relaxed
	if target == "--help" {
		fmt.Println(usage)
		os.Exit(0)
	}

	var vargs []string
	if v, ok := args["<command>"]; ok {
		vargs = v.([]string)
	}

	return cmd.Dock(target, vargs, c.Backend)
}

// Start activates the specified components.
func (c *Client) Start(argv []string) error {
	usage := `Activates the specified components.

Usage:
  deisctl start [<target>...]
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
  deisctl status [<target>...]
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
  deisctl stop [<target>...]
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
  deisctl uninstall [<target>...]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, argv, true, "", false)
	if err != nil {
		return err
	}

	return cmd.Uninstall(args["<target>"].([]string), c.Backend)
}
