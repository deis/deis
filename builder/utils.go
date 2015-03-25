package builder

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// YamlToJSON takes an input yaml string, parses it and returns a string formatted as json.
func YamlToJSON(bytes []byte) (string, error) {
	var anomaly map[string]string

	if err := yaml.Unmarshal(bytes, &anomaly); err != nil {
		return "", err
	}

	retVal, err := json.Marshal(&anomaly)

	if err != nil {
		return "", err
	}

	return string(retVal), nil
}

// ParseConfig takes a response body from the controller and returns a Config object.
func ParseConfig(body []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(body, &config)
	return &config, err
}

// ParseDomain returns the domain field from the bytes of a build hook response.
func ParseDomain(bytes []byte) (string, error) {
	var hook BuildHookResponse
	if err := json.Unmarshal(bytes, &hook); err != nil {
		return "", err
	}

	if hook.Domains == nil {
		return "", fmt.Errorf("invalid application domain")
	}

	if len(hook.Domains) < 1 {
		return "", fmt.Errorf("invalid application domain")
	}

	return hook.Domains[0], nil
}

// ParseReleaseVersion returns the version field from the bytes of a build hook response.
func ParseReleaseVersion(bytes []byte) (int, error) {
	var hook BuildHookResponse
	if err := json.Unmarshal(bytes, &hook); err != nil {
		return 0, fmt.Errorf("invalid application json configuration")
	}

	if hook.Release == nil {
		return 0, fmt.Errorf("invalid application version")
	}

	return hook.Release["version"], nil
}

// GetDefaultType returns the default process types given a YAML byte array.
func GetDefaultType(bytes []byte) (string, error) {
	type YamlTypeMap struct {
		DefaultProcessTypes ProcessType `yaml:"default_process_types"`
	}

	var p YamlTypeMap

	if err := yaml.Unmarshal(bytes, &p); err != nil {
		return "", err
	}

	retVal, err := json.Marshal(&p.DefaultProcessTypes)

	if err != nil {
		return "", err
	}

	if len(p.DefaultProcessTypes) == 0 {
		return "{}", nil
	}

	return string(retVal), nil
}

// ParseControllerConfig returns configuration key/value pair strings from a config.
func ParseControllerConfig(bytes []byte) ([]string, error) {
	var controllerConfig Config
	if err := json.Unmarshal(bytes, &controllerConfig); err != nil {
		return []string{}, err
	}

	if controllerConfig.Values == nil {
		return []string{""}, nil
	}

	retVal := []string{}
	for k, v := range controllerConfig.Values {
		retVal = append(retVal, fmt.Sprintf(" -e %s=\"%v\"", k, v))
	}
	return retVal, nil
}
