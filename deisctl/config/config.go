package config

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/deis/deis/deisctl/utils"
)

// fileKeys define config keys to be read from local files
var fileKeys = []string{
	"/deis/platform/sshPrivateKey",
	"/deis/router/sslCert",
	"/deis/router/sslKey"}

// b64Keys define config keys to be base64 encoded before stored
var b64Keys = []string{"/deis/platform/sshPrivateKey"}

// Config runs the config subcommand
func Config(args map[string]interface{}) error {
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

func SetConfig(target string, key string, val string) error {
	path := "/deis/" + target + "/" + key
	return SetConfigUsingFullKey(path, val)
}

func SetConfigUsingFullKey(key string, val string) error {
	client, err := getEtcdClient()
	if err != nil {
		return err
	}
	_, err = client.Set(key, val)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig(target string, key string) (string, error) {
	var value string
	client, err := getEtcdClient()
	if err != nil {
		return value, err
	}
	path := "/deis/" + target + "/" + key
	value, err = client.Get(path)
	if err != nil {
		return value, err
	}
	return value, nil
}

func GetConfigWithDefault(target string, key string, defaultValue string) (string, error) {
	var value string
	client, err := getEtcdClient()
	if err != nil {
		return value, err
	}
	path := "/deis/" + target + "/" + key
	value, err = client.GetWithDefault(path, defaultValue)
	if err != nil {
		return value, err
	}
	return value, nil
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
		split := strings.SplitN(kv, "=", 2)
		k, v := split[0], split[1]

		// prepare path and value
		path := root + k
		val, err := valueForPath(path, v)
		if err != nil {
			return result, err
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

// valueForPath returns the canonical value for a user-defined path and value
func valueForPath(path string, v string) (string, error) {

	// check if path is part of fileKeys
	for _, p := range fileKeys {

		if path == p {

			// read value from filesystem
			bytes, err := ioutil.ReadFile(utils.ResolvePath(v))
			if err != nil {
				return "", err
			}

			// see if we should return base64 encoded value
			for _, pp := range b64Keys {
				if path == pp {
					return base64.StdEncoding.EncodeToString(bytes), nil
				}
			}

			return string(bytes), nil
		}
	}

	return v, nil

}
