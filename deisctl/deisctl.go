package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/deis/deis/deisctl"
	docopt "github.com/docopt/docopt-go"
)

func cmdList(c deisctl.Client) error {
	err := c.List()
	return err
}

func cmdScale(c deisctl.Client, targets []string) error {
	for _, target := range targets {
		component, num, err := splitScaleTarget(target)
		if err != nil {
			return err
		}
		err = c.Scale(component, num)
		if err != nil {
			return err
		}
	}
	return nil
}

func cmdStart(c deisctl.Client, targets []string) error {
	for _, target := range targets {
		err := c.Start(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func cmdStop(c deisctl.Client, targets []string) error {
	for _, target := range targets {
		err := c.Stop(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func cmdStatus(c deisctl.Client, targets []string) error {
	for _, target := range targets {
		err := c.Status(target)
		if err != nil {
			return err
		}
	}
	return nil
}

func cmdInstall(c deisctl.Client) error {
	targets := []string{
		"database=1",
		"cache=1",
		"logger=1",
		"registry=1",
		"controller=1",
		"builder=1",
		"router=1"}
	fmt.Println("Scheduling units...")
	err := cmdScale(c, targets)
	fmt.Println("Activating units...")
	err = cmdStart(c, []string{"registry", "logger", "cache", "database"})
	if err != nil {
		return err
	}
	err = cmdStart(c, []string{"controller"})
	if err != nil {
		return err
	}
	err = cmdStart(c, []string{"builder"})
	if err != nil {
		return err
	}
	err = cmdStart(c, []string{"router"})
	if err != nil {
		return err
	}
	fmt.Println("Done.")
	return err
}

func cmdUninstall(c deisctl.Client) error {
	targets := []string{
		"database=0",
		"cache=0",
		"logger=0",
		"registry=0",
		"controller=0",
		"builder=0",
		"router=0"}
	err := cmdScale(c, targets)
	return err
}

func exit(err error, code int) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(code)
}

func splitScaleTarget(target string) (c string, num int, err error) {
	r := regexp.MustCompile(`([a-z-]+)=([\d]+)`)
	match := r.FindStringSubmatch(target)
	if len(match) == 0 {
		err = fmt.Errorf("Could not parse: %v", target)
		return
	}
	c = match[1]
	num, err = strconv.Atoi(match[2])
	if err != nil {
		return
	}
	return
}

func setGlobalFlags(args map[string]interface{}) {
	deisctl.Flags.Debug = args["--debug"].(bool)
	verbosity, _ := strconv.Atoi(args["--verbosity"].(string))
	deisctl.Flags.Verbosity = verbosity
	deisctl.Flags.Endpoint = args["--endpoint"].(string)
	deisctl.Flags.EtcdKeyPrefix = args["--etcd-key-prefix"].(string)
	deisctl.Flags.KnownHostsFile = args["--known-hosts-file"].(string)
	deisctl.Flags.StrictHostKeyChecking = args["--strict-host-key-checking"].(bool)
	tunnel := args["--tunnel"]
	if tunnel != nil {
		deisctl.Flags.Tunnel = tunnel.(string)
	} else {
		deisctl.Flags.Tunnel = os.Getenv("FLEETCTL_TUNNEL")
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
	c, err := deisctl.NewClient()
	if err != nil {
		exit(err, 1)
	}
	// dispatch the command
	switch command {
	case "list":
		err = cmdList(c)
	case "scale":
		err = cmdScale(c, targets)
	case "start":
		err = cmdStart(c, targets)
	case "stop":
		err = cmdStop(c, targets)
	case "status":
		err = cmdStatus(c, targets)
	case "install":
		err = cmdInstall(c)
	case "uninstall":
		err = cmdUninstall(c)
	default:
		fmt.Printf(usage)
		os.Exit(2)
	}
	if err != nil {
		exit(err, 1)
	}
}
