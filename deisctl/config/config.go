package config

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
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

// CheckConfig looks for a value at a keyspace path
// and returns an error if a value is not found
func CheckConfig(root string, k string) error {

	client, err := getEtcdClient()
	if err != nil {
		return err
	}

	_, err = doConfigGet(client, root, []string{k})
	if err != nil {
		return err
	}

	return nil
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

		// split k/v from args
		split := strings.Split(kv, "=")
		if len(split) != 2 {
			return result, fmt.Errorf("invalid argument: %v", kv)
		}
		k, v := split[0], split[1]

		// prepare path and value
		path := root + k
		var val string

		// special handling for sshKey
		if path == "/deis/platform/sshPrivateKey" {
			b64, err := readSSHPrivateKey(v)
			if err != nil {
				return result, err
			}
			val = b64
		} else {
			val = v
		}

		// set key/value in etcd
		ret, err := client.Set(path, val)
		if err != nil {
			return result, err
		}
		result = append(result, ret)

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

// readSSHPrivateKey reads the key file and returns a base64 encoded string
func readSSHPrivateKey(path string) (string, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}
