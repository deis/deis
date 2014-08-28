package config

import (
	"fmt"
	"strings"

	docopt "github.com/docopt/docopt-go"
)

// Config runs the config subcommand
func Config() error {
	usage := `Deis Cluster Configuration

    Usage:
    deisctl config <target> get [<key>...] [options]
    deisctl config <target> set <key=val>... [options]

    Options:
    --verbose                   print out the request bodies [default: false]
    `
	// parse command-line arguments
	args, err := docopt.Parse(usage, nil, true, "", true)
	if err != nil {
		return err
	}
	err = setConfigFlags(args)
	if err != nil {
		return err
	}
	return doConfig(args)
}

// Flags for config package
var Flags struct {
}

func setConfigFlags(args map[string]interface{}) error {
	return nil
}

func doConfig(args map[string]interface{}) error {
	client, err := getEtcdClient()
	if err != nil {
		return err
	}

	rootPath := "/deis/" + args["<target>"].(string) + "/"

	var vals []string
	if args["set"] == true {
		vals, err = doConfigSet(client, rootPath, args["<key=val>"].([]string))
	} else {
		vals, err = doConfigGet(client, rootPath, args["<key>"].([]string))
	}
	if err != nil {
		return err
	}

	// print results
	for _, v := range vals {
		fmt.Printf("%v\n", v)
	}
	return nil
}

func doConfigSet(client *etcdClient, root string, kvs []string) ([]string, error) {
	var result []string
	for _, kv := range kvs {
		split := strings.Split(kv, "=")
		if len(split) != 2 {
			return result, fmt.Errorf("invalid argument: %v", kv)
		}
		val, err := client.Set(root+split[0], split[1])
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}

func doConfigGet(client *etcdClient, root string, keys []string) ([]string, error) {
	var result []string
	for _, k := range keys {
		val, err := client.Get(root + k)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}
