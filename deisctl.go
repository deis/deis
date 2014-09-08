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
	client.Flags.Version = args["--version"].(bool)
	client.Flags.Endpoint = args["--endpoint"].(string)
	client.Flags.EtcdKeyPrefix = args["--etcd-key-prefix"].(string)
	client.Flags.EtcdKeyFile = args["--etcd-keyfile"].(string)
	client.Flags.EtcdCertFile = args["--etcd-certfile"].(string)
	client.Flags.EtcdCAFile = args["--etcd-cafile"].(string)
	//client.Flags.UseAPI = args["--experimental-api"].(bool)
	client.Flags.KnownHostsFile = args["--known-hosts-file"].(string)
	client.Flags.StrictHostKeyChecking = args["--strict-host-key-checking"].(bool)
	timeout, _ := strconv.ParseFloat(args["--request-timeout"].(string), 64)
	client.Flags.RequestTimeout = timeout
	tunnel := args["--tunnel"].(string)
	if tunnel != "" {
		client.Flags.Tunnel = tunnel
	} else {
		client.Flags.Tunnel = os.Getenv("DEISCTL_TUNNEL")
	}
}

func main() {
	usage := `Deis Control Utility

Usage:
  deisctl <command> [<target>...] [options]

Example Commands:

  deisctl install platform
  deisctl uninstall builder@1
  deisctl list
  deisctl scale router=2
  deisctl start router@2
  deisctl stop router builder
  deisctl status controller
  deisctl journal controller

Options:
  --debug                     print debug information to stderr
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
		err = cmd.ListUnits(c)
	case "list-units":
		err = cmd.ListUnits(c)
	case "list-unit-files":
		err = cmd.ListUnitFiles(c)
	case "scale":
		err = cmd.Scale(c, targets)
	case "start":
		err = cmd.Start(c, targets)
	case "restart":
		err = cmd.Restart(c, targets)
	case "stop":
		err = cmd.Stop(c, targets)
	case "status":
		err = cmd.Status(c, targets)
	case "journal":
		err = cmd.Journal(c, targets)
	case "install":
		err = cmd.Install(c, targets)
	case "uninstall":
		err = cmd.Uninstall(c, targets)
	case "config":
		err = cmd.Config()
	case "update":
		err = cmd.Update()
	case "refresh-units":
		err = cmd.RefreshUnits()
	default:
		fmt.Printf(usage)
		os.Exit(2)
	}
	if err != nil {
		exit(err, 1)
	}
}
