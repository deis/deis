package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/deis/deis/deisctl/backend/fleet"
	"github.com/deis/deis/deisctl/client"
	"github.com/deis/deis/deisctl/utils"

	docopt "github.com/docopt/docopt-go"
)

const (
	// Version of deisctl client
	Version string = "0.13.0-dev"
)

func exit(err error, code int) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(code)
}

func setGlobalFlags(args map[string]interface{}) {
	fleet.Flags.Endpoint = args["--endpoint"].(string)
	fleet.Flags.EtcdKeyPrefix = args["--etcd-key-prefix"].(string)
	fleet.Flags.EtcdKeyFile = args["--etcd-keyfile"].(string)
	fleet.Flags.EtcdCertFile = args["--etcd-certfile"].(string)
	fleet.Flags.EtcdCAFile = args["--etcd-cafile"].(string)
	//fleet.Flags.UseAPI = args["--experimental-api"].(bool)
	fleet.Flags.KnownHostsFile = args["--known-hosts-file"].(string)
	fleet.Flags.StrictHostKeyChecking = args["--strict-host-key-checking"].(bool)
	timeout, _ := strconv.ParseFloat(args["--request-timeout"].(string), 64)
	fleet.Flags.RequestTimeout = timeout
	tunnel := args["--tunnel"].(string)
	if tunnel != "" {
		fleet.Flags.Tunnel = tunnel
	} else {
		fleet.Flags.Tunnel = os.Getenv("DEISCTL_TUNNEL")
	}
}

func main() {
	deisctlMotd := utils.DeisIfy("Deis Control Utility")
	usage := deisctlMotd + `
Usage:
  deisctl <command> [<target>...] [options]

Commands:
  deisctl install [<service> | platform]
  deisctl uninstall [<service> | platform]
  deisctl list
  deisctl scale [<service>=<num>]
  deisctl start [<service> | platform]
  deisctl stop [<service> | platform]
  deisctl restart [<service> | platform]
  deisctl journal <service>
  deisctl config <component> <get|set> <args>
  deisctl update
  deisctl refresh-units

Example Commands:
  deisctl install platform
  deisctl uninstall builder@1
  deisctl scale router=2
  deisctl start router@2
  deisctl stop router builder
  deisctl status controller
  deisctl journal controller

Options:
  --version                   print version and exit
  --endpoint=<url>            etcd endpoint for fleet [default: http://127.0.0.1:4001]
  --etcd-key-prefix=<path>    keyspace for fleet data in etcd [default: /_coreos.com/fleet/]
  --etcd-keyfile=<path>       etcd key file authentication [default: ]
  --etcd-certfile=<path>      etcd cert file authentication [default: ]
  --etcd-cafile=<path>        etcd CA file authentication [default: ]
  --known-hosts-file=<path>   file used to store remote machine fingerprints [default: ~/.ssh/known_hosts]
  --strict-host-key-checking  verify SSH host keys [default: true]
  --tunnel=<host>             establish an SSH tunnel for communication with fleet and etcd [default: ]
  --request-timeout=<secs>    amount of time to allow a single request before considering it failed. [default: 3.0]
`
	// parse command-line arguments
	args, err := docopt.Parse(usage, nil, true, Version, true)
	if err != nil {
		exit(err, 2)
	}
	command := args["<command>"]
	targets := args["<target>"].([]string)
	setGlobalFlags(args)
	// construct a client
	c, err := client.NewClient("fleet")
	if err != nil {
		exit(err, 1)
	}
	// dispatch the command
	switch command {
	case "list":
		err = c.List()
	case "scale":
		err = c.Scale(targets)
	case "start":
		err = c.Start(targets)
	case "restart":
		err = c.Restart(targets)
	case "stop":
		err = c.Stop(targets)
	case "status":
		err = c.Status(targets)
	case "journal":
		err = c.Journal(targets)
	case "install":
		err = c.Install(targets)
	case "uninstall":
		err = c.Uninstall(targets)
	case "config":
		err = c.Config()
	case "update":
		err = c.Update()
	case "refresh-units":
		err = c.RefreshUnits()
	default:
		fmt.Printf(usage)
		os.Exit(2)
	}
	if err != nil {
		exit(err, 1)
	}
}
