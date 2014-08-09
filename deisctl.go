package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/deis/deisctl/client"
	"github.com/deis/deisctl/cmd"
	docopt "github.com/docopt/docopt-go"
)

func exit(err error, code int) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(code)
}

func setGlobalFlags(args map[string]interface{}) {
	client.Flags.Debug = args["--debug"].(bool)
	verbosity, _ := strconv.Atoi(args["--verbosity"].(string))
	client.Flags.Verbosity = verbosity
	client.Flags.Endpoint = args["--endpoint"].(string)
	client.Flags.EtcdKeyPrefix = args["--etcd-key-prefix"].(string)
	client.Flags.KnownHostsFile = args["--known-hosts-file"].(string)
	client.Flags.StrictHostKeyChecking = args["--strict-host-key-checking"].(bool)
	tunnel := args["--tunnel"]
	if tunnel != nil {
		client.Flags.Tunnel = tunnel.(string)
	} else {
		client.Flags.Tunnel = os.Getenv("FLEETCTL_TUNNEL")
	}
}

func main() {
	usage := `Deis Control Utility

Usage:
  deisctl <command> [<target>...] [options]

Example Commands:

  deisctl install
  deisctl uninstall
  deisctl list
  deisctl scale router=2
  deisctl start router.2
  deisctl stop router builder
  deisctl status controller

Options:
  --debug                     print debug information to stderr
  --endpoint=<url>            etcd endpoint for fleet [default: http://127.0.0.1:4001]
  --etcd-key-prefix=<path>    keyspace for fleet data in etcd [default: /_coreos.com/fleet/]
  --known-hosts-file=<path>   file used to store remote machine fingerprints [default: ~/.fleetctl/known_hosts]
  --strict-host-key-checking  verify SSH host keys [default: true]
  --tunnel=<host>             establish an SSH tunnel for communication with fleet and etcd
  --verbosity=<level>         log at a specified level of verbosity to stderr [default: 0]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, nil, true, "", true)
	if err != nil {
		exit(err, 2)
	}
	command := args["<command>"]
	targets := args["<target>"].([]string)
	setGlobalFlags(args)
	// construct a client
	c, err := client.NewClient()
	if err != nil {
		exit(err, 1)
	}
	// dispatch the command
	switch command {
	case "list":
		err = cmd.List(c)
	case "scale":
		err = cmd.Scale(c, targets)
	case "start":
		err = cmd.Start(c, targets)
	case "stop":
		err = cmd.Stop(c, targets)
	case "status":
		err = cmd.Status(c, targets)
	case "install":
		err = cmd.Install(c, targets)
	case "uninstall":
		err = cmd.Uninstall(c, targets)
	case "update":
		cmd.Update(os.Args)
	default:
		fmt.Printf(usage)
		os.Exit(2)
	}
	if err != nil {
		exit(err, 1)
	}
}
