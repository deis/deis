// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"launchpad.net/~niemeyer/goyaml/beta"
)

// Config ---------------------------------------------------------------------

// Config represents a configuration with convenient access methods.
type Config struct {
	Root interface{}
}

// Get returns a nested config according to a dotted path.
func (cfg *Config) Get(path string) (*Config, error) {
	n, err := Get(cfg.Root, path)
	if err != nil {
		return nil, err
	}
	return &Config{Root: n}, nil
}

// Bool returns a bool according to a dotted path.
func (cfg *Config) Bool(path string) (bool, error) {
	n, err := Get(cfg.Root, path)
	if err != nil {
		return false, err
	}
	switch n := n.(type) {
	case bool:
		return n, nil
	case string:
		return strconv.ParseBool(n)
	}
	return false, typeMismatch("bool or string", n)
}

// Float64 returns a float64 according to a dotted path.
func (cfg *Config) Float64(path string) (float64, error) {
	n, err := Get(cfg.Root, path)
	if err != nil {
		return 0, err
	}
	switch n := n.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	case string:
		return strconv.ParseFloat(n, 64)
	}
	return 0, typeMismatch("float64, int or string", n)
}

// Int returns an int according to a dotted path.
func (cfg *Config) Int(path string) (int, error) {
	n, err := Get(cfg.Root, path)
	if err != nil {
		return 0, err
	}
	switch n := n.(type) {
	case float64:
		// encoding/json unmarshals numbers into floats, so we compare
		// the string representation to see if we can return an int.
		if i := int(n); fmt.Sprint(i) == fmt.Sprint(n) {
			return i, nil
		} else {
			return 0, fmt.Errorf("Value can't be converted to int: %v", n)
		}
	case int:
		return n, nil
	case string:
		if v, err := strconv.ParseInt(n, 10, 0); err == nil {
			return int(v), nil
		} else {
			return 0, err
		}
	}
	return 0, typeMismatch("float64, int or string", n)
}

// List returns a []interface{} according to a dotted path.
func (cfg *Config) List(path string) ([]interface{}, error) {
	n, err := Get(cfg.Root, path)
	if err != nil {
		return nil, err
	}
	if value, ok := n.([]interface{}); ok {
		return value, nil
	}
	return nil, typeMismatch("[]interface{}", n)
}

// Map returns a map[string]interface{} according to a dotted path.
func (cfg *Config) Map(path string) (map[string]interface{}, error) {
	n, err := Get(cfg.Root, path)
	if err != nil {
		return nil, err
	}
	if value, ok := n.(map[string]interface{}); ok {
		return value, nil
	}
	return nil, typeMismatch("map[string]interface{}", n)
}

// String returns a string according to a dotted path.
func (cfg *Config) String(path string) (string, error) {
	n, err := Get(cfg.Root, path)
	if err != nil {
		return "", err
	}
	switch n := n.(type) {
	case bool, float64, int:
		return fmt.Sprint(n), nil
	case string:
		return n, nil
	}
	return "", typeMismatch("bool, float64, int or string", n)
}

// typeMismatch returns an error for an expected type.
func typeMismatch(expected string, got interface{}) error {
	return fmt.Errorf("Type mismatch: expected %s; got %T", expected, got)
}

// Fetching -------------------------------------------------------------------

// Get returns a child of the given value according to a dotted path.
func Get(cfg interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	// Normalize path.
	for k, v := range parts {
		if v == "" {
			if k == 0 {
				parts = parts[1:]
			} else {
				return nil, fmt.Errorf("Invalid path %q", path)
			}
		}
	}
	// Get the value.
	for pos, part := range parts {
		switch c := cfg.(type) {
		case []interface{}:
			if i, error := strconv.ParseInt(part, 10, 0); error == nil {
				if int(i) < len(c) {
					cfg = c[i]
				} else {
					return nil, fmt.Errorf(
						"Index out of range at %q: list has only %v items",
						strings.Join(parts[:pos+1], "."), len(c))
				}
			} else {
				return nil, fmt.Errorf("Invalid list index at %q",
					strings.Join(parts[:pos+1], "."))
			}
		case map[string]interface{}:
			if value, ok := c[part]; ok {
				cfg = value
			} else {
				return nil, fmt.Errorf("Nonexistent map key at %q",
					strings.Join(parts[:pos+1], "."))
			}
		default:
			return nil, fmt.Errorf(
				"Invalid type at %q: expected []interface{} or map[string]interface{}; got %T",
				strings.Join(parts[:pos+1], "."), cfg)
		}
	}
	return cfg, nil
}

// Parsing --------------------------------------------------------------------

// Must is a wrapper for parsing functions to be used during initialization.
// It panics on failure.
func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

// normalizeValue normalizes a unmarshalled value. This is needed because
// encoding/json doesn't support marshalling map[interface{}]interface{}.
func normalizeValue(value interface{}) (interface{}, error) {
	switch value := value.(type) {
	case map[interface{}]interface{}:
		node := make(map[string]interface{}, len(value))
		for k, v := range value {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Unsupported map key: %#v", k)
			}
			item, err := normalizeValue(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported map value: %#v", v)
			}
			node[key] = item
		}
		return node, nil
	case map[string]interface{}:
		node := make(map[string]interface{}, len(value))
		for key, v := range value {
			item, err := normalizeValue(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported map value: %#v", v)
			}
			node[key] = item
		}
		return node, nil
	case []interface{}:
		node := make([]interface{}, len(value))
		for key, v := range value {
			item, err := normalizeValue(v)
			if err != nil {
				return nil, fmt.Errorf("Unsupported list item: %#v", v)
			}
			node[key] = item
		}
		return node, nil
	case bool, float64, int, string:
		return value, nil
	}
	return nil, fmt.Errorf("Unsupported type: %T", value)
}

// JSON -----------------------------------------------------------------------

// ParseJson reads a JSON configuration from the given string.
func ParseJson(cfg string) (*Config, error) {
	return parseJson([]byte(cfg))
}

// ParseJsonFile reads a JSON configuration from the given filename.
func ParseJsonFile(filename string) (*Config, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseJson(cfg)
}

// parseJson performs the real JSON parsing.
func parseJson(cfg []byte) (*Config, error) {
	var out interface{}
	var err error
	if err = json.Unmarshal(cfg, &out); err != nil {
		return nil, err
	}
	if out, err = normalizeValue(out); err != nil {
		return nil, err
	}
	return &Config{Root: out}, nil
}

// RenderJson renders a YAML configuration.
func RenderJson(cfg interface{}) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// YAML -----------------------------------------------------------------------

// ParseYaml reads a YAML configuration from the given string.
func ParseYaml(cfg string) (*Config, error) {
	return parseYaml([]byte(cfg))
}

// ParseYamlFile reads a YAML configuration from the given filename.
func ParseYamlFile(filename string) (*Config, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseYaml(cfg)
}

// parseYaml performs the real YAML parsing.
func parseYaml(cfg []byte) (*Config, error) {
	var out interface{}
	var err error
	if err = goyaml.Unmarshal(cfg, &out); err != nil {
		return nil, err
	}
	if out, err = normalizeValue(out); err != nil {
		return nil, err
	}
	return &Config{Root: out}, nil
}

// RenderYaml renders a YAML configuration.
func RenderYaml(cfg interface{}) (string, error) {
	b, err := goyaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
