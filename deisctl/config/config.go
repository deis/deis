package config

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/deis/deis/deisctl/utils"
)

// fileKeys define config keys to be read from local files
var fileKeys = []string{
	"/deis/platform/sshPrivateKey",
	"/deis/router/sslCert",
	"/deis/router/sslKey",
	"/deis/router/sslDhparam"}

// b64Keys define config keys to be base64 encoded before stored
var b64Keys = []string{"/deis/platform/sshPrivateKey"}

// Config runs the config subcommand
func Config(target string, action string, key []string, cb Backend) error {
	return doConfig(target, action, key, cb, os.Stdout)
}

// CheckConfig looks for a value at a keyspace path
// and returns an error if a value is not found
func CheckConfig(root string, k string, cb Backend) error {

	_, err := doConfigGet(cb, root, []string{k})
	if err != nil {
		return err
	}

	return nil
}

func doConfig(target string, action string, key []string, cb Backend, w io.Writer) error {
	rootPath := "/deis/" + target + "/"

	var vals []string
	var err error

	switch action {
	case "rm":
		vals, err = doConfigRm(cb, rootPath, key)
	case "set":
		vals, err = doConfigSet(cb, rootPath, key)
	default:
		vals, err = doConfigGet(cb, rootPath, key)
	}
	if err != nil {
		return err
	}

	// print results
	for _, v := range vals {
		fmt.Fprintf(w, "%v\n", v)
	}
	return nil
}

func doConfigSet(cb Backend, root string, kvs []string) ([]string, error) {
	var result []string
	regex := regexp.MustCompile(`^(.+)=([\s\S]+)$`)

	for _, kv := range kvs {

		if !regex.MatchString(kv) {
			return []string{}, fmt.Errorf("'%s' does not match the pattern 'key=var', ex: foo=bar\n", kv)
		}

		// split k/v from args
		captures := regex.FindStringSubmatch(kv)
		k, v := captures[1], captures[2]

		// prepare path and value
		path := root + k
		val, err := valueForPath(path, v)
		if err != nil {
			return result, err
		}

		// set key/value in config backend
		ret, err := cb.Set(path, val)
		if err != nil {
			return result, err
		}
		result = append(result, ret)

	}
	return result, nil
}

func doConfigGet(cb Backend, root string, keys []string) ([]string, error) {
	var result []string
	for _, k := range keys {
		val, err := cb.Get(root + k)
		if err != nil {
			return result, err
		}
		result = append(result, val)
	}
	return result, nil
}

func doConfigRm(cb Backend, root string, keys []string) ([]string, error) {
	var result []string
	for _, k := range keys {
		err := cb.Delete(root + k)
		if err != nil {
			return result, err
		}
		result = append(result, k)
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
